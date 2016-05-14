package schemas

import (
	"github.com/graphql-go/graphql"
	snaphyInterface "snaphyAuth/Interfaces"
	"snaphyAuth/models"
)


var (
	AuthUserType    *graphql.Object
	ApplicationType *graphql.Object
	TokenType       *graphql.Object
)





func init(){
	//Defigning some fields now..
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
				Description:"status of User. Required field",
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

			//TODO ADD RELAY CONNECTION FOR ADDING APPLICATIONS RELATIONS LATER.
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

	TokenType = graphql.NewObject(graphql.ObjectConfig{
		Name:"Token",
		Description:"Token containing main application token for the user for access management.",
		Fields:graphql.Fields{
			"id" : &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
				Description:"Unique identity of the token.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if token, ok := p.Source.(models.Token); ok {
						return token.Id
					}

					return nil, nil
				},
			},
			"publicKey": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description: "Public key for encrypting an application",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if token, ok := p.Source.(models.Token); ok {
						return token.PublicKey
					}
					return nil, nil
				},
			},
			"privateKey": &graphql.Field{
				Type: graphql.String,
				Description:"Private key for decrypting an application.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if app, ok := p.Source.(models.Token); ok {
						return app.PrivateKey
					}
					return nil, nil
				},
			},
			"hashType": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description: "The hash type used for encoding/decoding of jwt token",
				Resolve:func(p graphql.ResolveParams) (interface{}, error){
					if token, ok := p.Source.(models.Token); ok {
						return token.HashType
					}
					return nil, nil
				},

			},
			"added": &graphql.Field{
				Type: graphql.String,
				Description:"DateTime when the user is added",
				Resolve:func(p graphql.ResolveParams) (interface{}, error){
					if token, ok := p.Source.(models.Token); ok {
						return token.HashType
					}
					return nil, nil
				},
			},
			"lastUpdated": &graphql.Field{
				Type: graphql.String,
				Description:"DateTime when the user last update their data",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if token, ok := p.Source.(models.Token); ok {
						return  token.LastUpdated
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