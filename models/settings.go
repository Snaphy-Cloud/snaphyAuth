package models

import "time"

type Settings struct {
	Id int
	Added time.Time `orm:"auto_now_add;type(datetime)"`
	LastUpdated time.Time `orm:"auto_now;type(datetime)"`
	Application *Application `orm:"reverse(one)"` // Reverse relationship (optional)
	TokenExpired int //Time duration in minutes after which token will expired.
}
