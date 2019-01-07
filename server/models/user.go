package models

import "gopkg.in/mgo.v2/bson"

type User struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	Name     string        `json:"name" bson:"name"`
	EmailID  string        `json:"emailID" bson:"emailID"`
	Password string        `json:"-" bson:"password"`
}

// Login provides struct for user credentails
type Login struct {
	Name     string `json:"name"`
	EmailID  string `json:"emailID"`
	Password string `json:"password"`
}

//LoggedInUser provides struct for capture logged-in user data
type LoggedInUser struct {
	ID      bson.ObjectId `json:"id" bson:"_id"`
	Name    string        `json:"firstName" bson:"firstName"`
	EmailID string        `json:"emailID" bson:"emailID"`
}
