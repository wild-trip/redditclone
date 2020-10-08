package posts

import (
	"reddit/pkg/user"

	"gopkg.in/mgo.v2/bson"
)

type Post struct {
	Author           *user.User      `bson:"author"`
	Category         string          `bson:"category"`
	CommentsID       []bson.ObjectId `bson:"comments"`
	Created          string          `bson:"created"`
	ID               bson.ObjectId   `bson:"_id"`
	Score            int             `bson:"score"`
	Text             string          `bson:"text,omitempty"`
	Link             string          `bson:"url,omitempty"`
	Title            string          `bson:"title"`
	Type             string          `bson:"type"`
	UpvotePercentage int             `bson:"upvotePercentage"`
	Views            int             `bson:"views"`
	Votes            []Vote          `bson:"votes"` // userID, vote=1,-1
}

type Vote struct {
	UserID int64 `json:"user,string" bson:"user"`
	Rating int   `json:"vote" bson:"vote"`
}
