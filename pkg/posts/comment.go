package posts

import (
	"reddit/pkg/user"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Comment struct {
	ID      bson.ObjectId `json:"id,string" bson:"_id"`
	Autor   *user.User    `json:"author" bson:"autor"`
	Body    string        `json:"body" bson:"body"`
	Created string        `json:"created" bson:"created"`
}

type CommentsRepo struct {
	DB *mgo.Collection
}
