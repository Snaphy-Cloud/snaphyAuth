package models

import (
	_ "github.com/lib/pq"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"fmt"
	_ "os"
	"os"
)


func init(){
	//First register the model..
	RegisterModel( new(AuthUser), new(Application), new(Token), new(DbIndex) )
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
	verbose := true
	//Default value of debug is false
	debug := beego.AppConfig.DefaultBool("model:prefix", false)
	orm.Debug = debug
	err = orm.RunSyncdb(name, force, verbose)
	return
}



func RegisterModel(snaphyModel ...interface{}){
	modelPrefix := beego.AppConfig.String("model::prefix")
	if modelPrefix != "" {
		orm.RegisterModelWithPrefix(modelPrefix, snaphyModel...)
	}else{
		orm.RegisterModel(snaphyModel...)
	}
}



//Get the database name, username and password info for postgresql.
func getDatabaseCredentials() (string, string, string){
	database := beego.AppConfig.String("postgres::database")
	user := beego.AppConfig.String("postgres::user")
	password := beego.AppConfig.String("postgres::password")
	return  database, user, password
}
