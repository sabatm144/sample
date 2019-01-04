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

type voteError struct {
	Message string `json:"message"`
	Err     string `json:"error,omitempty"`
}

func (v *voteError) Error() string {
	return fmt.Sprintf("%s: %s", v.Message, v.Err)
}

func getVoteCount(like string, userID, contentID bson.ObjectId) (int, error) {

	db := dbCon.CopyMongoDB()
	defer db.Session.Close()

	voterIns := models.Voter{}
	votequery := bson.M{"userID": userID, "contentID": contentID, "status": false}
	e := &voteError{}

	if like == "1" {
		err := db.C("votes").Find(votequery).One(&voterIns)
		if err != nil {
			e.Message = "Voter not found"
			e.Err = err.Error()
			return 0, e
		}

		if voterIns.ID.Hex() != "" {
			voterIns.Status = true
			err = db.C("votes").UpdateId(voterIns.ID, &voterIns)
			if err != nil {
				e.Message = "Voter not updated"
				e.Err = err.Error()
				return 0, e
			}
		}

		votequery["status"] = true
		count, countErr := db.C("votes").Find(votequery).Count()
		if countErr != nil {
			e.Message = "Could count Votes"
			e.Err = err.Error()
			return 0, e
		}

		return count, nil
	} else {
		votequery["status"] = true
		err := db.C("votes").Find(votequery).One(&voterIns)
		if err != nil {
			e.Message = "Voter not found"
			e.Err = err.Error()
			return 0, e
		}

		if voterIns.ID.Hex() != "" {
			voterIns.Status = false
			err = db.C("votes").UpdateId(voterIns.ID, &voterIns)
			if err != nil {
				e.Message = "Vote not updated"
				e.Err = err.Error()
				return 0, err
			}
		}

		votequery["status"] = false
		count, countErr := db.C("votes").Find(votequery).Count()
		if countErr != nil {
			e.Message = "Could not count votes"
			e.Err = err.Error()
			return 0, e
		}

		return count, nil
	}
}

func VoteContent(w http.ResponseWriter, r *http.Request) {

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
	like := params.ByName("like")
	log.Printf("Content ID %s %s", contentID, like)

	res := make(map[string]int)
	res["noOfLikes"], err = getVoteCount(like, userID, contentID)
	res["noOfDisLikes"], err = getVoteCount(like, userID, contentID)
	renderJSON(w, http.StatusNotFound, res)
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
	log.Printf("Content ID %s %s", contentID)
	type commentDesc struct {
		Text      string        `json:"text" bson:"text"`
		CommentID bson.ObjectId `json:"commentID" json:"commentID"`
	}
	commentIns := commentDesc{}

	if !parseJSON(w, r.Body, &commentIns) {
		return
	}

	nestedComments := &models.Comment{}
	if commentIns.CommentID.Hex() != "" {
		nestedComments.IsAParent = true
		err := db.C("comments").Find(bson.M{"_id": commentIns.CommentID, "contentID": contentID}).One(nestedComments)
		if err != nil {
			e = appErr{Message: "Could n't update user post!", Error: err.Error()}
			renderJSON(w, http.StatusBadRequest, e)
			return
		}

		if nestedComments.ID.Hex() != "" {
			nestedComments.IsAParent = true
			childComments := models.Comment{}
			childComments.ID = bson.NewObjectId()
			childComments.ContentID = contentID
			childComments.UserID = userID
			childComments.Text = commentIns.Text
			nestedComments.Child = append(nestedComments.Child, childComments)
		}

		if nestedComments.ID.Hex() == "" {
			nestedComments.ID = bson.NewObjectId()
			nestedComments.ContentID = contentID
			nestedComments.UserID = userID
			nestedComments.Text = commentIns.Text
		}
	}

	renderJSON(w, http.StatusNotFound, nestedComments)
	return
}
