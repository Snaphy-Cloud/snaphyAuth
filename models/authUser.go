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
func (user *AuthUser)GetUser() (err error){
	if user.Id != 0{
		o := orm.NewOrm()
		o.Using("default")
		err = o.Read(user)
	}else if(user.Email != ""){
		o := orm.NewOrm()
		o.Using("default")
		err = o.Read(user, "Email")
	}

	return
}

//Get user..
func (user *AuthUser)GetCustomUser(key string) (err error){
	o := orm.NewOrm()
	o.Using("default")
	err = o.Read(user, key)
	return
}


func (user *AuthUser) FetchApps() (num int64, err error) {
	o := orm.NewOrm()
	o.Using("default")
	num, err = o.LoadRelated(user, "Apps")
	return
}




//Used for registering a user....
func (user *AuthUser) Create() (id int64, err error) {
	// insert
	o := orm.NewOrm()
	o.Using("default")
	id, err = o.Insert(user)
	return id, err
}


//Deactivate a user account...
func (user *AuthUser) Deactivate() (int64, error){
	o := orm.NewOrm()
	o.Using("default")
	user.Status = StatusMap["DEACTIVATED"]
	num, err := o.Update(user)
	if err != nil{
		return 0, err
	}
	//Now also deactivate all its application
	_, err = o.QueryTable(new(Application)).Filter("owner_id", user.Id).Update(orm.Params{
		"status": StatusMap["DEACTIVATED"],
	})

	return num, err
}

//Deactivate a user account...
func (user *AuthUser) Activate() (int64, error){
	o := orm.NewOrm()
	o.Using("default")
	user.Status = StatusMap["ACTIVE"]
	num, err := o.Update(user)
	if err != nil{
		return 0, err
	}
	return num, err
}


//Only delete a user by ID
func (user *AuthUser) Delete() (num int64, err error){
	o := orm.NewOrm()
	o.Using("default")
	num, err = o.Delete(user)
	return
}



