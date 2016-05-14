package schemas

import (
	"github.com/graphql-go/graphql"
	snaphyInterface "snaphyAuth/Interfaces"
	"snaphyAuth/models"
)

var (
	TokenType      *graphql.Object
)


func init(){
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
			"status": &graphql.Field{
				Type: snaphyInterface.StatusEnum,
				Description:"status of Application",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if token, ok := p.Source.(models.Token); ok {
						return token.Status
					}

					return nil, nil
				},
			},
			"added": &graphql.Field{
				Type: graphql.String,
				Description:"DateTime when the user is added",
				Resolve:func(p graphql.ResolveParams) (interface{}, error){
					if token, ok := p.Source.(models.Token); ok {
						return token.Added
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

			//TODO ADD RELAY CONNECTION FOR ADDING APPLICATION RELATIONS LATER.

		},
		Interfaces: [] *graphql.Interface{
			snaphyInterface.TokenInterface,
			snaphyInterface.InfoInterface,
			snaphyInterface.CreatedOnInterface,
		},
	})
}
