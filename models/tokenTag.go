package models

import (
	"github.com/jmcvetta/neoism"
	"snaphyAuth/Interfaces"
)

type TokenTag struct{
	AppId int
	RealmName string
	Name string //unique among a particular realm and application.
}

//Check interface implementation..
//Will throw error if the struct doesn't implements Graph Interface..
var _ Interfaces.Graph = (*TokenTag)(nil)






func (tag *TokenTag) Exist() (isExist bool, err error)  {
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


func (tag *TokenTag) Delete() (err error){
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


