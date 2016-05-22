package Interfaces

//Interfaces define method which must be implememted by all graph databases..
type Graph interface {
	//Will merge the database if present already..
	Create() (err error)
	//Will return an error if the application already exists..
	CreateIfNotExist() (err error)
	Exist() (exist bool, err error)
	Delete() (err error)
	Update() (err error)
	Read() (err error)
	ReadAll([] *interface{}) (err error)
}