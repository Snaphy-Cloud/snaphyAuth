package schemas

import (
	"github.com/graphql-go/graphql"
	snaphyInterface "snaphyAuth/Interfaces"
	"snaphyAuth/models"
)


var (
	AuthUserType    *graphql.Object
	ApplicationType *graphql.Object
)





func init(){
	AuthUserType = graphql.NewObject(graphql.ObjectConfig {
		Name: "AuthUser",
		Description: "Snaphy cloud main Auth type for storing all application register info.",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:graphql.NewNonNull(graphql.ID),
				Description: "id for authuser",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(models.AuthUser); ok {
						return user.Id, nil
					}
					return nil, nil
				},
			},


			"firstName": &graphql.Field{
				Type:graphql.NewNonNull(graphql.String),
				Description:"First Name of User. Required field",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(models.AuthUser); ok {
						return user.FirstName, nil
					}
					return nil, nil
				},
			},


			"lastName": &graphql.Field{
				Type: graphql.String,
				Description:"Last Name of User. Required field",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(models.AuthUser); ok {
						return user.LastName, nil
					}
					return nil, nil
				},
			},


			"email": &graphql.Field{
				Type: graphql.NewNonNull(snaphyInterface.EmailType),
				Description:"Email of User. Required field",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(models.AuthUser); ok {
						return user.Email, nil
					}
					return nil, nil
				},
			},


			"status": &graphql.Field{
				Type: snaphyInterface.StatusEnum,
				Description:"Email of User. Required field",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(models.AuthUser); ok {
						return  user.Status
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
		},
		Interfaces: [] *graphql.Interface{
			snaphyInterface.UserInterface,
			snaphyInterface.InfoInterface,
			snaphyInterface.CreatedOnInterface,
		},
	})



	ApplicationType = graphql.NewObject(graphql.ObjectConfig{
		Name:"Application",
		Description:"Application model contains info about user assosiated with application",
		Fields:graphql.Fields{
			"id" : &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
				Description:"Unique identity of the application.",
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description: "Application name for graphql",
			},
		},
	})
}