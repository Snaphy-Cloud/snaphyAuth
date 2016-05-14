package schemas

import (
	"github.com/graphql-go/graphql"
	snaphyInterface "snaphyAuth/Interfaces"
	"snaphyAuth/models"
)



//Schemas used in application..
var (
	ApplicationType *graphql.Object
)



func init(){
	ApplicationType = graphql.NewObject(graphql.ObjectConfig{
		Name:"Application",
		Description:"Application model contains info about user assosiated with application",
		Fields:graphql.Fields{
			"id" : &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
				Description:"Unique identity of the application.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if app, ok := p.Source.(models.Application); ok {
						return app.Id
					}

					return nil, nil
				},
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description: "Application name for graphql",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if app, ok := p.Source.(models.Application); ok {
						return app.Name
					}

					return nil, nil

				},
			},
			"status": &graphql.Field{
				Type: snaphyInterface.StatusEnum,
				Description:"status of Application",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if app, ok := p.Source.(models.Application); ok {
						return app.Status
					}

					return nil, nil
				},
			},
			"added": &graphql.Field{
				Type: graphql.String,
				Description:"DateTime when the user is added",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(models.AuthUser); ok {
						return  user.Added
					}

					return nil, nil

				},
			},
			"lastUpdated": &graphql.Field{
				Type: graphql.String,
				Description:"DateTime when the user last update their data",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(models.AuthUser); ok {
						return  user.LastUpdated
					}

					return nil, nil
				},
			},

			//TODO ADD RELAY CONNECTION FOR ADDING USER|DATABASE|TOKEN RELATIONS LATER.

		},
		Interfaces: [] *graphql.Interface{
			snaphyInterface.ApplicationInterface,
			snaphyInterface.InfoInterface,
			snaphyInterface.CreatedOnInterface,
		},
	})


}