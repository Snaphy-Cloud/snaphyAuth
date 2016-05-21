package models

import (
	"github.com/jmcvetta/neoism"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"time"
	"fmt"
	"strconv"
	"strings"

)

//Writing node model definitons..
type NodeApp Application




type NodeRealm struct{
	Name string //Must be unique among the app
	AppId int
}


type NodeGroup struct {
	Name string
	AppId int
	RealmName string
}


type NodeToken struct{
	TokenString string //JWT TOKEN INFO
	AppId int
	RealmName string
	GroupName string
	Status string
	UserId int
	Added int64
	Expiry int64
	LastUpdated int64
	JTI string //Unique token provider.
	//TODO LATER CONVERT IT TO ARRAY TYPE
	Roles string
}

type NodeTag struct{
	AppId int
	RealmName string
	Name string //unique among a particular realm and application.
}


//Writing relationship struct..
type  RelIdentity struct{
	userId string
}


type node struct {
	N neoism.Node // Column "n" gets automagically unmarshalled into field N
}


var (
	db *neoism.Database
	ErrorAlreadyPresent = errors.New("Error Node already present")
)


func init() {
	var err error
	db, err = neoism.Connect("http://neo4j:myfunzone2030@localhost:7474/db/data")
	if err != nil {
		panic(err)
	}
	//Remove this after test..
	test()

}



func test(){
	var err error
	nodeApp := new(NodeApp)
	nodeApp.Id = 3
	nodeApp.Name = "snaphyAdminAuth"
	//Adding unique constraint for name...
	nodeApp.AddUniqueConstraint()
	//Create app..
	nodeApp.CreateIfNotExist()
	nodeApp.Status = StatusMap["ACTIVE"]
	nodeApp.UpdateStatus()
	//nodeApp.DeleteApp()

	nodeRealm := new(NodeRealm)
	nodeRealm.Name = "snaphyTest"

	nodeRealm.AppId = nodeApp.Id

	//Add realm.*:
	err = nodeApp.CreateRealm(nodeRealm)
	if err == nil{
		nodeGroup := new(NodeGroup)
		nodeGroup.AppId = nodeApp.Id
		nodeGroup.Name = "customer"
		err = nodeRealm.CreateGroup(nodeGroup)
		if err != nil{
			fmt.Println(err)
		}else{
			fmt.Println("Successfully created Group")
			tag :=  new (NodeTag)
			tag.RealmName = nodeRealm.Name
			tag.AppId = nodeApp.Id
			tag.Name = "Admin"
			err = nodeRealm.CreateTag(tag)
			if err != nil{
				fmt.Println(err)
			}else{

				nodeToken := new(NodeToken)
				//app *Application, settings *ApplicationSettings, tokenHelper TokenHelper, realm *NodeRealm, tag *NodeTag, userIdentity string
				app := new(Application)
				settings := new(ApplicationSettings)
				settings.Application = app
				settings.ExpiryDuration = 3600
				user := new(AuthUser)
				user.Id = 3
				user.GetUser()
				app.Id = 3
				app.GetApp()
				app.FetchAppTokens()
				if len(app.TokenInfo) != 0{
					//CreateToken(token *NodeToken, app *Application, settings *ApplicationSettings, tokenHelper TokenHelper, realm *NodeRealm, tag *NodeTag, userIdentity string)
					err = nodeGroup.CreateToken(nodeToken, app, settings, app.TokenInfo[0], nodeRealm, tag, strconv.Itoa(user.Id) )
					if err != nil{
						fmt.Println(err)
					}else{

						//(nodeToken *NodeToken) Verify() (valid bool, err error)
						nodeToken.VerifyHash(app.TokenInfo[0])
						fmt.Println("\n\n")
						nodeToken.VerifyAndParse(app.TokenInfo[0])
						expired, _ := nodeToken.CheckExpired()
						fmt.Println(expired)
					}
				}else{
					fmt.Println("Error: TokenInfo not present in helper file. add helper data first..")
				}
			}
		}
	}else{
		fmt.Println(err)
	}
}


//Create App in graphDb first find if any global application is present..
func (app *NodeApp)Exist()(exist bool, err error){
	var appExist []struct{
		Count int `json:"count"`
	}


	//stmt := `MATCH (app:Application{name:{name}, id:{id}}) RETURN app.id as id, app.name as name `
	stmt := `MATCH (app:Application) WHERE app.name = {name} AND app.id = {id} RETURN count(app) as count `
	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"name": app.Name, "id": app.Id},
		Result: &appExist,
	}

	// Issue the query.
	err = db.Cypher(&cq)

	if err != nil{
		return false, err
	}

	if len(appExist) != 0{
		if appExist[0].Count == 0{
			return false, err
		}else{
			return true, err
		}
	}else{
		return false, err
	}

}


func (app *NodeApp)AddUniqueConstraint() (err error){
	stmt := "CREATE CONSTRAINT ON (app:Application) ASSERT app.name IS UNIQUE"
	cq := neoism.CypherQuery{
		Statement: stmt,
	}
	// Issue the query.
	err = db.Cypher(&cq)

	return
}


func (app *NodeApp) CreateIfNotExist() (err error){
	var exist bool
	if exist, err = app.Exist(); err == nil && exist == false {
		stmt := `CREATE(app:Application{name: {name}, id: {id} })`
		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"name": app.Name, "id": app.Id},
		}

		// Issue the query.
		err = db.Cypher(&cq)
		return
	}else{
		return ErrorAlreadyPresent
	}
}


func (app *NodeApp) UpdateStatus() (err error){
	stmt :=  `MATCH (app:Application) WHERE app.name = {name} AND app.id = {id} SET app.status = {status} `
	if app.Status != ""{
		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"name": app.Name, "id": app.Id, "status": app.Status},
		}

		// Issue the query.
		err = db.Cypher(&cq)
		return
	}else {
		return errors.New("Application Status cannot be empty")
	}
}


func (app *NodeApp)DeleteApp()(err error){
	stmt :=  `MATCH q = (app:Application{id:1}) OPTIONAL MATCH p = (app)-[*]-() DETACH DELETE p, q`
	if app.Name != "" && app.Id != 0{
		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"name": app.Name, "id": app.Id},
		}

		// Issue the query.
		err = db.Cypher(&cq)
		return
	}else{
		return errors.New("Application ID and Name cannot be empty")
	}

}


func (app *NodeApp)DeactivateApp() (err error){
	app.Status = StatusMap["DEACTIVATED"]
	err = app.UpdateStatus()
	return
}


//Adding methods for nodeRealm
func (realm *NodeRealm) Exist() (exist bool, err error){
	var appExist []struct{
		Count int `json:"count"`
	}


	//stmt := `MATCH (app:Application{name:{name}, id:{id}}) RETURN app.id as id, app.name as name `
	stmt := `MATCH (realm:Realm) WHERE realm.name = {name} AND realm.appId = {appId} RETURN count(app) as count`
	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"name": realm.Name, "appId": realm.AppId},
		Result: &appExist,
	}

	// Issue the query.
	err = db.Cypher(&cq)
	if err != nil{
		return false, err
	}

	if len(appExist) != 0{
		if appExist[0].Count == 0{
			return false, err
		}else{
			return true, err
		}
	}else{
		return false, err
	}
}


//Create realm with relationship..if not exists..first check if it already exists for showing custom error
func (app *NodeApp) CreateRealm(realm *NodeRealm)(err error){
	stmt := `MATCH (app:Application{name: {appName}, id: {appId} })
		 MERGE(realm:Realm{name: {realmName}, appId: {appId} })
		 MERGE (app) - [org: Organization] -> (realm)`
	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"realmName": realm.Name, "appId": app.Id, "appName": app.Name},
	}

	// Issue the query.
	err = db.Cypher(&cq)
	return
}




func (realm *NodeRealm) CreateIfNotExist() (err error){
	var exist bool
	if exist, err = realm.Exist(); err == nil && exist == false {
		stmt := `CREATE(realm:Realm{name: {name}, appId: {id} })`
		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"name": realm.Name, "appId": realm.AppId},
		}
		// Issue the query.
		err = db.Cypher(&cq)
		return
	}else{
		return ErrorAlreadyPresent
	}
}



func (realm *NodeRealm) CreateGroup(group *NodeGroup) (err error)  {
	stmt := `MATCH (realm:Realm{name: {realmName}, appId: {appId} })
		 MERGE (grp:Group{name: {groupName}, appId: {appId}, realmName: {realmName} })
		 MERGE (realm) - [type: Type] -> (grp)`

	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"realmName": realm.Name, "appId": realm.AppId, "groupName": group.Name},
	}

	// Issue the query.
	err = db.Cypher(&cq)
	return
}



func (realm *NodeRealm) CreateTag(tag *NodeTag) (err error){
	stmt := `MERGE (tag:Label{ name:{labelName}, appId: {appId}, realmName: {realm} })`
	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"labelName": tag.Name, "appId": realm.AppId, "realm": realm.Name},
	}
	// Issue the query.
	err = db.Cypher(&cq)
	return
}



func (tag *NodeTag) Exist() (isExist bool, err error)  {
	var tagExist []struct{
		Count int `json:"count"`
	}
	stmt := `MATCH (tag: Label{ name:{labelName}, appId: {appId}, realmName: {realm} }) RETURN count(tag) AS count`

	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"labelName": tag.Name, "appId": tag.AppId, "realm": tag.RealmName},
		Result: &tagExist,
	}

	// Issue the query.
	err = db.Cypher(&cq)
	if err != nil{
		return false, err
	}

	if len(tagExist) != 0{
		if tagExist[0].Count == 0{
			return false, err
		}else{
			return true, err
		}
	}else{
		return false, err
	}
}


func (tag *NodeTag) Delete() (err error){
	stmt := `MATCH (tag:Label{ name:{name}, appId: {appId}, realmName: {realm} })
	         OPTIONAL MATCH (tag)- [role:Role] -> ()
	         DETACH DELETE tag, role`

	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"name": tag.Name, "appId": tag.AppId, "realm": tag.RealmName},
	}

	// Issue the query.
	err = db.Cypher(&cq)
	return
}



func (group *NodeGroup) Delete() (err error){
	stmt := `MATCH p =(begin:Group{name: {groupName}, appId: {appId}, realmName: {realmName} })-[r*]->(END:Token)  DETACH DELETE p`
	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"groupName": group.Name, "appId": group.AppId, "realmName": group.RealmName},
	}

	// Issue the query.
	err = db.Cypher(&cq)
	return
}


func (realm *NodeRealm) Delete() (err error){
	stmt := `MATCH p = (realm:Realm{name: {realmName}, appId: {appId} }) -[*]->(END) DETACH DELETE p`
	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"realmName": realm.Name, "appId": realm.AppId},
	}

	// Issue the query.
	err = db.Cypher(&cq)
	return
}



//Add a tag with relationship..
func (token *NodeToken) AddTag(tag *NodeTag, userIdentity string) (err error){
	stmt := `MATCH (tag: Label{ name:{labelName}, appId: {appId}, realmName: {realm} })
	         MATCH (token:Token) WHERE token.name = {tokenString} AND token.appId = {appId} AND token.realmName = {realm}
	         MERGE (tag) - [role:Role] -> (token)`

	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"labelName": tag.Name, "appId": token.AppId, "realm": token.RealmName, "tokenString": token.TokenString, "userId": userIdentity },
	}

	// Issue the query.
	err = db.Cypher(&cq)
	return
}



//Verify the token before parsing .. the token just checks if the tokens ia a valid one.
func (nodeToken *NodeToken) VerifyHash(tokenHelper *TokenHelper) (valid bool, err error){

	parts := strings.Split(nodeToken.TokenString, ".")

	method := jwt.GetSigningMethod(tokenHelper.HashType)
	err = method.Verify(strings.Join(parts[0:2], "."), parts[2], []byte(tokenHelper.PublicKey))

	if err != nil{

		return false, err
	}else{
		return true, err
	}

}

//Parses the token value..And also validates the algorithm..
func (nodeToken *NodeToken) VerifyAndParse(tokenHelper *TokenHelper) (valid bool, err error){
	var token *jwt.Token
	token, err = jwt.Parse(nodeToken.TokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenHelper.PublicKey), nil
	})

	nodeToken.Added = int64(token.Claims["iat"].(float64))
	nodeToken.Expiry = int64(token.Claims["exp"].(float64))
	nodeToken.JTI = token.Claims["jti"].(string)
	nodeToken.GroupName = token.Claims["grp"].(string)
	nodeToken.Roles = token.Claims["roles"].(string)
	nodeToken.RealmName = token.Claims["realm"].(string)
	return token.Valid, err
}

//Check if token is expired or not valid or valid...
//TRUE IF EXPIRED AND FALSE IF NOT
func (nodeToken *NodeToken) CheckExpired() (expired bool, err error){
	if nodeToken.Expiry == 0 {
		return true, errors.New("Expiry claim not present in NodeToken model.")
	}else{
		now := time.Now().Unix()
		//TODO ALSO CHECK THE STATUS IN GRAPHDATABASE FIRST
		//TODO CREATE METHOD GETSTATUS() in NODE TOKEN
		if now >= nodeToken.Expiry{
			return true, nil
		}else{
			return false, nil
		}
	}

}


//TODO WRITE A METHOD TO UPDATE STATUS IF FOUND INVALID..TOKEN
//TODO WRITE A METHOD TO FIND IF THE TOKEN ALREADY PRESENT IN THE DATABASE..






//Create and sign the token..
func (group *NodeGroup) CreateToken(token *NodeToken, app *Application, settings *ApplicationSettings, tokenHelper *TokenHelper, realm *NodeRealm, tag *NodeTag, userIdentity string) (err error){
	// Create the token
	var signToken *jwt.Token

	if tokenHelper.HashType == "" {
		return errors.New("No Signing method present in the TokenHelper Database")
	}else if(tokenHelper.HashType == "RS256"){
		signToken = jwt.New(jwt.SigningMethodRS256)
	}else{
		//Right now default is RS256
		signToken = jwt.New(jwt.SigningMethodRS256)
	}

	//signToken := jwt.GetSigningMethod(tokenHelper.HashType)

	duration := settings.ExpiryDuration
	if duration == 0 {
		duration = 3600
	}

	expiryTime := time.Now().Add(time.Second * time.Duration(duration)).Unix() //time after it will be invalid..
	issuedAt := time.Now().Unix() //Issuer at time..
	jti := uuid.NewV4().String() //unique number provider
	signToken.Claims["exp"] = expiryTime
	signToken.Claims["iat"] = issuedAt
	signToken.Claims["iss"] = userIdentity //user identity..
	signToken.Claims["grp"] = group.Name //Group with which user belong to.
	signToken.Claims["realm"] = realm.Name
	//TODO LATER ADD MORE ROLES BY FETCHING RELATION FROM GRAPH DATABASE..
	signToken.Claims["roles"] = tag.Name
	//JUST STORE THIS INSTEAD OF STORING TOKENS..
	signToken.Claims["jti"] = jti

	if tokenHelper.PrivateKey != ""{

		// Sign and get the complete encoded token as a string
		tokenString, err := signToken.SignedString([]byte(tokenHelper.PrivateKey))
		if err != nil{
			return err
		}else{
			userId, _ := strconv.Atoi(userIdentity)
			token.Added = issuedAt
			token.LastUpdated = issuedAt
			token.AppId = app.Id
			token.RealmName = realm.Name
			token.Status = StatusMap["ACTIVE"]
			token.TokenString = tokenString
			token.UserId = userId
			token.JTI = jti
		}

		return err
	}else{
		return errors.New("Private key not present in tokenHelper")
	}
}








//TODO ADD CREATE TOKEN METHOD AND ATTACH TO GRAPH
//TODO ARRANGE ALL MODELS TO DIFFERENT FOLDER
//TODO WRITE TEST CASES FOR METHODS
//TODO WRITE GRAPHQL METHOD FOR FOR DATA ENDPOINT
//TODO IN TOKEN FIRST ENCRYPT THE DATA WITH PRIVATE KEY and ADD USERID TO RELATION
//TODO ADD A METHOD IN TOKEN TO REFRESH WITH AUTOMATIC DELETE FOR BAD TOKENS
//TODO ADD A METHOD TO CHECK LOGIN
//TODO ADD METHOD TO DELETE GROUP AND REALM AND TOKENS AND TAG
//TODO ADD A METHOD TO LOGOUT FROM ALL TOKEN OF A USER.
//TODO ADD A METHOD TO SHOW DIFFERENT CURRENT LOGIN TOKENS







