package models

import (
	"github.com/jmcvetta/neoism"
	"errors"
	"fmt"
	"strconv"
)








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

	nodeRealm := new(Realm)
	nodeRealm.Name = "snaphyTest"

	nodeRealm.AppId = nodeApp.Id

	//Add realm.*:
	err = nodeApp.CreateRealm(nodeRealm)
	if err == nil{
		nodeGroup := new(Group)
		nodeGroup.AppId = nodeApp.Id
		nodeGroup.Name = "customer"
		err = nodeRealm.CreateGroup(nodeGroup)
		if err != nil{
			fmt.Println(err)
		}else{
			fmt.Println("Successfully created Group")
			tag :=  new (TokenTag)
			tag.RealmName = nodeRealm.Name
			tag.AppId = nodeApp.Id
			tag.Name = "Admin"
			err = nodeRealm.CreateTag(tag)
			if err != nil{
				fmt.Println(err)
			}else{

				nodeToken := new(Token)
				//app *Application, settings *ApplicationSettings, tokenHelper TokenHelper, realm *Realm, tag *NodeTag, userIdentity string
				app := new(Application)
				settings := new(ApplicationSettings)
				settings.Application = app
				settings.ExpiryDuration = 3600
				user := new(User)
				user.Id = 3
				user.GetUser()
				app.Id = 3
				app.GetApp()
				app.FetchAppTokens()
				if len(app.TokenInfo) != 0{
					//CreateToken(token *NodeToken, app *Application, settings *ApplicationSettings, tokenHelper TokenHelper, realm *Realm, tag *NodeTag, userIdentity string)
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




//Create realm with relationship..if not exists..first check if it already exists for showing custom error
func (app *NodeApp) CreateRealm(realm *Realm)(err error){
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







