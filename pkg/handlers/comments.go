package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reddit/pkg/session"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

type AddCommentRequest struct {
	Comment string `json:"comment"`
}

func (h *PostsHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, `Bad auth`, http.StatusBadRequest)
		h.Logger.Errorf("Bad auth. Error: %v", err)
		return
	}
	arg := mux.Vars(r)

	newRequest := new(AddCommentRequest)
	body, errReadBody := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, newRequest)
	if errReadBody != nil || err != nil {
		http.Error(w, "", http.StatusBadRequest)
		h.Logger.Errorf("Bad JSON. Error: %v, %v", errReadBody, err)
		return
	}
	if !bson.IsObjectIdHex(arg["POST_ID"]) {
		http.Error(w, "bad id", 500)
		return
	}
	postID := bson.ObjectIdHex(arg["POST_ID"])
	commentID, err := h.CommentRepo.NewComment(sess.User, newRequest.Comment)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		h.Logger.Errorf("Bad add comment to comment repo. Error: %v, %v", err)
		return
	}
	post, err := h.PostsRepo.AddComment(postID, commentID)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		h.Logger.Errorf("Bad add comment to post repo. Error: %v, %v", err)
		return
	}
	postResponse, err := PostToPostResponse(post, h.CommentRepo)
	if err != nil {
		http.Error(w, ``, http.StatusInternalServerError)
		h.Logger.Errorf("Post Transform error: %v", err)
		return
	}
	answer, errAnswer := json.Marshal(postResponse)
	if errAnswer != nil {
		http.Error(w, `Bad form`, http.StatusBadRequest)
		h.Logger.Errorf("Bad answer JSON. Error: %v", err)
		return
	}
	w.Write(answer)
	h.Logger.Infof("Post %v was updated by new comment", post.ID)
}

func (h *PostsHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	_, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, `Bad auth`, http.StatusBadRequest)
		h.Logger.Errorf("Bad auth. Error: %v", err)
		return
	}
	arg := mux.Vars(r)

	if !bson.IsObjectIdHex(arg["POST_ID"]) || !bson.IsObjectIdHex(arg["COMMENT_ID"]) {
		http.Error(w, "bad id", 500)
		return
	}
	postID := bson.ObjectIdHex(arg["POST_ID"])
	commentID := bson.ObjectIdHex(arg["COMMENT_ID"])
	isDelete, err := h.CommentRepo.DelComment(commentID)
	if err != nil || !isDelete {
		h.Logger.Errorf("Delete comment fall, %v", err)
		return
	}

	post, err := h.PostsRepo.DeleteComment(postID, commentID)

	postResponse, err := PostToPostResponse(post, h.CommentRepo)
	if err != nil {
		http.Error(w, "BD Error", http.StatusInternalServerError)
		h.Logger.Errorf("Post Transform error: %v", err)
		return
	}

	answer, errAnswer := json.Marshal(postResponse)
	if errAnswer != nil {
		http.Error(w, `Bad form`, http.StatusBadRequest)
		h.Logger.Errorf("Bad answer JSON. Error: %v", err)
		return
	}
	w.Write(answer)
	h.Logger.Infof("Post %v was udated by delete comment %v", post.ID, commentID)
}
