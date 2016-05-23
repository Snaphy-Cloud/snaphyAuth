package models

import (
	"github.com/jmcvetta/neoism"
	"snaphyAuth/Interfaces"
	"errors"
	"snaphyAuth/errorMessage"
)


//Writing node model definitons..
type GraphApp struct {
	Name string
	Id int
}



//Check interface implementation..
//Will throw error if the struct doesn't implements Graph Interface..
var _ Interfaces.Graph = (*GraphApp)(nil)


func init(){
	nodeApp := new(GraphApp)
	//Adding unique constraint for name...
	nodeApp.AddUniqueConstraint()
}



//Create App in graphDb first find if any global application name is present..
func (app *GraphApp) Exist()(exist bool, err error){
	var appExist []struct{
		Count int `json:"count"`
	}

	//stmt := `MATCH (app:Application{name:{name}, id:{id}}) RETURN app.id as id, app.name as name `
	stmt := `MATCH (app: GraphApp) WHERE app.name = {name}  RETURN count(app) as count `
	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"name": app.Name},
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


func (app *GraphApp)ReadAll(appList [] *interface{}) (err error){
	var(
		cq neoism.CypherQuery
		graphApp []*GraphApp
	)

	if app.Id != 0 && app.Name != "" {
		stmt := `MATCH (app: GraphApp) WHERE app.name = {name} AND app.id = {id}  RETURN app.name AS name, app.id as id LIMIT 1000`
		cq = neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"name": app.Name, "id" : app.Id},
			Result: &graphApp,
		}
	}else if app.Id != 0 && app.Name == "" {
		stmt := `MATCH (app: GraphApp) WHERE  app.id = {id}  RETURN app.name AS name, app.id as id LIMIT 1000`
		cq = neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"id" : app.Id},
			Result: &graphApp,
		}
	}else if  app.Id == 0 && app.Name != "" {
		stmt := `MATCH (app: GraphApp) WHERE  app.name = {name}  RETURN app.name AS name, app.id as id LIMIT 1000`
		cq = neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"name" : app.Name},
			Result: &graphApp,
		}
	}

	//Add value to readAll
	appList = &graphApp

	// Issue the query.
	err = db.Cypher(&cq)

	return err

}


//Read data and populate..app
func (app *GraphApp)Read() (err error){
	var graphApp []*GraphApp
	err = app.ReadAll(&graphApp)
	if len(graphApp) != 0{
		app =  graphApp[0]
	}
	return
}


//Check if app is present of the same name and id..
func (app *GraphApp)Get()(exist bool, err error){
	var appExist []struct{
		Count int `json:"count"`
	}

	stmt := `MATCH (app: GraphApp) WHERE app.name = {name} AND app.id = {id}  RETURN count(app) as count `
	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"name": app.Name, "id" : app.Id},
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


func (app *GraphApp)AddUniqueConstraint() (err error){
	stmt := "CREATE CONSTRAINT ON (app:GraphApp) ASSERT app.name IS UNIQUE"
	cq := neoism.CypherQuery{
		Statement: stmt,
	}
	// Issue the query.
	err = db.Cypher(&cq)

	return
}


//Will merge if node already exists.
func (app *GraphApp) Create() (err error){
	stmt := `MERGE (app:GraphApp{name: {name}, id: {id} })`
	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"name": app.Name, "id": app.Id},
	}

	// Issue the query.
	err = db.Cypher(&cq)
	return
}



//Will return an error if the node already exists..
func (app *GraphApp) CreateIfNotExist() (err error){
	var exist bool
	if exist, err = app.Exist(); err == nil && exist == false {
		stmt := `CREATE(app:GraphApp{name: {name}, id: {id} })`
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




//DELETE ALL DATA INCLUDING ITS NODE TO..
//TODO WARNING WHEN PREPARING GRAPHQL END POINT ADD A WEEK TIME AFTER WHICH THIS DATABASE WILL GET DELETED.
//TODO DEACTIVATE APP WITHIN THIS TIME AND PERMANENTLY DELETE AFTER ONE WEEK
//TODO ALSO SEND AN EMAIL WARNING DELETING OF DATA.
func (app *GraphApp)Delete()(err error){
	stmt :=  `MATCH q = (app:GraphApp) WHERE app.id = {id} AND app.name = {name} OPTIONAL MATCH p = (app)-[*]-() DETACH DELETE p, q`
	if app.Name != "" && app.Id != 0{
		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"name": app.Name, "id": app.Id},
		}

		// Issue the query.
		err = db.Cypher(&cq)
		return
	}else{
		return errors.New(errorMessage.GRAPH_APP_NAME_NOT_FOUND)
	}
}


//Only name can be updated..throw error if already exists..
func (app *GraphApp) Update(err error){
	var  exist bool = false
	if exist , err = app.Exist(); exist == false && err == nil{
		stmt := `MATCH (app: GraphApp { id: {id} })
			 SET app.name = {name} `
		cq := neoism.CypherQuery{
			Statement: stmt,
			Parameters: neoism.Props{"name": app.Name, "id": app.Id},
		}

		// Issue the query.
		err = db.Cypher(&cq)

	}else{
		return errors.New(errorMessage.GRAPH_APP_ALREADY_PRESENT)
	}

	return
}



