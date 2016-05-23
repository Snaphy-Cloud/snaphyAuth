package models

import (
	"github.com/jmcvetta/neoism"
	"strings"
	"fmt"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
	"snaphyAuth/Interfaces"
	"github.com/astaxie/beego"
)



type Token struct{
	IAT int64 //Issued at
	ISS int64 //User Identity
	EXP int64 //Expiry Time
	JTI string //Unique string identifies a token
	GRP string //Group
	KID string //AppId its not applicationID but AppId in TokenHelper file to track application.....
}


//Check interface implementation..
//Will throw error if the struct doesn't implements Graph Interface..
var _ Interfaces.Graph = (*Token)(nil)


func init(){
	token := new(Token)
	err := token.AddUniqueConstraint()
	if err != nil{
		beego.Trace("Error creating UNIQUE constraint on Token database")
		beego.Trace(err)
	}
}




func (token *Token)AddUniqueConstraint() (err error){
	stmt := "CREATE CONSTRAINT ON (token:Token) ASSERT token.JTI IS UNIQUE"
	cq := neoism.CypherQuery{
		Statement: stmt,
	}
	// Issue the query.
	err = db.Cypher(&cq)
	return
}



//Add a tag with relationship..
func (token *Token) AddTag(tag *TokenTag) (err error){
	stmt := `MATCH (tag: Label{ name:{labelName}, appId: {appId}, realmName: {realm} })
	         MATCH (token:Token) WHERE token.name = {tokenString} AND token.appId = {appId} AND token.realmName = {realm}
	         MERGE (tag) - [role:Role] -> (token)`

	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"labelName": tag.Name, "appId": token.AppId, "realm": token.RealmName, "tokenString": token.TokenString},
	}

	// Issue the query.
	err = db.Cypher(&cq)
	return
}



//Verify the token before parsing .. the token just checks if the tokens ia a valid one.
func (nodeToken *Token) VerifyHash(tokenHelper *TokenHelper) (valid bool, err error){

	parts := strings.Split(nodeToken.TokenString, ".")

	method := jwt.GetSigningMethod(tokenHelper.HashType)
	err = method.Verify(strings.Join(parts[0:2], "."), parts[2], []byte(tokenHelper.PublicKey))

	if err != nil{

		return false, err
	}else{
		return true, err
	}

}




//Parses the token value..And also validates the algorithm..
func (nodeToken *Token) VerifyAndParse() (valid bool, err error){
	var token *jwt.Token
	token, err = jwt.Parse(nodeToken.TokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		//Also check the token expiry date
		//Also check if the token in present in the db.

		return LookUpKey(token.Claims["kid"].(string))
	})

	nodeToken.Added = int64(token.Claims["iat"].(float64))
	nodeToken.Expiry = int64(token.Claims["exp"].(float64))
	nodeToken.JTI = token.Claims["jti"].(string)
	nodeToken.GroupName = token.Claims["grp"].(string)
	nodeToken.Roles = token.Claims["roles"].(string)
	nodeToken.RealmName = token.Claims["realm"].(string)
	return token.Valid, err
}


//Find the public key with the provided data..
func LookUpKey(appId string) (publicKey []byte, err error){
	if appId == ""{
		return nil, errors.New("Error: appId value is null")
	}
	tokenHelper := new(TokenHelper)
	tokenHelper.AppId = appId
	err = tokenHelper.GetToken()

	if err != nil {
		return nil, err
	}else{
		return []byte(tokenHelper.PublicKey), nil
	}
}



//Check if token is expired or not valid or valid...
//TRUE IF EXPIRED AND FALSE IF NOT
func (nodeToken *Token) CheckExpired() (expired bool, err error){
	if nodeToken.Expiry == 0 {
		return true, errors.New("Expiry claim not present in Token model.")
	}else{
		now := time.Now().Unix()
		//TODO ALSO CHECK THE STATUS IN GRAPHDATABASE FIRST
		//TODO CREATE METHOD GETSTATUS() in NODE TOKEN
		if now >= nodeToken.Expiry{
			return true, nil
		}else{
			return false, nil
		}
	}

}


//Find the token from the database and populate the data..
//Provide the nodeToken with jwt field.
func (nodeToken *Token) GetToken() (err error){
	/*
	TokenString string //JWT TOKEN INFO
	AppId int
	RealmName string
	GroupName string
	Status string
	UserId int
	Added int64
	Expiry int64
	LastUpdated int64
	JTI string //Unique token provider.
	Roles string*/
	//First try to match the token..
	/*stmt := `MATCH (app: Application{name: {appName}, id: {appId} })
	         MATCH (app)-[org:Organization]->()-[type:Type]->()-[identity:Identity{userId: {userId} }]->(token:Token{jti: {jti}})`*/
	stmt := `MATCH (token:Token{jti: {jti} }) RETURN token.status AS status`
}


//TODO WRITE A METHOD TO UPDATE STATUS IF FOUND INVALID..TOKEN
//TODO WRITE A METHOD TO FIND IF THE TOKEN ALREADY PRESENT IN THE DATABASE..


