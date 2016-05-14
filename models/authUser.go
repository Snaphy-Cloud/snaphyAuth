package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)



type AuthUser struct {
	Id int
	FirstName string
	LastName string
	Email string
	Status string `orm:"default('active')"`
	Added time.Time `orm:"auto_now_add;type(datetime)"`
	LastUpdated time.Time `orm:"auto_now;type(datetime)"`
	Application []*Application `orm:"null;reverse(many)"`
}





//Get user..
func (user *AuthUser)getUser() (err error){
	o := orm.NewOrm()
	o.Using("default")
	err = o.Read(&user)
	return err
}



func (user *AuthUser) Save() error {
	// insert
	o := orm.NewOrm()
	o.Using("default")
	_, err := o.Insert(&user)
	return err
}



