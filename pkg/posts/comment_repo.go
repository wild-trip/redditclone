package posts

import (
	"log"
	"reddit/pkg/user"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func NewCommentRepo(collection *mgo.Collection) *CommentsRepo {
	return &CommentsRepo{DB: collection}
}

func (repo *CommentsRepo) NewComment(autor *user.User, body string) (bson.ObjectId, error) {
	newCommment := &Comment{
		ID:      bson.NewObjectId(),
		Autor:   autor,
		Body:    body,
		Created: time.Now().Format(time.RFC3339),
	}
	repo.DB.Insert(&newCommment)
	return newCommment.ID, nil
}

func (repo *CommentsRepo) GetByID(id bson.ObjectId) (*Comment, error) {
	var comment *Comment
	err := repo.DB.Find(bson.M{"_id": id}).One(&comment)
	if err != nil {
		log.Printf("DB error")
		return nil, err
	}
	return comment, nil
}

func (repo *CommentsRepo) DelComment(commentID bson.ObjectId) (bool, error) {
	err := repo.DB.Remove(bson.M{"_id": commentID})
	if err != nil {
		return false, err
	}
	return true, nil
}
