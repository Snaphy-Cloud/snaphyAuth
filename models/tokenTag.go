package models

import (
	"github.com/jmcvetta/neoism"
	"snaphyAuth/Interfaces"
	"golang.org/x/text/internal/tag"
	"snaphyAuth/helper"
	"snaphyAuth/errorMessage"
)


type TokenTag struct{
	Id string //UUID unique identifier..
	AppId int
	RealmName string
	Name string //unique among a particular realm and application.
}



//Check interface implementation..
//Will throw error if the struct doesn't implements Graph Interface..
var _ Interfaces.Graph = (*TokenTag)(nil)



func (tokenTag *TokenTag)AddUniqueConstraint() (err error){
	stmt := "CREATE CONSTRAINT ON (tag:TokenTag) ASSERT tag.id IS UNIQUE"
	cq := neoism.CypherQuery{
		Statement: stmt,
	}

	// Issue the query.
	err = db.Cypher(&cq)
	return
}


func (tokenTag *TokenTag) Create() (err error){
	err = tokenTag.CreateIfNotExist()
	if err != nil{
		var tokenTagList []*TokenTag
		//Node exists just merge.
		stmt := `MATCH (realm:Realm{name: {realmName}, appId: {appId} })
			 MERGE (tokenTag:TokenTag{name: {tokenName}, appId: {appId}, realmName: {realmName}})
			 MERGE (realm) - [type:Tag] -> (tokenTag)
			 RETURN tokenTag.name AS name, tokenTag.id AS id, tokenTag.appId as appId, tokenTag.realmName AS realmName`

		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"realmName":  tokenTag.RealmName, "appId": tokenTag.AppId, "tokenName": tokenTag.Name},
			Result: &tokenTagList,
		}

		// Issue the query.
		err = db.Cypher(&cq)

		if len(tokenTagList) != 0 {
			g := tokenTagList[0]
			tokenTag.Id = g.Id
		}

		return
	}

	return
}


func (tokenTag *TokenTag)CreateIfNotExist() (err error){
	var exist bool
	//Also create relationship.
	if exist, err = tokenTag.Exist(); err == nil && exist == false {
		id := helper.CreateUUID()
		stmt := `MATCH (realm:Realm{name: {realmName}, appId: {appId} })
			 CREATE (tokenTag:TokenTag{name: {tokenName}, appId: {appId}, realmName: {realmName}, id: {id} })
			 CREATE (realm) - [type:Tag] -> (tokenTag)`
		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"realmName":  tokenTag.RealmName, "appId": tokenTag.AppId, "tokenName": tokenTag.Name, "id": id},
		}

		// Issue the query.
		err = db.Cypher(&cq)

		if err == nil{
			//Add id
			tokenTag.Id = id
		}
		return
	}else{
		return errorMessage.ErrorAlreadyPresent
	}
}




func (tokenTag *TokenTag) Exist() (exist bool, err error)  {
	var tagExist []struct{
		Count int `json:"count"`
	}
	stmt := `MATCH (tag: TokenTag{ name:{tokenName}, appId: {appId}, realmName: {realmName} }) RETURN count(tag) AS count`

	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"tokenName": tokenTag.Name, "appId": tokenTag.AppId, "realmName": tokenTag.RealmName},
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



func (tokenTag *TokenTag) Delete() (err error){
	stmt := `MATCH (tag:TokenTag{ name:{tokenName}, appId: {appId}, realmName: {realm} })
	         DETACH DELETE tag`

	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"tokenName": tokenTag.Name, "appId": tokenTag.AppId, "realm": tokenTag.RealmName},
	}

	// Issue the query.
	err = db.Cypher(&cq)
	return
}



func (tokenTag *TokenTag) Update() (err error){
	var  exist bool = false
	exist, err = tokenTag.Exist()
	if tokenTag.Id == 0 {
		return errorMessage.ErrorIdNotPresent
	}

	if exist == false && err == nil{
		stmt := `MATCH (tokenTag:Group) WHERE tokenTag.id = {id} SET tokenTag.name = {name}`
		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{ "name": tokenTag.Name, "id": tokenTag.Id },
		}
		// Issue the query.
		err = db.Cypher(&cq)
	}else{
		return errorMessage.ErrorAlreadyPresent
	}

	return
}



func (tokenTag *TokenTag) ReadAll(tokenTagListInterface [] *interface{}) (err error){
	var tokenTagList []*TokenTag

	if(tokenTag.Id != ""){
		//Node exists just merge.
		stmt := `MATCH (tokenTag:TokenTag) WHERE tokenTag.id = {id}
			 RETURN tokenTag.name AS name, tokenTag.id AS id, tokenTag.appId as appId, tokenTag.realmName AS realmName`


		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"id": tokenTag.Id},
			Result: &tokenTagList,
		}

		// Issue the query.
		err = db.Cypher(&cq)

	}else if (tokenTag.AppId != "" && tokenTag.RealmName != "" && tokenTag.Name == ""){
		//Node exists just merge.
		stmt := `MATCH (tokenTag:TokenTag) WHERE tokenTag.appId = {appId}  AND tokenTag.realmName = {realmName}
			 RETURN tokenTag.name AS name, tokenTag.id AS id, tokenTag.appId as appId, tokenTag.realmName AS realmName`


		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"appId": tokenTag.AppId,  "realmName": tokenTag.RealmName},
			Result: &tokenTagList,
		}

		// Issue the query.
		err = db.Cypher(&cq)
	}else{
		//Node exists just merge.//Else name, realmName, appId is compulsary..
		stmt := `MATCH (tokenTag:TokenTag) WHERE tokenTag.appId = {appId} AND tokenTag.name = {tokenTag.name} AND tokenTag.realmName = {realmName}
			 RETURN tokenTag.name AS name, tokenTag.id AS id, tokenTag.appId as appId, tokenTag.realmName AS realmName`


		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"appId": tokenTag.AppId, "name": tokenTag.Name, "realmName": tokenTag.RealmName},
			Result: &tokenTagList,
		}

		// Issue the query.
		err = db.Cypher(&cq)
	}


	if len(tokenTagList) != 0 {
		tokenTagListInterface = &tokenTagList
	}
	return
}


