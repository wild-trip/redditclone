package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"reddit/pkg/session"
	"reddit/pkg/user"

	"reddit/pkg/posts"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
)

type PostsRepositoryInterface interface {
	GetAll() ([]*posts.Post, error)
	GetCategory(string) ([]*posts.Post, error)
	GetByID(bson.ObjectId) (*posts.Post, error)
	GetByUserLogin(string) ([]*posts.Post, error)
	Add(*user.User, string, string, string, string, string) (*posts.Post, error)
	AddComment(bson.ObjectId, bson.ObjectId) (*posts.Post, error)
	UpViews(bson.ObjectId) error
	DeleteComment(bson.ObjectId, bson.ObjectId) (*posts.Post, error)
	Upvote(*user.User, bson.ObjectId) (*posts.Post, error)
	Downvote(*user.User, bson.ObjectId) (*posts.Post, error)
	Delete(bson.ObjectId) (bool, error)
}
type CommentsRepositoryInterface interface {
	NewComment(*user.User, string) (bson.ObjectId, error)
	GetByID(bson.ObjectId) (*posts.Comment, error)
	DelComment(bson.ObjectId) (bool, error)
}

type PostsHandler struct {
	Tmpl        *template.Template
	PostsRepo   PostsRepositoryInterface
	CommentRepo CommentsRepositoryInterface
	Logger      *zap.SugaredLogger
}

type NewPostRequest struct {
	Category string `json:"category"`
	Text     string `json:"text,omitempty"`
	Link     string `json:"url,omitempty"`
	Title    string `json:"title"`
	Type     string `json:"type"`
}

type PostResponse struct {
	Author           *user.User       `json:"author"`
	Category         string           `json:"category"`
	Comments         []*posts.Comment `json:"comments"`
	Created          string           `json:"created"`
	ID               bson.ObjectId    `json:"id,string"`
	Score            int              `json:"score"`
	Text             string           `json:"text,omitempty"`
	Link             string           `json:"url,omitempty"`
	Title            string           `json:"title"`
	Type             string           `json:"type"`
	UpvotePercentage int              `json:"upvotePercentage"`
	Views            int              `json:"views"`
	Votes            []posts.Vote     `json:"votes"`
}

func PostToPostResponse(post *posts.Post, commentsRepo CommentsRepositoryInterface) (*PostResponse, error) {
	comments := make([]*posts.Comment, 0)
	for _, commentID := range post.CommentsID {
		comment, err := commentsRepo.GetByID(commentID)
		if err != nil {
			return nil, fmt.Errorf("Can't get comment: %v", err)
		}
		comments = append(comments, comment)
	}
	postResponse := &PostResponse{
		Author:           post.Author,
		Category:         post.Category,
		Comments:         comments,
		Created:          post.Created,
		ID:               post.ID,
		Score:            post.Score,
		Text:             post.Text,
		Link:             post.Link,
		Title:            post.Title,
		Type:             post.Type,
		UpvotePercentage: post.UpvotePercentage,
		Views:            post.Views,
		Votes:            post.Votes,
	}
	return postResponse, nil
}

var TemplateName = "index.html"

func (h *PostsHandler) Init(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, TemplateName, nil)
	if err != nil {
		http.Error(w, `Template error`, http.StatusInternalServerError)
		h.Logger.Errorf("Template error: %v", err)
		return
	}
	h.Logger.Infof("Main page")
}

func (h *PostsHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	posts, err := h.PostsRepo.GetAll()
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		h.Logger.Errorf("DB err: %v", err)
		return
	}

	postsResponse := make([]*PostResponse, 0)

	for _, post := range posts {
		postResponse, err := PostToPostResponse(post, h.CommentRepo)
		if err != nil {
			http.Error(w, "BD Error", http.StatusInternalServerError)
			h.Logger.Errorf("Post Transform error: %v", err)
			return
		}
		postsResponse = append(postsResponse, postResponse)
	}

	resp, _ := json.Marshal(postsResponse)
	w.Write(resp)
	h.Logger.Infof("List all")
}

func (h *PostsHandler) ListCategory(w http.ResponseWriter, r *http.Request) {
	arg := mux.Vars(r)
	cat, _ := arg["CATEGORY"]
	posts, err := h.PostsRepo.GetCategory(cat)
	if err != nil {
		http.Error(w, `DB error`, http.StatusInternalServerError)
		h.Logger.Errorf("Bad category: %v", err)
		return
	}

	postsResponse := make([]*PostResponse, 0)

	for _, post := range posts {
		postResponse, err := PostToPostResponse(post, h.CommentRepo)
		if err != nil {
			http.Error(w, `DB error`, http.StatusInternalServerError)
			h.Logger.Errorf("Post Transform error: %v", err)
			return
		}
		postsResponse = append(postsResponse, postResponse)
	}

	resp, _ := json.Marshal(postsResponse)
	w.Write(resp)
	h.Logger.Infof("List category")
}

func (h *PostsHandler) ListByUserLogin(w http.ResponseWriter, r *http.Request) {
	arg := mux.Vars(r)
	userLogin, _ := arg["USER_LOGIN"]
	posts, err := h.PostsRepo.GetByUserLogin(userLogin)
	if err != nil {
		http.Error(w, `DB error`, http.StatusInternalServerError)
		h.Logger.Errorf("Bad login: %v", err)
		return
	}

	postsResponse := make([]*PostResponse, 0)

	for _, post := range posts {
		postResponse, err := PostToPostResponse(post, h.CommentRepo)
		if err != nil {
			http.Error(w, `DB error`, http.StatusInternalServerError)
			h.Logger.Errorf("Post Transform error: %v", err)
			return
		}
		postsResponse = append(postsResponse, postResponse)
	}

	resp, _ := json.Marshal(postsResponse)
	w.Write(resp)
	h.Logger.Infof("List posts by user ligin")
}

func (h *PostsHandler) ListByID(w http.ResponseWriter, r *http.Request) {
	arg := mux.Vars(r)
	if !bson.IsObjectIdHex(arg["ID"]) {
		http.Error(w, "Bad ID", http.StatusInternalServerError)
		h.Logger.Errorf("Bad ID: %v")
		return
	}
	postID := bson.ObjectIdHex(arg["ID"])
	post, err := h.PostsRepo.GetByID(postID)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		h.Logger.Errorf("DB err: %v", err)
		return
	}
	err = h.PostsRepo.UpViews(postID)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		h.Logger.Errorf("Up view err: %v", err)
		return
	}
	postResponse, err := PostToPostResponse(post, h.CommentRepo)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		h.Logger.Errorf("Post Transform error: %v", err)
		return
	}

	resp, _ := json.Marshal(postResponse)
	w.Write(resp)
	h.Logger.Infof("List posts by post id")
}

func (h *PostsHandler) Add(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, `Bad auth`, http.StatusBadRequest)
		h.Logger.Errorf("Bad auth. Error: %v", err)
		return
	}
	newRequest := new(NewPostRequest)
	body, errReadBody := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, newRequest)
	if errReadBody != nil || err != nil {
		http.Error(w, "Internal error", http.StatusBadRequest)
		h.Logger.Errorf("Bad JSON. Error: %v, %v", errReadBody, err)
		return
	}

	newPost, err := h.PostsRepo.Add(
		sess.User,
		newRequest.Category,
		newRequest.Title,
		newRequest.Type,
		newRequest.Text,
		newRequest.Link,
	)
	if err != nil {
		http.Error(w, "BD error", http.StatusInternalServerError)
		h.Logger.Errorf("Bad add post. Error: %v, %v", err)
		return
	}

	postResponse, err := PostToPostResponse(newPost, h.CommentRepo)
	if err != nil {
		http.Error(w, `BD error`, http.StatusInternalServerError)
		h.Logger.Errorf("Post Transform error: %v", err)
		return
	}

	answer, _ := json.Marshal(postResponse)
	w.Write(answer)
	h.Logger.Infof("New post was created with ID: %v", newPost.ID)
}

func (h *PostsHandler) Upvote(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, `Bad auth`, http.StatusBadRequest)
		h.Logger.Errorf("Bad auth. Error: %v", err)
		return
	}
	arg := mux.Vars(r)
	if !bson.IsObjectIdHex(arg["POST_ID"]) {
		http.Error(w, "Bad id", http.StatusBadRequest)
		h.Logger.Errorf("Bad post id")
		return
	}
	postID := bson.ObjectIdHex(arg["POST_ID"])
	post, err := h.PostsRepo.Upvote(sess.User, postID)
	if err != nil {
		http.Error(w, `Bad upvote`, http.StatusInternalServerError)
		h.Logger.Errorf("Bad upvote. Error: %v", err)
		return
	}

	postResponse, err := PostToPostResponse(post, h.CommentRepo)
	if err != nil {
		http.Error(w, `Internal error`, http.StatusInternalServerError)
		h.Logger.Errorf("Post Transform error: %v", err)
		return
	}

	resp, _ := json.Marshal(postResponse)
	w.Write(resp)
	h.Logger.Infof("Upvote post")
}

func (h *PostsHandler) Downvote(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, `Bad auth`, http.StatusBadRequest)
		h.Logger.Errorf("Bad auth. Error: %v", err)
		return
	}
	arg := mux.Vars(r)
	if !bson.IsObjectIdHex(arg["POST_ID"]) {
		http.Error(w, "Bad id", http.StatusBadRequest)
		h.Logger.Errorf("Bad post id")
		return
	}
	postID := bson.ObjectIdHex(arg["POST_ID"])
	post, err := h.PostsRepo.Downvote(sess.User, postID)
	if err != nil {
		http.Error(w, `Bad downvote`, http.StatusInternalServerError)
		h.Logger.Errorf("Bad downvote. Error: %v", err)
		return
	}

	postResponse, err := PostToPostResponse(post, h.CommentRepo)
	if err != nil {
		http.Error(w, `Internal error`, http.StatusInternalServerError)
		h.Logger.Errorf("Post Transform error: %v", err)
		return
	}

	resp, _ := json.Marshal(postResponse)
	w.Write(resp)
	h.Logger.Infof("Downvote post")
}

func (h *PostsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	_, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, `Bad auth`, http.StatusBadRequest)
		h.Logger.Errorf("Bad auth. Error: %v", err)
		return
	}
	arg := mux.Vars(r)
	if !bson.IsObjectIdHex(arg["POST_ID"]) {
		http.Error(w, "Bad id", http.StatusBadRequest)
		return
	}
	postID := bson.ObjectIdHex(arg["POST_ID"])
	ok, err := h.PostsRepo.Delete(postID)
	if err != nil {
		http.Error(w, `Delete error`, http.StatusInternalServerError)
		h.Logger.Errorf("Delete error: %v", err)
		return
	}
	if ok {
		w.Write([]byte("{\"message\": \"success\"}"))
		h.Logger.Infof("Delete post success")
	} else {
		w.Write([]byte("{\"message\": \"failure\"}"))
		h.Logger.Infof("Delete post failure")
	}
}
