package controller

import (
	"log"
	"net/http"
	"sample/server/dbCon"
	"sample/server/models"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

func Vote(w http.ResponseWriter, r *http.Request) {

	ID := r.Context().Value("loggedInUserId").(string)
	if !bson.IsObjectIdHex(ID) {
		renderJSON(w, http.StatusBadRequest, "Not a valid user ID")
		return
	}
	userID := bson.ObjectIdHex(ID)
	log.Printf("LoggedIn UserID %s", userID)

	db := dbCon.CopyMongoDB()
	defer db.Session.Close()

	e := appErr{}
	//User is present or n't
	userIns := models.User{}
	err := db.C("users").FindId(userID).One(&userIns)
	if err != nil {
		e = appErr{Message: "User not found", Error: err.Error()}
		renderJSON(w, http.StatusNotFound, e)
		return
	}

	params := r.Context().Value("params").(httprouter.Params)
	if !bson.IsObjectIdHex(params.ByName("id")) {
		renderJSON(w, http.StatusBadRequest, "Not a valid content ID")
		return
	}
	contentID := bson.ObjectIdHex(params.ByName("id"))
	log.Printf("Content ID %s", contentID)

	contentIns := models.Content{}
	err = db.C("contents").FindId(contentID).One(&contentIns)
	if err != nil {
		e = appErr{Message: "User not found", Error: err.Error()}
		renderJSON(w, http.StatusNotFound, e)
		return
	}

	type Vote struct {
		Status bool `json:status`
	}

	voteIns := Vote{}
	if !parseJSON(w, r.Body, &voteIns) {
		return
	}

	vIns := models.Voter{}
	db.C("votes").Find(bson.M{"contentID": contentIns.ID, "userID": userID}).One(&vIns)
	log.Println(vIns, voteIns.Status)

	vIns.Status = voteIns.Status
	if vIns.ID.Hex() == "" {
		vIns.ID = bson.NewObjectId()
		vIns.ContentID = contentIns.ID
		vIns.UserID = userID
	}

	log.Println(vIns)
	_, err = db.C("votes").UpsertId(vIns.ID, bson.M{"contentID": vIns.ContentID, "userID": vIns.UserID, "status": vIns.Status})
	if err != nil {
		e = appErr{Message: "Could n't update/insert vote!", Error: err.Error()}
		renderJSON(w, http.StatusBadRequest, e)
		return
	}

	e = appErr{Message: "Status Updated!"}
	renderJSON(w, http.StatusOK, e)
}
