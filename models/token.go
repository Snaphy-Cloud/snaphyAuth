package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"crypto/rand"
	"crypto/rsa"
	"github.com/satori/go.uuid"
	"github.com/astaxie/beego"
)

type Token struct  {
	Id int
	PublicKey string `orm:"unique"`
	PrivateKey string `orm:"unique"`
	HashType string
	AppSecret string `orm:"unique"`
	AppId string `orm:"unique"`
	Application *Application `orm:"rel(fk)"`
	Status string `orm:"default(active)"`
	Added time.Time `orm:"auto_now_add;type(datetime)"`
	LastUpdated time.Time `orm:"auto_now;type(datetime)"`
}


//Method for generating token..

//Get token
func (token *Token) getToken()(err error){
	o := orm.NewOrm()
	o.Using("default")
	err = o.Read(&token)
	return
}



//Used for creating a token..
//Only Application
func (token *Token) create() (id int, err error) {
	// insert
	o := orm.NewOrm()
	o.Using("default")
	var privateKey *rsa.PrivateKey
	//generate private keys.
	privateKey, err = generateKeys()
	if err != nil{
		return
	}
	//Now store private key
	token.PrivateKey = ""+privateKey
	token.PublicKey = ""+ &privateKey.PublicKey
	token.HashType = beego.AppConfig.String("jwt::algorithm")
	token.AppId = uuid.NewV4()
	token.AppSecret = uuid.NewV4()

	//Get the appId.
	id, err = o.Insert(&token)
	return
}



//Only delete a token by ID
func (token *Token) delete() (num int64, err error){
	o := orm.NewOrm()
	o.Using("default")
	num, err = o.Delete(token)
	return
}









//Method to generate public and private keys..
func generateKeys() (privateKey *rsa.PrivateKey, err error){
	// Generate RSA Keys
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	return
}


