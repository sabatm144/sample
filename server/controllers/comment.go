package controller

import (
	"log"
	"net/http"
	"test_2/server/dbCon"
	"test_2/server/models"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

func ListComments(w http.ResponseWriter, r *http.Request) {

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
	res["commentList"] = comments
	res["comments"] = len(comments)
	renderJSON(w, http.StatusOK, res)
	return
}

func NestedComments(w http.ResponseWriter, r *http.Request) {

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
		Text      string        `json:"text" bson:"text"`
		CommentID bson.ObjectId `json:"id,omitempty" bson:"id,omitempty"`
		// ChildID bson.ObjectId `json:"childID,omitempty" bson:"childID,omitempty"`
	}
	commentIns := commentDesc{}
	if !parseJSON(w, r.Body, &commentIns) {
		return
	}

	log.Println(commentIns.CommentID, contentID)
	nestedComments := &models.Comment{}
	if commentIns.CommentID != "" {
		nestedComments.IsAParent = true

		//parent level
		err := db.C("comments").Find(bson.M{"_id": commentIns.CommentID, "contentID": contentID}).One(nestedComments)
		if err != nil {
			e = appErr{Message: "Could n't update user post!", Error: err.Error()}
			renderJSON(w, http.StatusBadRequest, e)
			return
		}

		log.Println(nestedComments)
		if nestedComments.ID.Hex() != "" {
			nestedComments.IsAParent = true
			childComments := models.Comment{}
			childComments.ID = bson.NewObjectId()
			childComments.ContentID = contentID
			childComments.UserID = userID
			childComments.Text = commentIns.Text
			nestedComments.Child = append(nestedComments.Child, childComments)
		}
	}

	if commentIns.CommentID == "" {
		nestedComments.ID = bson.NewObjectId()
		nestedComments.ContentID = contentID
		nestedComments.UserID = userID
		nestedComments.Text = commentIns.Text
	}

	_, err = db.C("comments").UpsertId(nestedComments.ID, &nestedComments)
	if err != nil {
		e = appErr{Message: "Could n't update/insert user comment!", Error: err.Error()}
		renderJSON(w, http.StatusBadRequest, e)
		return
	}
	renderJSON(w, http.StatusOK, "Comment posted successfully")
	return
}
