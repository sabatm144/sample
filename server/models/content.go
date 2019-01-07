package models

import "gopkg.in/mgo.v2/bson"

//Content :User created post struct
type Content struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	UserID      bson.ObjectId `json:"userID" bson:"userID"`
	Link        string        `json:"link" bson:"link"`
	Title       string        `json:"title" bson:"title"`
	Description string        `json:"description" bson:"description"`
	ContentType string        `json:"contentType" bson:"contentType"`
	Like        int           `json:"like" bson:"-"`
	DisLike     int           `json:"dLike" bson:"-"`
}

//Comment :User comment struct
type Comment struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Text      string        `json:"text,omitempty" bson:"text,omitempty"`
	ContentID bson.ObjectId `json:"contentID,omitempty" bson:"contentID,omitempty"`
	UserID    bson.ObjectId `json:"-" bson:"userID,omitempty"`
	Replies   []Comment     `json:"replies" bson:"replies"`
}
