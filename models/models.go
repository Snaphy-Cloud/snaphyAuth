package models

import (
	_ "github.com/lib/pq"
	"github.com/astaxie/beego/orm"
	"fmt"
	"os"
	"github.com/astaxie/beego"
)


func init(){
	//First register the model..
	RegisterModel( new(AuthUser), new(Application), new(TokenHelper), new(DbIndex), new(ApplicationSettings) )
	//Now register the database..
	err := registerDb()
	if err != nil{
		beego.Trace(err)
		os.Exit(1)
	}
	//Initialize status..
	initStatus()
}




var (
	StatusMap map[string]string
)


func initStatus(){
	StatusMap = make(map[string]string)
	StatusMap["ACTIVE"] = "active"
	StatusMap["INACTIVE"] = "inactive"
	StatusMap["DISABLED"] = "disabled"
	StatusMap["DEACTIVATED"] = "DEACTIVATED"
}



func registerDb() (err error){
	// orm.RegisterDataBase("default", "mysql", "root:root@/orm_test?charset=utf8")
	orm.RegisterDriver("postgres", orm.DRPostgres)
	database, user, password := getDatabaseCredentials()
	connString := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", user, password, database)
	orm.RegisterDataBase("default", "postgres", connString )
	name := "default"
	force := false
	verbose := beego.AppConfig.DefaultBool("model::debug", false)
	//Default value of debug is false
	debug := beego.AppConfig.DefaultBool("model::debug", false)
	orm.Debug = debug
	err = orm.RunSyncdb(name, force, verbose)
	return
}



func RegisterModel(snaphyModel ...interface{}){
	modelPrefix := beego.AppConfig.DefaultString("model::prefix", "snaphy_auth_")
	if modelPrefix != "" {
		orm.RegisterModelWithPrefix(modelPrefix, snaphyModel...)
	}else{
		orm.RegisterModel(snaphyModel...)
	}
}



//Get the database name, username and password info for postgresql.
func getDatabaseCredentials() (string, string, string){
	database := beego.AppConfig.DefaultString("postgres::database", "robins")
	user := beego.AppConfig.DefaultString("postgres::user", "robins")
	password := beego.AppConfig.DefaultString("postgres::password", "12345")
	return  database, user, password
}
