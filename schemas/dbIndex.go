package schemas


import (
	"github.com/graphql-go/graphql"
	snaphyInterface "snaphyAuth/Interfaces"
	"snaphyAuth/models"
)



//Schemas used in application..
var (
	DbIndexType    *graphql.Object
)


func init(){
	DbIndexType = graphql.NewObject(graphql.ObjectConfig{
		Name:"DbIndex",
		Description:"Database index mapping a database.",
		Fields: graphql.Fields{
			"id" : &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
				Description:"Unique identity of database index lookup table.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if db, ok := p.Source.(models.DbIndex); ok {
						return db.Id
					}
					return nil, nil
				},
			},
			"name" : &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description:"Name of the database.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if db, ok := p.Source.(models.DbIndex); ok {
						return db.Name
					}
					return nil, nil
				},
			},
			"name" : &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description:"Name of the database.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if db, ok := p.Source.(models.DbIndex); ok {
						return db.Name
					}
					return nil, nil
				},
			},
			"dbType": &graphql.Field{
				Type: graphql.NewNonNull(snaphyInterface.DatabaseTypeEnum),
				Description: "Type of database used. postgres|mysql|mongodb",
				Resolve: func(p graphql.ResolveParams) (interface{}, error){
					if db, ok := p.Source.(models.DbIndex); ok{
						return db.DbType
					}
					return nil, nil
				},
			},

			"dbPass": &graphql.Field{
				Type: graphql.NewNonNull(snaphyInterface.DatabaseTypeEnum),
				Description: "Password used for database",
				Resolve: func(p graphql.ResolveParams) (interface{}, error){
					if db, ok := p.Source.(models.DbIndex); ok{
						return db.DbPass
					}
					return nil, nil
				},
			},

			"dbUser": &graphql.Field{
				Type: graphql.NewNonNull(snaphyInterface.DatabaseTypeEnum),
				Description: "Username authorized for this database",
				Resolve: func(p graphql.ResolveParams) (interface{}, error){
					if db, ok := p.Source.(models.DbIndex); ok{
						return db.DbUser
					}
					return nil, nil
				},
			},
			"status": &graphql.Field{
				Type: snaphyInterface.StatusEnum,
				Description:"status of Application",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if db, ok := p.Source.(models.DbIndex); ok {
						return db.Status
					}

					return nil, nil
				},
			},
			"added": &graphql.Field{
				Type: graphql.String,
				Description:"DateTime when the user is added",
				Resolve:func(p graphql.ResolveParams) (interface{}, error){
					if db, ok := p.Source.(models.DbIndex); ok {
						return db.Added
					}
					return nil, nil
				},
			},
			"lastUpdated": &graphql.Field{
				Type: graphql.String,
				Description:"DateTime when the user last update their data",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if db, ok := p.Source.(models.DbIndex); ok {
						return  db.LastUpdated
					}
					return nil, nil
				},
			},
		},
		Interfaces: [] *graphql.Interface{
			snaphyInterface.DbIndexInterface,
			snaphyInterface.InfoInterface,
			snaphyInterface.CreatedOnInterface,
		},
	})

}
