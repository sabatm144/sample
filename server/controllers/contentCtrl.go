package controller

import (
	"log"
	"net/http"
	"test_2/server/dbCon"
	"test_2/server/models"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

func ListContents(w http.ResponseWriter, r *http.Request) {

	db := dbCon.CopyMongoDB()
	defer db.Session.Close()

	e := appErr{}
	contents := []models.Content{}
	err := db.C("contents").Find(bson.M{}).All(&contents)
	if err != nil {
		e = appErr{Message: "Contents not found", Error: err.Error()}
		renderJSON(w, http.StatusNotFound, e)
		return
	}

	log.Println(contents)
	renderJSON(w, http.StatusOK, contents)
}

func GetContent(w http.ResponseWriter, r *http.Request) {

	userID := bson.ObjectIdHex(r.Context().Value("loggedInUserId").(string))
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
	id := params.ByName("id")
	log.Printf("Content ID %s", id)

	contentIns := models.Content{}
	err = db.C("contents").FindId(bson.ObjectIdHex(id)).One(&contentIns)
	if err != nil {
		e = appErr{Message: "Content not found", Error: err.Error()}
		renderJSON(w, http.StatusNotFound, e)
		return
	}

	log.Println(contentIns)
	renderJSON(w, http.StatusOK, contentIns)
}

func CreateContent(w http.ResponseWriter, r *http.Request) {

	userID := bson.ObjectIdHex(r.Context().Value("loggedInUserId").(string))
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

	contentIns := models.Content{}
	if !parseJSON(w, r.Body, &contentIns) {
		return
	}
	contentIns.ID = bson.NewObjectId()
	contentIns.UserID = userIns.ID
	//Create post for the specific user
	err = db.C("contents").Insert(contentIns)
	if err != nil {
		e = appErr{Message: "Post not created", Error: err.Error()}
		renderJSON(w, http.StatusBadRequest, e)
		return
	}

	e = appErr{Message: "Created!"}
	renderJSON(w, http.StatusOK, contentIns)
}

func UpdateContent(w http.ResponseWriter, r *http.Request) {

	userID := bson.ObjectIdHex(r.Context().Value("loggedInUserId").(string))
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
	id := params.ByName("id")
	log.Printf("Content ID %s", id)

	contentIns := models.Content{}
	err = db.C("contents").FindId(bson.ObjectIdHex(id)).One(&contentIns)
	if err != nil {
		e = appErr{Message: "User not found", Error: err.Error()}
		renderJSON(w, http.StatusNotFound, e)
		return
	}

	if !parseJSON(w, r.Body, &contentIns) {
		return
	}

	log.Println(contentIns)
	err = db.C("contents").Update(bson.M{"_id": contentIns.ID, "userID": contentIns.UserID}, &contentIns)
	if err != nil {
		e = appErr{Message: "Could n't update user post!", Error: err.Error()}
		renderJSON(w, http.StatusBadRequest, e)
		return
	}

	e = appErr{Message: "Updated!"}
	renderJSON(w, http.StatusOK, e)
}

func DeleteContent(w http.ResponseWriter, r *http.Request) {

	userID := bson.ObjectIdHex(r.Context().Value("loggedInUserId").(string))
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

	contentIns := models.Content{}
	params := r.Context().Value("params").(httprouter.Params)
	id := params.ByName("id")

	log.Printf("Content ID %s", id)
	err = db.C("contents").Find(bson.M{"_id": bson.ObjectIdHex(id), "userID": userID}).One(&contentIns)
	if err != nil {
		e = appErr{Message: "Could n't find the post specific to user!", Error: err.Error()}
		renderJSON(w, http.StatusOK, e)
		return
	}

	log.Printf("Content Ins ID %s", contentIns)

	err = db.C("contents").RemoveId(contentIns.ID)
	if err != nil {
		e = appErr{Message: "Could n't delete!", Error: err.Error()}
		renderJSON(w, http.StatusOK, e)
		return
	}

	e = appErr{Message: "Deleted!"}
	renderJSON(w, http.StatusOK, e)
}

func ListComments(w http.ResponseWriter, r *http.Request) {

	if bson.IsObjectIdHex(r.Context().Value("loggedInUserId").(string)) {
		renderJSON(w, http.StatusBadRequest, "Not a valid user ID")
		return	
	}

	userID := bson.ObjectIdHex(r.Context().Value("loggedInUserId").(string))
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
	if bson.IsObjectIdHex(params.ByName("id")) {
		renderJSON(w, http.StatusBadRequest, "Not a valid ID")
		return	
	}

	contentID := bson.ObjectIdHex(params.ByName("id"))
	log.Printf("Content ID %s %s", contentID)

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

	userID := bson.ObjectIdHex(r.Context().Value("loggedInUserId").(string))
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
