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
)



type Token struct{
	IAT int64 //Issued at
	ISS int64 //User Identity
	EXP int64 //Expiry Time
	JTI string //Unique string identifies a token
	GRP string //Group
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



//Add a tag with relationship..
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

		tokenHelper := new (TokenHelper)
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



//TODO WRITE A METHOD TO UPDATE STATUS IF FOUND INVALID..TOKEN
//TODO WRITE A METHOD TO FIND IF THE TOKEN ALREADY PRESENT IN THE DATABASE..


