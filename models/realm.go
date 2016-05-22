package models

import (
	"github.com/jmcvetta/neoism"
	"snaphyAuth/Interfaces"
	"snaphyAuth/errorMessage"
	"snaphyAuth/helper"
)


type Realm struct{
	Name string //Must be unique among the app
	AppId int
	Id string //uuid unique identifier..
}


//Check interface implementation..
//Will throw error if the struct doesn't implements Graph Interface..
var _ Interfaces.Graph = (*Realm)(nil)


func init(){
	realm := new(Realm)
	realm.AddUniqueConstraint()
}


func (realm *Realm)AddUniqueConstraint() (err error){
	stmt := "CREATE CONSTRAINT ON (realm:Realm) ASSERT realm.id IS UNIQUE"
	cq := neoism.CypherQuery{
		Statement: stmt,
	}
	// Issue the query.
	err = db.Cypher(&cq)

	return
}


//Adding methods for nodeRealm
func (realm *Realm) Exist() (exist bool, err error){
	var appExist []struct{
		Count int `json:"count"`
	}

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


//Throw error if not exist..
func (realm *Realm) CreateIfNotExist() (err error){
	var exist bool
	//Also create relationship.
	if exist, err = realm.Exist(); err == nil && exist == false {
		id := helper.CreateUUID()
		stmt := `MATCH (app:GraphApp{id: {appId} })
			 CREATE(realm:Realm{name: {name}, appId: {appId}, id: {id} })
			 MERGE (app) - [org: Organization] -> (realm)`
		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"name": realm.Name, "appId": realm.AppId, "id": id},
		}
		// Issue the query.
		err = db.Cypher(&cq)
		if err == nil{
			//Add id
			realm.Id = id
		}
		return
	}else{
		return errorMessage.ErrorAlreadyPresent
	}
}

//Merge if exist also create relationship with app..dont create realm without relationship
func (realm *Realm) Create() (err error){

	var exist bool
	if exist, err = realm.Exist(); err == nil && exist == false {
		id := helper.CreateUUID()
		stmt := `MATCH (app:GraphApp{id: {appId} })
			 CREATE (realm:Realm{name: {name}, appId: {appId}, id: {id}})
			 MERGE (app) - [org: Organization] -> (realm)`

		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"name": realm.Name, "appId": realm.AppId, "id": id},
		}
		// Issue the query.
		err = db.Cypher(&cq)
		if err == nil{
			//Add id
			realm.Id = id
		}
	}else{
		var realmList []*Realm
		stmt := `MATCH (app:GraphApp{id: {appId} })
			 MERGE (realm:Realm{name: {name}, appId: {appId}})
			 MERGE (app) - [org: Organization] -> (realm)
			 RETURN realm.name AS name, realm.id AS id, realm.appId AS appId`

		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"name": realm.Name, "appId": realm.AppId },
			Result: &realmList,
		}
		// Issue the query.
		err = db.Cypher(&cq)

		if len(realmList) != 0{
			r := realmList[0]
			realm.Id = r.Id
		}


	}

	return
}



//Return the realm value
func (realm *Realm)Read() (err error){
	var(
		realmList []*Realm
	)

	err = realm.ReadAll(&realmList)
	if len(realmList) != 0{
		realm = realmList[0]
	}
	return
}




func (realm *Realm)ReadAll(realmListInterface [] *interface{}) (err error){
	var(
		cq neoism.CypherQuery
		realmList []*Realm
	)

	if realm.AppId != 0 && realm.Name != ""{
		stmt := `MATCH (realm: Realm) WHERE realm.name = {name} AND realm.appId = {appId}  RETURN realm.name AS name, realm.appId as appId`
		cq = neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"name": realm.Name, "appId" : realm.AppId},
			Result: &realmList,
		}
	}else if realm.AppId != 0 && realm.Name == "" {
		stmt := `MATCH (realm: Realm)  WHERE  realm.appId = {appId}  RETURN realm.name AS name, realm.appId as appId`
		cq = neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"appId" : realm.AppId},
			Result: &realmList,
		}
	}else if  realm.AppId == 0 && realm.Name != "" {
		stmt := `MATCH (realm: Realm) WHERE  realm.name = {name}  RETURN realm.name AS name, realm.appId as appId`
		cq = neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"name" : realm.Name},
			Result: &realmList,
		}
	}

	//Add value to readAll
	realmListInterface = &realmList
	// Issue the query.
	err = db.Cypher(&cq)

	return err
}



//Update realm name..
func  (realm *Realm) Update() (err error){
	var  exist bool = false
	exist, err = realm.Exist()
	if realm.Id == 0 {
		return errorMessage.ErrorIdNotPresent
	}

	if exist == false && err == nil{
		stmt := `MATCH (realm:Realm) WHERE realm.id = {id} SET realm.name = {name}`
		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{ "name": realm.Name, "id": realm.Id },
		}
		// Issue the query.
		err = db.Cypher(&cq)
	}else{
		return errorMessage.ErrorAlreadyPresent
	}

	return
}





func (realm *Realm) Delete() (err error){
	stmt := `MATCH p = (realm:Realm{name: {realmName}, appId: {appId} }) -[*]->(END) DETACH DELETE p`
	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"realmName": realm.Name, "appId": realm.AppId},
	}

	// Issue the query.
	err = db.Cypher(&cq)
	return
}







func (realm *Realm) CreateTag(tag *TokenTag) (err error){
	stmt := `MERGE (tag:Label{ name:{labelName}, appId: {appId}, realmName: {realm} })`
	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"labelName": tag.Name, "appId": realm.AppId, "realm": realm.Name},
	}
	// Issue the query.
	err = db.Cypher(&cq)
	return
}

