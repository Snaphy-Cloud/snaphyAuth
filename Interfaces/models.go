/*Written by Robins*/
package Interfaces

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/ast"
	valid "github.com/asaskevich/govalidator"
)

var (
	UserInterface *graphql.Interface
	Application   *graphql.Interface
	Info          *graphql.Interface
	CreatedOn     *graphql.Interface

	EmailType     *graphql.Scalar

	StatusEnum    *graphql.Enum
)



func init() {
	//Defigning a custom variable email type.
	EmailType = graphql.NewScalar(graphql.ScalarConfig{
		Name: "Email",
		Serialize:func(value interface{}) (interface{}) {
			return value
		},
		ParseValue:func(value interface{}) interface{}{
			return value
		},
		ParseLiteral:func(valueAST ast.Value) interface{}{
			if(valueAST.GetKind() != kinds.StringValue){
				if ok := valid.IsEmail(valueAST.GetValue()); ok{
					return valueAST.GetValue()
				}
			}
			return nil
		},
	})


	StatusEnum = graphql.NewEnum(graphql.EnumConfig{
		Name: "Status",
		Description:"Status enum showing active/inactive status",
		Values:graphql.EnumValueConfigMap{
			"ACTIVE": graphql.EnumValueConfig{
				Value: "active",
				Description: "Shows active status",
			},

			"INACTIVE": graphql.EnumValueConfig{
				Value: "inactive",
				Description: "Shows inactive status",
			},

			"DISABLED": graphql.EnumValueConfig{
				Value: "disabled",
				Description: "Shows disabled status",
			},

			"DEACTIVATED": graphql.EnumValueConfig{
				Value: "deactivated",
				Description: "Shows deactivated status",
			},

		},
	})


	UserInterface = graphql.NewInterface(graphql.InterfaceConfig{
		Name: "User",
		Description: "A user defined in snaphy auth cloud",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The id of the user.",
			},
			"firstName": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description:"First Name of User. Required field",
			},
			"lastname":&graphql.Field{
				Type: graphql.String,
				Description: "Last Name of User",
			},
			"email": &graphql.Field{
				Type: graphql.NewNonNull(EmailType),
				Description: "Email of the User",
			},
		},
	})


	Application = graphql.NewInterface(graphql.InterfaceConfig{
		Name: "Application",
		Description:"Application model interface snaphy auth cloud",
		Fields:graphql.Fields{
			"id" : &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description:"Unique identity of the application.",
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description: "Application name for graphql",
			},
		},
	})


	Info = graphql.NewInterface(graphql.InterfaceConfig{
		Name: "Info",
		Description: "Interface for displaying information.",
		Fields: graphql.Fields{
			"status": &graphql.Field{
				Type: StatusEnum,
				Description: "Show status of a data",
			},
		},
	})



	CreatedOn = graphql.NewInterface(graphql.InterfaceConfig{
		Name: "CreatedOn",
		Description: "Interface for displaying date for vaious type of events like added/ last updated.",
		Fields: graphql.Fields{
			"added": &graphql.Field{
				Type: graphql.String,
				Description: "Show  date when the model is created",
			},

			"lastUpdated": &graphql.Field{
				Type: graphql.String,
				Description: "Show  date when the model is last updated",
			},
		},
	})


}