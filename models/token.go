package models

import (
	"github.com/jmcvetta/neoism"
	"strings"
	"fmt"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
	"snaphyAuth/Interfaces"
	"github.com/astaxie/beego"
	"snaphyAuth/errorMessage"
	"snaphyAuth/helper"
	"go/token"
)



type Token struct{
	IAT int64 //Issued at
	ISS int64 //User Identity
	EXP int64 //Expiry Time
	JTI string //Unique string identifies a token
	GRP string //Group
	REALM string //realm name
	KID string //AppId its not applicationID but AppId in TokenHelper file to track application.....
	STATUS string //Status showing the token is invalid or what.
}


//Check interface implementation..
//Will throw error if the struct doesn't implements Graph Interface..
var _ Interfaces.Graph = (*Token)(nil)


func init(){
	token := new(Token)
	err := token.AddUniqueConstraint()
	if err != nil{
		beego.Trace("Error creating UNIQUE constraint on Token database")
		beego.Trace(err)
	}
}




func (token *Token)AddUniqueConstraint() (err error){
	stmt := "CREATE CONSTRAINT ON (token:Token) ASSERT token.JTI IS UNIQUE"
	cq := neoism.CypherQuery{
		Statement: stmt,
	}
	// Issue the query.
	err = db.Cypher(&cq)
	return
}



//Add a tag to token with relationship..
func (token *Token) AddTag(tokenTag *TokenTag) (err error){
	stmt := `MATCH (tokenTag: TokenTag) WHERE tokenTag.id = {tokenTagId}
	         MATCH (token:Token) WHERE token.JTI = {JTI} AND token.KID = {KID}
	         MERGE (tokenTag) - [role:Role] -> (token)`

	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"tokenTagId": tokenTag.Id, "KID": token.KID, "JTI": token.JTI},
	}

	// Issue the query.
	err = db.Cypher(&cq)
	return
}



//Verify the token before parsing .. the token just checks if the tokens ia a valid one.
//NOTE: This method just verify the encryption of algorithm and remains silent about expiry of tokens or token not present in the graph.
func  (token *Token)  VerifyHash(tokenString string) (ok bool, err error){
	if tokenString == ""{
		return false, errors.New(errorMessage.TokenNotValid)
	}

	//Now check if the application status is active or not...
	tokenHelper := new(TokenHelper)

	tokenHelper, err = token.GetTokenHelper(token.KID)
	if ok, err = tokenHelper.CheckAppStatus(); ok && err == nil{
		parts := strings.Split(tokenString, ".")

		method := jwt.GetSigningMethod(tokenHelper.HashType)
		err = method.Verify(strings.Join(parts[0:2], "."), parts[2], []byte(tokenHelper.PublicKey))

		if err != nil{
			return false, err
		}else{
			//Valid return..
			return true, nil
		}
	}

	return
}




//Parses the token value..And also validates the algorithm..
//Return invalid if any error occures..
func (loginToken *Token) VerifyAndParse(tokenString string) (valid bool, err error){
	var token *jwt.Token
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		var (
			publicKey []byte
			ok bool
		)
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}


		//Parse the tokens claims and put it in model
		loginToken.KID = token.Claims["kid"].(string)
		loginToken.JTI = token.Claims["jti"].(string)
		loginToken.EXP = token.Claims["exp"].(int64)
		loginToken.IAT = token.Claims["iat"].(int64)
		loginToken.GRP = token.Claims["grp"].(string)
		loginToken.ISS = token.Claims["iss"].(int64)
		loginToken.REALM = token.Claims["realm"].(string)

		tokenHelper := new(TokenHelper)
		//Complete check if the token is valid..
		ok, err = loginToken.CheckIfTokenValid(tokenHelper)
		if ok == false || err != nil{
			return nil, err
		}else{
			//Now get the public key..
			publicKey, err = loginToken.LookUpKey(loginToken.KID, tokenHelper)
			if err != nil{
				return nil, err
			}else{
				//Token is valid return..
				return publicKey, nil
			}
		}
	})
	return token.Valid, err
}


//Find the public key with the provided data..
func (token *Token)LookUpKey(appId string, tokenHelper *TokenHelper) (publicKey []byte, err error){
	if appId == ""{
		return nil, errors.New("Error: appId value is null")
	}

	if tokenHelper.Id == 0{
		tokenHelper = new(TokenHelper)

		tokenHelper, err = token.GetTokenHelper(appId)
	}

	if err != nil {
		return nil, err
	}else{
		return []byte(tokenHelper.PublicKey), nil
	}
}



func (token *Token) GetTokenHelper(appId string) (tokenHelper *TokenHelper, err error){
	if appId == ""{
		return nil, errorMessage.AppIdNull
	}

	tokenHelper = new(TokenHelper)
	tokenHelper.AppId = appId
	err = tokenHelper.GetToken()
	return
}



//Complete method for checking if the token is valid or not..
//Required a parsed token..
func (token *Token)CheckIfTokenValid(tokenHelper *TokenHelper) (ok bool, err error){
	if tokenHelper.Id == 0{
		tokenHelper = new(TokenHelper)
		tokenHelper, err = token.GetTokenHelper(token.KID)
	}

	//Now check if Application of User is Active or Not
	if ok, err = tokenHelper.CheckAppStatus(); ok && err == nil{
		//Now check token expiry status..
		ok, err = token.CheckTokenExpiry()
		return

	}else{
		//Application is Not in Active state or TokenHelper value is disabled..
		return false, err
	}
}


//This doesn't checks if the main Application is Active or not.
//Check if token is expired or present in the node or status is invalid or what not ...
//TRUE IF TOKEN IS VALID AND FALSE IF NOT
func (token *Token) CheckTokenExpiry() (ok bool, err error){
	if token.EXP == 0 {
		return true, errors.New("Expiry claim not present in Token model.")
	}else{
		now := time.Now().Unix()
		if now >= token.EXP{
			//Reject the token
			return false, nil
		}else{
			//Now check if the token is present in the database....
			err := token.Read()
			if err != nil {
				return false, err
			}else{
				//Now check the status of the token..
				if token.STATUS != StatusMap["ACTIVE"] {
					return false, err
				}else{
					return true, err
				}
			}
		}
	}
}



//Read the token by token JTI token..
func (token *Token) Read() (err error){
	var(
		tokenList []*Token
	)
	if token.JTI == "" {
		return errorMessage.TokenJTINotPresent
	}

	err = token.ReadAll(&tokenList)
	if len(tokenList) != 0{
		token = tokenList[0]
	}
	return
}



//Find the token from the database and populate the data..
//Provide the nodeToken with jwt field or by userIdentity.
func (token *Token) ReadAll(tokenTagListInterface [] *interface{}) (err error){
	var tokenList []*Token
	if token.JTI != "" {
		stmt := `MATCH (token:Token) WHERE token.JTI = {JTI} RETURN token.IAT AS IAT, token.ISS AS ISS, token.EXP AS EXP, token.JTI AS JTI, token.GRP AS GRP, token.KID AS KID`
		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"JTI": token.JTI},
			Result: &tokenList,
		}

		// Issue the query.
		err = db.Cypher(&cq)

	}else if(token.ISS != 0){
		stmt := `MATCH (token:Token) WHERE token.ISS = {ISS} RETURN token.IAT AS IAT, token.ISS AS ISS, token.EXP AS EXP, token.JTI AS JTI, token.GRP AS GRP, token.KID AS KID`
		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"ISS": token.ISS},
			Result: &tokenList,
		}

		// Issue the query.
		err = db.Cypher(&cq)

	}else{
		panic("Unhandled condition.")
	}

	if len(tokenList) != 0 {
		tokenTagListInterface = &tokenList
		return
	}
}


//Check if token exist by checking JTI value
func (token *Token) Exist() (exist bool, err error)  {
	var tokenExist []struct{
		Count int `json:"count"`
	}
	if token.JTI == "" {
		return false, errorMessage.TokenJTINotPresent
	}

	stmt := `MATCH (token: Token) WHERE token.JTI = {JTI} RETURN count(tag) AS count`

	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"JTI": token.JTI },
		Result: &tokenExist,
	}

	// Issue the query.
	err = db.Cypher(&cq)
	if err != nil{
		return false, err
	}

	if len(tokenExist) != 0{
		if tokenExist[0].Count == 0{
			return false, err
		}else{
			return true, err
		}
	}else{
		return false, err
	}
}


//Delete the token by JTI value..
func (token *Token) Delete(err error){
	if token.JTI == "" {
		return false, errorMessage.TokenJTINotPresent
	}

	stmt := `MATCH (token: Token) WHERE token.JTI = {JTI}
	         OPTIONAL MATCH p = (token)-[*]->(END)
	         DETACH DELETE p`

	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"JTI": token.JTI },
	}

	// Issue the query.
	err = db.Cypher(&cq)
	return
}


//Only update status..
func (token *Token)Update(err error){
	//First check if token exist or not.
	var  exist bool = false
	exist, err = token.Exist()
	if exist == false {
		return errorMessage.TokenNotValid
	}else if err != nil{
		return err
	}

	if token.STATUS == "" {
		return errorMessage.TokenStatusNotPresent
	}

	stmt := `MATCH (token: Token) WHERE token.JTI = {JTI}
	         SET token.STATUS = {status}`
	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"JTI": token.JTI, "STATUS": token.STATUS },
	}
	// Issue the query.
	err = db.Cypher(&cq)
	return
}


//Throw error if not exist..
func (token *Token) CreateIfNotExist() (err error){

	tokenHelper := new (TokenHelper)
	//Find tokenHelper, Application and ApplicationSettings
	err = getTokenHelperApp(tokenHelper, token)

	if err != nil{
		return err
	}

	if token.IAT == 0{
		token.IAT = time.Now().Unix() //Issuer at time..

	}
	if token.ISS == 0{
		return errorMessage.TokenUserIdNotPresent
	}
	if token.JTI == ""{
		token.JTI = helper.CreateUUID()
	}
	if token.KID == ""{
		return errorMessage.AppIdNull
	}
	if token.EXP == 0{

		if tokenHelper.Application.Settings != nil{
			duration := tokenHelper.Application.Settings.ExpiryDuration
			if duration == 0 {
				//Default expiry after 1 hour
				duration = beego.AppConfig.DefaultInt("jwt::expiry", 3600)
			}
			expiryTime := time.Now().Add(time.Second * time.Duration(duration)).Unix() //time after it will be invalid..
			token.EXP = expiryTime
		}

	}

	if token.GRP != ""{
		return errorMessage.TokenGroupNotPresent
	}
	if token.REALM != ""{
		return errorMessage.TokenRealmNotPresent
	}
	//Putting default status of token to active..
	token.STATUS = StatusMap["ACTIVE"]

	var exist bool
	//Also create relationship.
	if exist, err = token.Exist(); err == nil && exist == false {
		stmt := `MATCH (grp:Group{name: {groupName}, appId: {appId}, realmName: {realmName}})
			 CREATE (token:Token{iat: {IAT}, iss: {ISS}, exp: {EXP}, jti: {JTI}, grp: {GRP}, realm: {REALM}, kid: {KID}, status:{STATUS} })
			 CREATE (grp) - [type: Type{userIdentity: {ISS} }] -> (grp)`
		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{
				"groupName": token.GRP,
				"appId": tokenHelper.Application.Id,
				"realmName": token.REALM,
				"IAT": token.IAT,
				"ISS": token.ISS,
				"EXP": token.EXP,
				"JTI": token.JTI,
				"GRP": token.GRP,
				"realm": token.REALM,
				"KID": token.KID,
				"STATUS": token.STATUS,
			},
		}

		// Issue the query.
		err = db.Cypher(&cq)

		return
	}else{
		return errorMessage.ErrorAlreadyPresent
	}
}


//Add a logout token method..delete the whole token tree..

func (token *Token) RefreshToken(previousToken *Token) (err error){
	//Renew token after one is invalid..
	//Find the previous token if not found the reject
	//Now if the status of previous token is already other than Active then reject RefreshToken request
	//Update the status of the previous token to Expired or Invalid
	//Now RENEW a new token adding a `Refresh` relationship to the token..
	panic("Under Construction Method")

}


func (token *Token) GenerateLoginToken(previousToken *Token) (err error){
	//First find the previous token and silently ignore if not present however if present
	//if previous token is provided  first delete the point upto which the previousToken including the previousToken node
	//Now find the node through which previous token was connected it may be a group or a token handle both cases differently...
	//Now Refresh the new token with the previous token
	//IF it was a group the simple create a new node and delete the previous one or if its a token and refresh a new token adding it base as previous
	panic("Under Construction Method")
}

//Will simply add token data to graph database with depedencies..
//token does not support merge create.
//It will create only if token doesnot exist and will throw an error if it exist..
func (token *Token)Create() (err error){
	err = token.CreateIfNotExist()
	return
}



func getTokenHelperApp(tokenHelper *TokenHelper, token *Token) (err error){
	tokenHelper, err = token.GetTokenHelper(token.KID)
	if err != nil{
		return err
	}
	//Now fetch application..
	_, err = tokenHelper.FetchTokenHelperApp()
	if err != nil{
		return err
	}
	//Now fetch the application settings...
	_, err = tokenHelper.Application.FetchAppSettings()
	if err != nil{
		return err
	}
	return
}





//Create first time Token and Generate Signature..
//Generate Signature for the token and Also add tokens to graph Database..
func (token *Token) GenerateSignature(tokenHelper *TokenHelper, realm *Realm, group *Group, tag *TokenTag, userIdentity string) (tokenString string, err error){
	// Create the token
	var signToken *jwt.Token

	if tokenHelper.HashType == "" {
		return errors.New("No Signing method present in the TokenHelper Database")
	}else if(tokenHelper.HashType == "RS256"){
		signToken = jwt.New(jwt.SigningMethodRS256)
	}else{
		//Right now default is RS256
		//TODO Add some more signing methods..
		signToken = jwt.New(jwt.SigningMethodRS256)
	}


	token.ISS = userIdentity
	token.GRP = group.Name
	token.KID = tokenHelper.AppId
	token.REALM = realm.Name

	//Now store token to graph Database...
	err = token.CreateIfNotExist()
	if err != nil{
		return
	}

	//Now add tag to database if..tag present
	if tag.Name != ""{
		token.AddTag(tag)
	}

	signToken.Claims["exp"]   = token.EXP
	signToken.Claims["iat"]   = token.IAT
	signToken.Claims["iss"]   = userIdentity //user identity..
	signToken.Claims["grp"]   = group.Name //Group with which user belong to.
	signToken.Claims["realm"] = realm.Name
	signToken.Claims["jti"]   = token.JTI
	signToken.Claims["kid"]   = tokenHelper.AppId

	//TODO LATER ADD MORE ROLES BY FETCHING RELATION FROM GRAPH DATABASE..
	if tag.Name != ""{
		signToken.Claims["roles"] = tag.Name
	}



	if tokenHelper.PrivateKey != "" {
		// Sign and get the complete encoded token as a string
		tokenString, err = signToken.SignedString([]byte(tokenHelper.PrivateKey))
		//Finally return the tokenString..
		return tokenString, err
	}else{
		return  "", errors.New("Private key not present in tokenHelper")
	}
}




//TODO WRITE ALL GRAPHQL METHODS....
//TODO WRITE A CRON METHOD TO REMOVE TOKENS AFTER THEIR DEFAULT DAYS that is 30 days-60 days




