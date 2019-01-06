package controller

import (
	"fmt"
	"log"
	"net/http"
	"test_2/server/dbCon"
	"test_2/server/models"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

type voteError struct {
	Message string `json:"message"`
	Err     string `json:"error,omitempty"`
}

func (v *voteError) Error() string {
	return fmt.Sprintf("%s: %s", v.Message, v.Err)
}

func getVoteCount(like string, contentID bson.ObjectId) (int, error) {

	db := dbCon.CopyMongoDB()
	defer db.Session.Close()

	votequery := bson.M{"contentID": contentID}	

	if like == "1" {
		//Count
		votequery["status"] = true
		count, _ := db.C("votes").Find(votequery).Count()
		return count, nil
	} else {
		votequery["status"] = false
		count, _ := db.C("votes").Find(votequery).Count()
		return count, nil
	}
}

func CountVotes(w http.ResponseWriter, r *http.Request) {

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

	res := make(map[string]interface{})
	res["noOfLikes"], err = getVoteCount("1", contentID)
	res["noOfDisLikes"], err = getVoteCount("", contentID)
	res["contentID"] = contentID.Hex()
	renderJSON(w, http.StatusOK, res)
	return
}

func LikeContent(w http.ResponseWriter, r *http.Request) {

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