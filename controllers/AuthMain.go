package controllers

import "github.com/astaxie/beego"



// Operations about object
type AuthMainController struct {
	beego.Controller
}



// @router / [get]
func (authCtrl *AuthMainController) GetAll(){
	authCtrl.graphQLRequest()
}


// @router / [post]
func (authCtrl *AuthMainController) Post(){
	authCtrl.graphQLRequest()
}



func (this *AuthMainController) graphQLRequest(){

}