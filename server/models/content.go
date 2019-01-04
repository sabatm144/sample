package models

import "gopkg.in/mgo.v2/bson"

//Content :User created post struct
type Content struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	UserID      bson.ObjectId `json:"userID" bson:"userID"`
	Link        string        `json:"link" bson:"link"`
	Title       string        `json:"title" bson:"title"`
	Description string        `json:"description" bson:"description"`
	//IsALink: decides whether it's a article link or text post
	IsALink bool `json:"isALink" bson:"isALink"`
}

//Comment :User comment struct
type Comment struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Text      string `json:"text" bson:"text"`
	ContentID bson.ObjectId `json:"contentID" bson:"contentID"`
	UserID    bson.ObjectId `json:"userID" bson:"userID"`
	//IsAParent: decides whether comments has child
	IsAParent bool      `json:"isAParent" bson:"isAParent"`
	Child     []Comment `json:"child" bson:"child"`
}
