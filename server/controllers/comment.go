package controller

import (
	"log"
	"net/http"
	"sample/server/dbCon"
	"sample/server/models"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

//Comments returns the comments for a content
func Comments(w http.ResponseWriter, r *http.Request) {

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

	comments := []models.Comment{}
	err = db.C("comments").Find(bson.M{"contentID": contentID}).All(&comments)
	if err != nil {
		e = appErr{Message: "Couldn't find comments!", Error: err.Error()}
		renderJSON(w, http.StatusBadRequest, e)
		return
	}

	res := make(map[string]interface{})
	res["comments"] = comments
	renderJSON(w, http.StatusOK, res)
	return
}

//Comment is to add comment for a content
func Comment(w http.ResponseWriter, r *http.Request) {

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

	type commentDesc struct {
		Text string `json:"text" bson:"text"`
	}

	commentIns := commentDesc{}
	if !parseJSON(w, r.Body, &commentIns) {
		return
	}

	comment := &models.Comment{}
	comment.ID = bson.NewObjectId()
	comment.ContentID = contentID
	comment.UserID = userID
	comment.Text = commentIns.Text
	err = db.C("comments").Insert(comment)
	if err != nil {
		e = appErr{Message: "Could n't insert user comment!", Error: err.Error()}
		renderJSON(w, http.StatusBadRequest, e)
		return
	}
	renderJSON(w, http.StatusOK, "Comment posted successfully")
	return
}

//Reply is to add reply for a comment
func Reply(w http.ResponseWriter, r *http.Request) {

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
		renderJSON(w, http.StatusBadRequest, "Not a valid commentID ID")
		return
	}
	commentID := bson.ObjectIdHex(params.ByName("id"))
	log.Printf("Content ID %s", commentID)

	type commentDesc struct {
		Text string `json:"text" bson:"text"`
	}
	commentIns := commentDesc{}
	if !parseJSON(w, r.Body, &commentIns) {
		return
	}

	comment := &models.Comment{}
	err = db.C("comments").FindId(commentID).One(comment)
	if err != nil {
		e = appErr{Message: "Could n't find comment!", Error: err.Error()}
		renderJSON(w, http.StatusBadRequest, e)
		return
	}

	reply := models.Comment{}
	reply.ID = bson.NewObjectId()
	reply.ContentID = comment.ContentID
	reply.UserID = comment.UserID
	reply.Text = commentIns.Text
	comment.Replies = append(comment.Replies, reply)
	_, err = db.C("comments").UpsertId(commentID, &comment)
	if err != nil {
		e = appErr{Message: "Could n't update/insert user reply!", Error: err.Error()}
		renderJSON(w, http.StatusBadRequest, e)
		return
	}
	renderJSON(w, http.StatusOK, "Reply posted successfully")
	return
}
