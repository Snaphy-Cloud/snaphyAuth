package schemas

import (
	"github.com/graphql-go/graphql"
	snaphyInterface "snaphyAuth/Interfaces"
	"snaphyAuth/models"
)



//Schemas used in application..
var (
	AuthUserType         *graphql.Object
	AuthUserQueryType    *graphql.Object
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


	AuthUserQueryType = graphql.NewObject(graphql.ObjectConfig {
		Name: "AuthUserQuery",
		Description: "Snaphy cloud main Auth type for storing all application register info.",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: AuthUserType,
				Description: "Return the AuthUser model",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Description: "find AuthUser by ID",
						Type: graphql.ID,
					},
					"email": &graphql.ArgumentConfig{
						Description: "find AuthUser by Email",
						Type: snaphyInterface.EmailType,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if user, ok := p.Source.(models.AuthUser); ok {
						return user.Id, nil
					}

					id  := p.Args["id"].(int)
					email  := p.Args["email"].(string)

					if(id != nil){
						//Now find auth user by ID
						authUser := models.AuthUser{Id: id}
						if err := authUser.GetUser(); err == nil{
							return authUser
						}else{
							nil, err
						}
					}else if(email != nil){
						//Now find auth user by ID
						authUser := models.AuthUser{Email: email}
						if err := authUser.GetUser(); err == nil{
							return authUser
						}else{
							nil, err
						}

					}

					return nil, nil
				},
			},

		},
	})




}


