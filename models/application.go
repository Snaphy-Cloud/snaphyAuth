package models

import "time"

type Application struct {
	Id int
	Name string
	Status string `orm:"default(active)"`
	Added time.Time `orm:"auto_now_add;type(datetime)"`
	LastUpdated time.Time `orm:"auto_now;type(datetime)"`
	Owner *AuthUser `orm:"null;rel(fk)"`
	TokenInfo []* Token `orm:"null;reverse(many)"`
	Database *DbIndex `orm:"null;reverse(one)"`
}


