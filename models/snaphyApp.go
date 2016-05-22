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
	nodeApp := new(GraphApp)
	nodeApp.Id = 3
	nodeApp.Name = "snaphyAdminAuth"

	//Create app..
	nodeApp.CreateIfNotExist()
	//nodeApp.DeleteApp()

	nodeRealm := new(Realm)
	nodeRealm.Name = "snaphyTest"

	nodeRealm.AppId = nodeApp.Id

	//Add realm.*:
	err = nodeRealm.Create()
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







