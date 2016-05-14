package models

import "time"

type DbIndex struct  {
	Id int
	Name string
	DbUser string
	DbPass string
	Application *Application `orm:"rel(fk)"`
	Status string `orm:"default('active')"`
	Added time.Time `orm:"auto_now_add;type(datetime)"`
	LastUpdated time.Time `orm:"auto_now;type(datetime)"`
}

