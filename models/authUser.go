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
	Status string `orm:"default(active)"`
	Added time.Time `orm:"auto_now_add;type(datetime)"`
	LastUpdated time.Time `orm:"auto_now;type(datetime)"`
	Apps []*Application `orm:"null;reverse(many)"`
}





//Get user..
func (user *AuthUser)getUser() (err error){
	o := orm.NewOrm()
	o.Using("default")
	err = o.Read(&user)
	return
}


func (user *AuthUser) fetchApps() (num int, err error) {
	o := orm.NewOrm()
	o.Using("default")
	num, err = o.LoadRelated(&user, "Apps")
	return
}




//Used for registering a user....
func (user *AuthUser) create() (id int, err error) {
	// insert
	o := orm.NewOrm()
	o.Using("default")
	id, err = o.Insert(&user)
	return err
}


//Deactivate a user account...
func (user *AuthUser) deactivate() (int, error){
	o := orm.NewOrm()
	o.Using("default")
	user.Status = StatusMap["DEACTIVATED"]
	num, err := o.Update(&user)
	if err != nil{
		return nil, err
	}
	//Now also deactivate all its application
	_, err = o.QueryTable(new(Application)).Filter("owner_id", user.Id).Update(orm.Params{
		"status": StatusMap["DEACTIVATED"],
	})

	return num, err
}



