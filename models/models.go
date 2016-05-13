package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"fmt"
)

func init(){
	//Register model
	registerModel()
	// orm.RegisterDataBase("default", "mysql", "root:root@/orm_test?charset=utf8")
	orm.RegisterDriver("postgres", orm.DRPostgres)
	database, user, password := getDatabaseCredentials()
	connString := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", user, password, database)
	orm.RegisterDataBase("default", "postgres", connString )
	name := "default"
	force := false
	verbose := true
	//Default value of debug is false
	debug := beego.AppConfig.DefaultBool("model:prefix", false)
	orm.Debug = debug
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		fmt.Println(err)
	}
}


func registerModel(){
	modelPrefix := beego.AppConfig.String("model::prefix")
	if modelPrefix != nil || modelPrefix != "" {
		return orm.RegisterModelWithPrefix(modelPrefix, new(AuthUser))
	}

	return orm.RegisterModel(new(AuthUser))
}



//Get the database name, username and password info for postgresql.
func getDatabaseCredentials() (string, string, string){
	database := beego.AppConfig.String("postgres::database")
	user := beego.AppConfig.String("postgres::user")
	password := beego.AppConfig.String("postgres::password")
	return  database, user, password
}
