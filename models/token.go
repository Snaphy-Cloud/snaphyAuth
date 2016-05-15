package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"crypto/rand"
	"crypto/rsa"
	"github.com/satori/go.uuid"
	"github.com/astaxie/beego"
	"encoding/pem"
	"crypto/x509"
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


func init(){
	//TODO Test performance benchmark for these key generation
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
func (token *Token) create() (id int64, err error) {
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

	token.PrivateKey, err = GeneratePem(privateKey)
	token.PublicKey, err = GeneratePub(privateKey)
	token.HashType = beego.AppConfig.String("jwt::algorithm")
	token.AppId = uuid.NewV4().String()
	token.AppSecret = uuid.NewV4().String()

	if err != nil{
		return  0, err
	}
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





//Generate private key in pem format..
func GeneratePem(privateKey *rsa.PrivateKey)(string, error){
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type: "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	return string(pemdata), nil
}

//http://stackoverflow.com/questions/13555085/save-and-load-crypto-rsa-privatekey-to-and-from-the-disk
//Generate public  key file pub format..
func GeneratePub(privateKey *rsa.PrivateKey)(string, error){
	PubASN1, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		// do something about it
		return "", err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: PubASN1,
	})

	return string(pubBytes), err
}






//Method to generate public and private keys..
func generateKeys() (privateKey *rsa.PrivateKey, err error){
	// Generate RSA Keys
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	return
}


