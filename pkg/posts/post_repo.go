package posts

import (
	"fmt"
	"log"
	"reddit/pkg/user"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PostsRepo struct {
	DB PostRepositoryDBInterface
}

type FindInterface interface {
	One(interface{}) error
	All(interface{}) error
}

type PostRepositoryDBInterface interface {
	Find(interface{}) *mgo.Query // FindInterface
	Insert(...interface{}) error
	Update(interface{}, interface{}) error
	Remove(interface{}) error
}

func NewRepo(collection PostRepositoryDBInterface) *PostsRepo {
	return &PostsRepo{DB: collection}
}

func (repo *PostsRepo) GetAll() ([]*Post, error) {
	posts := []*Post{}
	err := repo.DB.Find(bson.M{}).All(&posts)
	if err != nil {
		log.Printf("DB error")
		return nil, fmt.Errorf("DB err: %v", err)
	}
	return posts, nil
}

func (repo *PostsRepo) GetCategory(category string) ([]*Post, error) {
	var posts []*Post
	err := repo.DB.Find(bson.M{"category": category}).All(&posts)
	if err != nil {
		log.Printf("DB error")
		return nil, err
	}
	return posts, nil
}

func (repo *PostsRepo) GetByID(id bson.ObjectId) (*Post, error) {
	var post *Post
	err := repo.DB.Find(bson.M{"_id": id}).One(&post)
	if err != nil {
		log.Printf("DB error: %v", err)
		return nil, err
	}
	return post, nil
}

func (repo *PostsRepo) GetByUserLogin(login string) ([]*Post, error) {
	var posts []*Post
	err := repo.DB.Find(bson.M{"author.username": login}).All(&posts)
	if err != nil {
		log.Printf("DB error")
		return nil, err
	}
	return posts, nil
}

func (repo *PostsRepo) Add(
	user *user.User,
	category string,
	title string,
	typePost string,
	text string,
	link string) (*Post, error) {
	nowTime := time.Now()
	timestamp := nowTime.Format(time.RFC3339)
	newPost := &Post{
		ID:               bson.NewObjectId(),
		Author:           user,
		Category:         category,
		CommentsID:       make([]bson.ObjectId, 0),
		Created:          timestamp,
		Score:            1,
		Title:            title,
		UpvotePercentage: 100,
		Views:            0,
		Votes:            make([]Vote, 0),
	}
	if typePost == "link" {
		newPost.Link = link
		newPost.Text = ""
		newPost.Type = typePost
	} else {
		newPost.Link = ""
		newPost.Text = text
		newPost.Type = typePost
	}
	newPost.Votes = append(newPost.Votes, Vote{
		UserID: user.ID,
		Rating: 1,
	})

	err := repo.DB.Insert(&newPost)
	if err != nil {
		log.Printf("Insert error")
		return nil, err
	}
	return newPost, nil
}

func (repo *PostsRepo) AddComment(postID, commentID bson.ObjectId) (*Post, error) {
	post, err := repo.GetByID(postID)
	if err != nil {
		return nil, fmt.Errorf("Error getting post from BD: %v", err)
	}
	post.CommentsID = append(post.CommentsID, commentID)
	err = repo.DB.Update(
		bson.M{"_id": postID},
		post)
	if err != nil {
		return nil, fmt.Errorf("Error update BD: %v", err)
	}
	return post, nil
}

func (repo *PostsRepo) UpViews(postID bson.ObjectId) error {
	post, err := repo.GetByID(postID)
	if err != nil {
		return fmt.Errorf("Error getting post from BD: %v", err)
	}
	post.Views++
	err = repo.DB.Update(
		bson.M{"_id": postID},
		post)
	if err != nil {
		return fmt.Errorf("Error update BD: %v", err)
	}
	return nil
}

func (repo *PostsRepo) DeleteComment(postID, commentID bson.ObjectId) (*Post, error) {
	post, err := repo.GetByID(postID)
	if err != nil {
		return nil, fmt.Errorf("Error getting post from BD: %v", err)
	}

	for iComment := range post.CommentsID {
		if post.CommentsID[iComment] == commentID {
			post.CommentsID = append(post.CommentsID[:iComment], post.CommentsID[iComment+1:]...)
			break
		}
	}

	err = repo.DB.Update(
		bson.M{"_id": postID},
		post)
	if err != nil {
		return nil, fmt.Errorf("Error update BD: %v", err)
	}
	return post, nil
}

func (repo *PostsRepo) Upvote(user *user.User, postID bson.ObjectId) (*Post, error) {
	elem, err := repo.GetByID(postID)
	if err != nil {
		return nil, fmt.Errorf("DB err: %v", err)
	}
	nUpVotes := 0
	iUserVote := -1
	for i, vote := range elem.Votes {
		if vote.UserID == user.ID {
			iUserVote = i
		}
		if vote.Rating == 1 {
			nUpVotes++
		}
	}
	if iUserVote != -1 {
		if elem.Votes[iUserVote].Rating == -1 {
			nUpVotes++
			elem.Score += 2
			elem.Votes[iUserVote].Rating = 1
		}
	} else {
		nUpVotes += 1
		elem.Score++
		elem.Votes = append(elem.Votes, Vote{
			UserID: user.ID,
			Rating: 1,
		})
	}

	elem.UpvotePercentage = nUpVotes * 100 / len(elem.Votes)

	err = repo.DB.Update(
		bson.M{"_id": postID},
		elem)
	if err != nil {
		return nil, fmt.Errorf("Error update BD: %v", err)
	}
	return elem, nil
}

func (repo *PostsRepo) Downvote(user *user.User, postID bson.ObjectId) (*Post, error) {
	elem, err := repo.GetByID(postID)
	if err != nil {
		return nil, fmt.Errorf("DB err: %v", err)
	}
	nUpVotes := 0
	iUserVote := -1
	for i, vote := range elem.Votes {
		if vote.UserID == user.ID {
			iUserVote = i
		}
		if vote.Rating == 1 {
			nUpVotes++
		}
	}
	if iUserVote != -1 {
		if elem.Votes[iUserVote].Rating == 1 {
			elem.Score -= 2
			nUpVotes--
			elem.Votes[iUserVote].Rating = -1
		}
	} else {
		elem.Score--
		elem.Votes = append(elem.Votes, Vote{
			UserID: user.ID,
			Rating: -1,
		})
	}

	elem.UpvotePercentage = nUpVotes * 100 / len(elem.Votes)
	err = repo.DB.Update(
		bson.M{"_id": postID},
		elem)
	if err != nil {
		return nil, fmt.Errorf("Error update BD: %v", err)
	}
	return elem, nil
}

func (repo *PostsRepo) Delete(postID bson.ObjectId) (bool, error) {
	err := repo.DB.Remove(bson.M{"_id": postID})
	if err == mgo.ErrNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
