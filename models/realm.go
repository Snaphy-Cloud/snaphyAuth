package models

import (
	"github.com/jmcvetta/neoism"
	"snaphyAuth/Interfaces"
)


type Realm struct{
	Name string //Must be unique among the app
	AppId int
}


//Check interface implementation..
//Will throw error if the struct doesn't implements Graph Interface..
var _ Interfaces.Graph = (*Realm)(nil)



//Adding methods for nodeRealm
func (realm *Realm) Exist() (exist bool, err error){
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



func (realm *Realm) CreateIfNotExist() (err error){
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




func (realm *Realm) CreateGroup(group *NodeGroup) (err error)  {
	stmt := `MATCH (realm:Realm{name: {realmName}, appId: {appId} })
		 MERGE (grp:Group{name: {groupName}, appId: {appId}, realmName: {realmName} })
		 MERGE (realm) - [type: Type] -> (grp)`
	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"realmName":  realm.Name, "appId": realm.AppId, "groupName": group.Name},
	}
	// Issue the query.
	err = db.Cypher(&cq)
	return
}



func (realm *Realm) CreateTag(tag *NodeTag) (err error){
	stmt := `MERGE (tag:Label{ name:{labelName}, appId: {appId}, realmName: {realm} })`
	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"labelName": tag.Name, "appId": realm.AppId, "realm": realm.Name},
	}
	// Issue the query.
	err = db.Cypher(&cq)
	return
}

