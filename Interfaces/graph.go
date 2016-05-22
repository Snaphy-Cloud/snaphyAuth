package Interfaces

//Interfaces define method which must be implememted by all graph databases..
type Graph interface {
	Create() (err error)
	Exist() (exist bool, err error)
	Delete() (err error)
	CreateIfNotExist() (err error)
	Update() (err error)
	Read() (err error)
	ReadAll([] *interface{}) (err error)
}