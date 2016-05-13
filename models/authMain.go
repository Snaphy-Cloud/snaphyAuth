package models

import (

)
import (
	"time"
	"github.com/astaxie/beego/orm"
	"fmt"
)

type AuthUser struct {
	Id int
	FirstName string
	LastName string
	Email string
	Added time.Time `orm:"auto_now_add;type(datetime)"`
}


func init(){
	//orm.RegisterModel(new(AuthUser))
	orm.RegisterModelWithPrefix("snaphy_auth_", new(AuthUser))
	// orm.RegisterDataBase("default", "mysql", "root:root@/orm_test?charset=utf8")
	orm.RegisterDriver("postgres", orm.DRPostgres)
	database, user, password := getDatabaseCredentials()
	connString := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", user, password, database)

	orm.RegisterDataBase("default", "postgres", connString )
	name := "default"
	force := false
	verbose := true
	orm.Debug = true
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		fmt.Println(err)
	}
}

