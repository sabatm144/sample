package controller

import (
	"log"
	"net/http"
	"strconv"
	"test_2/server/dbCon"
	"test_2/server/models"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)


func getVoteCount(status  bool, contentID bson.ObjectId) (int) {

	db := dbCon.CopyMongoDB()
	defer db.Session.Close()

	votequery := bson.M{"contentID": contentID, "status": status}	
	count, _ := db.C("votes").Find(votequery).Count()
	return count
}

func formContent(contents []models.Content) {
	for idx := range contents {
		contents[idx].Like = getVoteCount(true, contents[idx].ID)
		contents[idx].DisLike = getVoteCount(false, contents[idx].ID)
	}
}

//Contents returns the contents with paginated
func Contents(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 0)
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 0)

	db := dbCon.CopyMongoDB()
	defer db.Session.Close()

	if limit == 0 {
		limit = 5
	}

	if page == 0 {
		page = 1
	}
	skip := limit * (page - 1)

	contents := []models.Content{}
	db.C("contents").Find(bson.M{}).Skip(int(skip)).Limit(int(limit)).All(&contents)

	total, _ := db.C("contents").Count()

	formContent(contents)
	result := struct {
		Contents    []models.Content `json:"contents"`
		Total       int              `json:"total"`
		CurrentPage int64            `json:"currentPage"`
		Limit       int64            `json:"limit"`
	}{
		contents, total, page, limit,
	}
	
	renderJSON(w, http.StatusOK, result)
}

func GetContent(w http.ResponseWriter, r *http.Request) {

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
	contentID := params.ByName("id")
	if !bson.IsObjectIdHex(contentID) {
		renderJSON(w, http.StatusBadRequest, "Not a valid content ID")
		return
	}
	log.Printf("Content ID %s", contentID)

	contentIns := models.Content{}
	err = db.C("contents").FindId(bson.ObjectIdHex(contentID)).One(&contentIns)
	if err != nil {
		e = appErr{Message: "Content not found", Error: err.Error()}
		renderJSON(w, http.StatusNotFound, e)
		return
	}

	log.Println(contentIns)
	contentIns.Like = getVoteCount(true, contentIns.ID)
	contentIns.DisLike = getVoteCount(false, contentIns.ID)
	renderJSON(w, http.StatusOK, contentIns)
}

func CreateContent(w http.ResponseWriter, r *http.Request) {

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

	ID := r.Context().Value("loggedInUserId").(string)
	if !bson.IsObjectIdHex(ID) {
		renderJSON(w, http.StatusNotFound, "Not a valid user ID")
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
	contentID := params.ByName("id")
	if !bson.IsObjectIdHex(contentID) {
		renderJSON(w, http.StatusNotFound, "Not a valid content ID")
		return
	}
	log.Printf("Content ID %s", contentID)
	contentIns := models.Content{}
	if !parseJSON(w, r.Body, &contentIns) {
		return
	}

	count, _ := db.C("contents").FindId(bson.ObjectIdHex(contentID)).Count()
	if count == 0 {
		e = appErr{Message: "Content not found", Error: err.Error()}
		renderJSON(w, http.StatusNotFound, e)
		return
	}

	err = db.C("contents").Update(bson.M{"_id": contentIns.ID, "userID": contentIns.UserID}, &contentIns)
	if err != nil {
		e = appErr{Message: "Could n't update post!", Error: err.Error()}
		renderJSON(w, http.StatusBadRequest, e)
		return
	}

	renderJSON(w, http.StatusOK, "Post updated successfully!")
}

func DeleteContent(w http.ResponseWriter, r *http.Request) {

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

	contentIns := models.Content{}
	params := r.Context().Value("params").(httprouter.Params)
	contentID := params.ByName("id")
	if !bson.IsObjectIdHex(contentID) {
		renderJSON(w, http.StatusNotFound, "Not a valid content ID")
		return
	}
	log.Printf("Content ID %s", contentID)
	err = db.C("contents").Find(bson.M{"_id": bson.ObjectIdHex(contentID), "userID": userID}).One(&contentIns)
	if err != nil {
		e = appErr{Message: "Could n't find the post specific to user!", Error: err.Error()}
		renderJSON(w, http.StatusOK, e)
		return
	}

	log.Printf("ContentIns ID %s", contentIns)
	err = db.C("contents").RemoveId(contentIns.ID)
	if err != nil {
		e = appErr{Message: "Could n't delete!", Error: err.Error()}
		renderJSON(w, http.StatusOK, e)
		return
	}

	renderJSON(w, http.StatusOK, "Post deleted successfully!")
}
