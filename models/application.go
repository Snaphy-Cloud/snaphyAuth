package models

import "time"

type Application struct {
	Id int
	UserId int
	DatabaseId int
	Name string
	Status string
	Added time.Time `orm:"auto_now_add;type(datetime)"`
	LastUpdated time.Time `orm:"auto_now_add;type(datetime)"`
	owner AuthUser `-`
}


func init(){
	RegisterModel(new(AuthUser))
}
