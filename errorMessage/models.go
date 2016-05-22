package errorMessage

import "errors"

//APPLICATION AREA
var (
	GRAPH_APP_NAME_NOT_FOUND  = errors.New("Application name or id not found")
	GRAPH_APP_ALREADY_PRESENT = errors.New("Application already present")
	ErrorAlreadyPresent  	  = errors.New("Error Node already present")
	ErrorIdNotPresent         = errors.New("Error Id property is not present in the node")

)
