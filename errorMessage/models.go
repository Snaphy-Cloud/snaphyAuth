package errorMessage

import "errors"

//APPLICATION AREA
var (
	GRAPH_APP_NAME_NOT_FOUND  = errors.New("Application name or id not found")
	GRAPH_APP_ALREADY_PRESENT = errors.New("Application already present")
	ErrorAlreadyPresent  	  = errors.New("Error Node already present")
	ErrorIdNotPresent         = errors.New("Error Id property is not present in the node")
	TokenNotValid             = errors.New("Token string not valid")
	AppIdNull                 = errors.New("Token helper app id cannot be null")
	TokenJTINotPresent        = errors.New("Tokens JTI key not present in the token model")
	TokenStatusNotPresent     = errors.New("Tokens model doesn't have status present")
	TokenUserIdNotPresent     = errors.New("User Id  not present in the token model")
	TokenRealmNotPresent      = errors.New("Token model realm name cannot be empty")
	TokenGroupNotPresent      = errors.New("Token model group name cannot be empty")
)
