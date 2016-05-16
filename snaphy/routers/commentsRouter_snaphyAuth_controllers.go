package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["snaphyAuth/controllers:AuthUserController"] = append(beego.GlobalControllerRouter["snaphyAuth/controllers:AuthUserController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["snaphyAuth/controllers:AuthUserController"] = append(beego.GlobalControllerRouter["snaphyAuth/controllers:AuthUserController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

}
