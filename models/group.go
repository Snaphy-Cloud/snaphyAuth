package models

import (
	"github.com/jmcvetta/neoism"
	"snaphyAuth/Interfaces"
	"snaphyAuth/helper"
	"snaphyAuth/errorMessage"
)



type Group struct {
	Id string //uuid unique identifier..
	Name string
	AppId int
	RealmName string
}





//Check interface implementation..
//Will throw error if the struct doesn't implements Graph Interface..
var _ Interfaces.Graph = (*Group)(nil)



func init(){
	group := new(Group)
	group.AddUniqueConstraint()
}



func (group *Group)AddUniqueConstraint() (err error){
	stmt := "CREATE CONSTRAINT ON (group:Group) ASSERT group.id IS UNIQUE"
	cq := neoism.CypherQuery{
		Statement: stmt,
	}
	// Issue the query.
	err = db.Cypher(&cq)

	return
}


//Adding methods for nodeRealm
func (group *Group) Exist() (exist bool, err error){
	var groupExist []struct{
		Count int `json:"count"`
	}

	stmt := `MATCH (group:Group) WHERE group.name = {name} AND group.appId = {appId} AND group.realmName = {realmName} RETURN count(app) AS count`
	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"name": group.Name, "appId": group.AppId, "realmName": group.RealmName},
		Result: &groupExist,
	}

	// Issue the query.
	err = db.Cypher(&cq)
	if err != nil{
		return false, err
	}

	if len(groupExist) != 0{
		if groupExist[0].Count == 0{
			return false, err
		}else{
			return true, err
		}
	}else{
		return false, err
	}
}


//Throw error if not exist..
func (group *Group) CreateIfNotExist() (err error){
	var exist bool
	//Also create relationship.
	if exist, err = group.Exist(); err == nil && exist == false {
		id := helper.CreateUUID()
		stmt := `MATCH (realm:Realm{name: {realmName}, appId: {appId} })
			 CREATE (grp:Group{name: {groupName}, appId: {appId}, realmName: {realmName}, id: {id} })
			 CREATE (realm) - [type: Type] -> (grp)`
		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"realmName":  group.RealmName, "appId": group.AppId, "groupName": group.Name, "id": id},
		}

		// Issue the query.
		err = db.Cypher(&cq)

		if err == nil{
			//Add id
			group.Id = id
		}
		return
	}else{
		return errorMessage.ErrorAlreadyPresent
	}
}



//Create and Merge if exists..
func (group *Group) Create() (err error)  {
	err = group.CreateIfNotExist()
	if err != nil{
		var grpList []*Group
		//Node exists just merge.
		stmt := `MATCH (realm:Realm{name: {realmName}, appId: {appId} })
			 MERGE (grp:Group{name: {groupName}, appId: {appId}, realmName: {realmName}})
			 MERGE (realm) - [type: Type] -> (grp)
			 RETURN grp.name AS name, grp.id AS id, grp.appId as appId, grp.realmName AS realmName`
		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"realmName":  group.RealmName, "appId": group.AppId, "groupName": group.Name},
			Result: &grpList,
		}

		// Issue the query.
		err = db.Cypher(&cq)

		if len(grpList) != 0 {
			g := grpList[0]
			group.Id = g.Id
		}

		return
	}

	return
}


func (group *Group) Delete() (err error){
	stmt := `MATCH p = (begin:Group{name: {groupName}, appId: {appId}, realmName: {realmName} })-[r*]->(END:Token)  DETACH DELETE p`
	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"groupName": group.Name, "appId": group.AppId, "realmName": group.RealmName},
	}
	// Issue the query.
	err = db.Cypher(&cq)
	return
}



func (group *Group) ReadAll(groupListInterface [] *interface{}) (err error){
	var grpList []*Group

	if(group.Id != ""){
		//Node exists just merge.
		stmt := `MATCH (grp:Group) WHERE grp.id = {id}
			 RETURN grp.name AS name, grp.id AS id, grp.appId as appId, grp.realmName AS realmName LIMIT 1000`


		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"id": group.Id},
			Result: &grpList,
		}

		// Issue the query.
		err = db.Cypher(&cq)

	}else{
		//Node exists just merge.
		stmt := `MATCH (grp:Group) WHERE grp.appId = {appId} AND grp.name = {grp.name} AND grp.realmName = {realmName}
			 RETURN grp.name AS name, grp.id AS id, grp.appId as appId, grp.realmName AS realmName LIMIT 1000`


		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"appId": group.AppId, "name": group.Name, "realmName": group.RealmName},
			Result: &grpList,
		}

		// Issue the query.
		err = db.Cypher(&cq)
	}


	if len(grpList) != 0 {
		groupListInterface = &grpList
	}
	return
}


func  (group *Group) Read()  (err error){
	var(
		grpList []*Group
	)

	err = group.ReadAll(&grpList)
	if len(grpList) != 0{
		group = grpList[0]
	}
	return
}


func (group *Group) Update() (err error){
	var  exist bool = false
	exist, err = group.Exist()
	if group.Id == 0 {
		return errorMessage.ErrorIdNotPresent
	}

	if exist == false && err == nil{
		stmt := `MATCH (group:Group) WHERE group.id = {id} SET group.name = {name}`
		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{ "name": group.Name, "id": group.Id },
		}
		// Issue the query.
		err = db.Cypher(&cq)
	}else{
		return errorMessage.ErrorAlreadyPresent
	}

	return
}
