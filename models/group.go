package models

import (
	"github.com/jmcvetta/neoism"
	"snaphyAuth/Interfaces"
	"strconv"
	"github.com/dgrijalva/jwt-go"
	"errors"
	"time"
	"snaphyAuth/helper"
)



type Group struct {
	Name string
	AppId int
	RealmName string
}



//Check interface implementation..
//Will throw error if the struct doesn't implements Graph Interface..
var _ Interfaces.Graph = (*Group)(nil)



func (group *Group) Delete() (err error){
	stmt := `MATCH p =(begin:Group{name: {groupName}, appId: {appId}, realmName: {realmName} })-[r*]->(END:Token)  DETACH DELETE p`
	cq := neoism.CypherQuery{
		Statement: stmt,
		Parameters: neoism.Props{"groupName": group.Name, "appId": group.AppId, "realmName": group.RealmName},
	}

	// Issue the query.
	err = db.Cypher(&cq)
	return
}








//Create and sign the token..
func (group *Group) CreateToken(token *NodeToken, app *Application, settings *ApplicationSettings, tokenHelper *TokenHelper, realm *Realm, tag *NodeTag, userIdentity string) (err error){
	// Create the token
	var signToken *jwt.Token

	if tokenHelper.HashType == "" {
		return errors.New("No Signing method present in the TokenHelper Database")
	}else if(tokenHelper.HashType == "RS256"){
		signToken = jwt.New(jwt.SigningMethodRS256)
	}else{
		//Right now default is RS256
		signToken = jwt.New(jwt.SigningMethodRS256)
	}

	//signToken := jwt.GetSigningMethod(tokenHelper.HashType)

	duration := settings.ExpiryDuration
	if duration == 0 {
		duration = 3600
	}

	expiryTime := time.Now().Add(time.Second * time.Duration(duration)).Unix() //time after it will be invalid..
	issuedAt := time.Now().Unix() //Issuer at time..
	jti := helper.CreateUUID()
	signToken.Claims["exp"] = expiryTime
	signToken.Claims["iat"] = issuedAt
	signToken.Claims["iss"] = userIdentity //user identity..
	signToken.Claims["grp"] = group.Name //Group with which user belong to.
	signToken.Claims["realm"] = realm.Name
	//TODO LATER ADD MORE ROLES BY FETCHING RELATION FROM GRAPH DATABASE..
	signToken.Claims["roles"] = tag.Name
	//JUST STORE THIS INSTEAD OF STORING TOKENS..
	signToken.Claims["jti"] = jti
	//Add kid to track token to public and private keys..
	signToken.Claims["kid"] = tokenHelper.AppId

	if tokenHelper.PrivateKey != ""{

		// Sign and get the complete encoded token as a string
		tokenString, err := signToken.SignedString([]byte(tokenHelper.PrivateKey))
		if err != nil{
			return err
		}else{
			userId, _ := strconv.Atoi(userIdentity)
			token.Added = issuedAt
			token.LastUpdated = issuedAt
			token.AppId = app.Id
			token.RealmName = realm.Name
			token.Status = StatusMap["ACTIVE"]
			token.TokenString = tokenString
			token.UserId = userId
			token.JTI = jti
		}

		return err
	}else{
		return errors.New("Private key not present in tokenHelper")
	}
}

