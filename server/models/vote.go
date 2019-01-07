package models

import "gopkg.in/mgo.v2/bson"

//Vote struct for capturing likes and dislikes specific to user for a content 
type Vote struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	ContentID bson.ObjectId `json:"contentID" bson:"contentID"`
	UserID    bson.ObjectId `json:"userID" bson:"userID"`
	Status    bool          `json:"status" bson:"status"`
}
