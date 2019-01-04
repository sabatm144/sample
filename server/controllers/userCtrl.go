package controller

import (
	"log"
	"net/http"
	"test_2/server/authentication"
	"test_2/server/dbCon"
	"test_2/server/models"
	"test_2/server/utils"

	"gopkg.in/mgo.v2/bson"
)

// RegisterUser
func RegisterUser(w http.ResponseWriter, r *http.Request) {

	login := models.Login{}
	if !parseJSON(w, r.Body, &login) {
		return
	}

	log.Println(login)
	db := dbCon.CopyMongoDB()
	defer db.Session.Close()

	//Validate user is already present
	query := bson.M{
		"emailID": bson.RegEx{Pattern: "^" + login.EmailID + "$", Options: "i"},
	}

	if c, _ := db.C("users").Find(query).Count(); c > 0 {
		renderERROR(w, http.StatusBadRequest, "Email already present, Login to continue")
		return
	}

	userIns := models.User{
		ID:       bson.NewObjectId(),
		EmailID:  login.EmailID,
		Password: utils.SHAEncoding(login.Password),
	}

	if err := db.C("users").Insert(&userIns); err != nil {
		renderERROR(w, http.StatusBadRequest, "Error in creating account, please try again")
		return
	}

	authBackend := authentication.InitJWTAuthenticationBackend()
	auth := authentication.AuthStruct{}
	auth.ID = userIns.ID.Hex()
	token, _ := authBackend.GenerateToken(auth)

	resp := struct {
		Message string              `json:"message"`
		Token   string              `json:"token"`
		User    models.LoggedInUser `json:"user"`
	}{
		"Added successfully, Complete your profile",
		token,
		models.LoggedInUser{
			userIns.ID, userIns.Name, userIns.EmailID,
		},
	}

	renderJSON(w, http.StatusOK, resp)
}

// AuthenticateUser
func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	login := models.Login{}
	if !parseJSON(w, r.Body, &login) {
		return
	}

	log.Println(login)
	db := dbCon.CopyMongoDB()
	defer db.Session.Close()

	user := models.LoggedInUser{}
	query := bson.M{
		"emailID":      bson.RegEx{Pattern: "^" + login.EmailID + "$", Options: "i"},
		"password":     utils.SHAEncoding(login.Password),
	}

	if err := db.C("users").Find(query).One(&user); err != nil {
		renderERROR(w, http.StatusBadRequest, "Error in authentication, please try again")
		return
	}

	authBackend := authentication.InitJWTAuthenticationBackend()
	auth := authentication.AuthStruct{}
	auth.ID = user.ID.Hex()
	token, _ := authBackend.GenerateToken(auth)

	resp := struct {
		Message  string              `json:"message"`
		Token    string              `json:"token"`
		Customer models.LoggedInUser `json:"customer"`
	}{
		"LoggedIn successfully",
		token,
		user,
	}

	renderJSON(w, http.StatusOK, resp)
}
