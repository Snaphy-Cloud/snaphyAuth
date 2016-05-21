package models

import "time"

type ApplicationSettings struct {
	Id int
	ExpiryDuration int64 //Time in seconds after which it token will expired..
	Added time.Time `orm:"auto_now_add;type(datetime)"`
	LastUpdated time.Time `orm:"auto_now;type(datetime)"`
	Application *Application `orm:"reverse(one)"` // Reverse relationship (optional)
}