package models

import "time"

type Token struct  {
	Id int
	PublicKey string
	PrivateKey string
	HashType string
	AppSecret string
	AppId string
	Application *Application `orm:"rel(fk)"`
	Status string `orm:"default(active)"`
	Added time.Time `orm:"auto_now_add;type(datetime)"`
	LastUpdated time.Time `orm:"auto_now;type(datetime)"`
}

