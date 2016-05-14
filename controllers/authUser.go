package controllers

import "github.com/astaxie/beego"



// Operations about object
type AuthUserController struct {
	beego.Controller
}



// @router / [get]
func (authCtrl *AuthUserController) GetAll(){
	authCtrl.graphQLRequest()
}


// @router / [post]
func (authCtrl *AuthUserController) Post(){
	authCtrl.graphQLRequest()
}



func (this *AuthUserController) graphQLRequest(){

}