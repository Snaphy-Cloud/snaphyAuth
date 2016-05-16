package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

type Application struct {
	Id int
	Name string  `orm:"unique"`
	Status string `orm:"default(active)"`
	Added time.Time `orm:"auto_now_add;type(datetime)"`
	LastUpdated time.Time `orm:"auto_now;type(datetime)"`
	Owner *AuthUser `orm:"null;rel(fk)"`
	TokenInfo []* Token `orm:"null;reverse(many)"`
}


//Get the  application listed whose application id is given..
func (app *Application) GetApp()(err error){
	o := orm.NewOrm()
	o.Using("default")
	if app.Id != 0{
		err = o.Read(app)
	}else{
		err = o.Read(app, "Name")
	}
	return
}



//Used for creating an application....
func (app *Application) Create() (error) {
	// insert
	o := orm.NewOrm()
	o.Using("default")
	_, err := o.Insert(app)
	return err
}



//Deactivate a user account...
func (app *Application) Deactivate() (num int64, err error){
	o := orm.NewOrm()
	o.Using("default")
	app.Status = StatusMap["DEACTIVATED"]
	num, err = o.Update(app)
	//Now also deactivate all token whose...status is active..
	if err != nil{
		return 0, err
	}
	//Now also deactivate all its application
	//Only change status of those token whose status is Active
	_, err = o.QueryTable(new(Token)).Filter("application_id", app.Id).Filter("status", StatusMap["ACTIVE"]).Update(orm.Params{
		"status": StatusMap["DEACTIVATED"],
	})
	return
}

//Deactivate a user account...
func (app *Application) Activate() (num int64, err error){
	o := orm.NewOrm()
	o.Using("default")
	app.Status = StatusMap["ACTIVATE"]
	num, err = o.Update(app)
	//Now also deactivate all token whose...status is active..
	if err != nil{
		return 0, err
	}

	return
}


//Only delete a user by ID
func (app *Application) Delete() (num int64, err error){
	o := orm.NewOrm()
	o.Using("default")
	num, err = o.Delete(app)
	return
}


