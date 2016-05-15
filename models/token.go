package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

type Token struct  {
	Id int
	PublicKey string `orm:"unique"`
	PrivateKey string `orm:"unique"`
	HashType string
	AppSecret string `orm:"unique"`
	AppId string `orm:"unique"`
	Application *Application `orm:"rel(fk)"`
	Status string `orm:"default(active)"`
	Added time.Time `orm:"auto_now_add;type(datetime)"`
	LastUpdated time.Time `orm:"auto_now;type(datetime)"`
}


//Method for generating token..

//Get token
func (token *Token) getToken()(err error){
	o := orm.NewOrm()
	o.Using("default")
	err = o.Read(&token)
	return
}



//Used for registering a user....
//Only Application
func (token *Token) save() (error) {
	// insert
	o := orm.NewOrm()
	o.Using("default")
	_, err := o.Insert(&app)
	return err
}


