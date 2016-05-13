package models

import (

)
import (
	"time"
)

type AuthUser struct {
	Id int
	FirstName string
	LastName string
	Email string
	Added time.Time `orm:"auto_now_add;type(datetime)"`
}


func init(){

}

