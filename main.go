package main

import (
	_ "snaphyAuth/docs"
	_ "snaphyAuth/routers"
	_ "snaphyAuth/models"
	"github.com/astaxie/beego"

	"fmt"
)


func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	configLog()
	beego.Run()
}



func configLog(){
	fileName := beego.AppConfig.String("logs:fileName")
	if  fileName != "" {
		beego.SetLogger("file", fmt.Sprintf(`{"filename": %s}`, fileName))
	}
	if beego.BConfig.RunMode != "dev"{
		//Remove data from console..
		beego.BeeLogger.DelLogger("console")
	}
	beego.BeeLogger.Async()
	//beego.SetLogFuncCall(true)
}
