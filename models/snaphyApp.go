package models

//Writing node model definitons..
type NodeApp struct {
	Id   int
	Name string //Must be an unique name
}


type NodeRealm struct{
	Name string //Must be unique among the app
	AppId int
}


type NodeGroup struct {
	Name string
	AppId int
	Realm string
}


type NodeToken struct{
	TokenId string //JWT TOKEN INFO
	AppId int
	Realm string
	Status string
	added string
	lastUpdated string

}

type NodeTag struct{
	AppId int
	Realm string
	Name string //unique among a particular realm and application.
}


//Writing relationship struct..
type  RelIdentity struct{
	userId string
}


