package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	posts "reddit/pkg/posts"
	"reddit/pkg/session"
	user "reddit/pkg/user"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
)

var testPost1 = &posts.Post{
	Author: &user.User{
		ID:       1,
		Username: "rvasily",
	},
	Category: "music",
	CommentsID: []bson.ObjectId{
		"^\xba\xf9\xf2<\x04\xc1|V\xf5\x12E",
	},
	Created:          "2020-05-12T22:33:02+03:00",
	ID:               "^\xba\xf9\xee<\x04\xc1|V\xf5\x12D",
	Score:            1,
	Text:             "Something 1",
	Link:             "",
	Title:            "Lorem",
	Type:             "text",
	UpvotePercentage: 100,
	Views:            10,
	Votes: []posts.Vote{
		{
			UserID: 1,
			Rating: 1,
		},
	},
}
var testUser = &user.User{
	ID:       1,
	Username: "rvasily",
}
var testComment = &posts.Comment{
	ID:      "^\xba\xf9\xf2<\x04\xc1|V\xf5\x12E",
	Autor:   testUser,
	Body:    "something",
	Created: "2020-05-12T22:33:02+03:00",
}

type Result struct {
	Body []byte
	Code int
}

type TestCase struct {
	Request        *http.Request
	ExpectMockFunc []*gomock.Call
	ReturnMockFunc [][]interface{}
	HandlerFunc    func(w http.ResponseWriter, r *http.Request)
	ExpectResult   Result
}

func TestListAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostsRepo := NewMockPostsRepositoryInterface(ctrl)
	mockCommentsRepo := NewMockCommentsRepositoryInterface(ctrl)
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	postsTestHandler := &PostsHandler{
		Tmpl:        template.Must(template.ParseFiles("../../template/index.html")),
		Logger:      logger,
		PostsRepo:   mockPostsRepo,
		CommentRepo: mockCommentsRepo,
	}
	testUser := &user.User{
		ID:       1,
		Username: "rvasily",
	}
	//password := "lovelove"
	//testPost1JSON, _ := json.Marshal([]*posts.Post{testPost1})
	testPosts := []*posts.Post{testPost1}

	testCases := []TestCase{
		{ //Init SUCCESS
			Request:        httptest.NewRequest("GET", "/", nil),
			ExpectMockFunc: []*gomock.Call{},
			ReturnMockFunc: [][]interface{}{},
			HandlerFunc:    postsTestHandler.Init,
			ExpectResult: Result{
				Body: func() []byte {
					w := httptest.NewRecorder()
					postsTestHandler.Tmpl.ExecuteTemplate(w, "index.html", nil)
					body, _ := ioutil.ReadAll(w.Result().Body)
					return body
				}(),
				Code: http.StatusOK,
			},
		},
		{ //List All SUCCESS
			Request: httptest.NewRequest("GET", "/api/posts/", nil),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().GetAll(),
				mockCommentsRepo.EXPECT().GetByID(gomock.Any()),
			},
			ReturnMockFunc: [][]interface{}{
				{testPosts, nil},
				{testComment, nil},
			},
			HandlerFunc: postsTestHandler.ListAll,
			ExpectResult: Result{
				Body: []byte(`[{"author":{"username":"rvasily","id":"1"},"category":"music","comments":[{"id":"5ebaf9f23c04c17c56f51245","author":{"username":"rvasily","id":"1"},"body":"something","created":"2020-05-12T22:33:02+03:00"}],"created":"2020-05-12T22:33:02+03:00","id":"5ebaf9ee3c04c17c56f51244","score":1,"text":"Something 1","title":"Lorem","type":"text","upvotePercentage":100,"views":10,"votes":[{"user":"1","vote":1}]}]`),
				Code: http.StatusOK,
			},
		},
		{ //List ALL. Post repo error
			Request: httptest.NewRequest("GET", "/api/posts/", nil),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().GetAll(),
			},
			ReturnMockFunc: [][]interface{}{
				{nil, fmt.Errorf("Internal error")},
			},
			HandlerFunc: postsTestHandler.ListAll,
			ExpectResult: Result{
				Body: []byte("DB err\n"),
				Code: http.StatusInternalServerError,
			},
		},
		{ //List All. Comment repo error
			Request: httptest.NewRequest("GET", "/api/posts/", nil),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().GetAll(),
				mockCommentsRepo.EXPECT().GetByID(gomock.Any()),
			},
			ReturnMockFunc: [][]interface{}{
				{testPosts, nil},
				{nil, fmt.Errorf("Internal error")},
			},
			HandlerFunc: postsTestHandler.ListAll,
			ExpectResult: Result{
				Body: []byte("BD Error\n"),
				Code: http.StatusInternalServerError,
			},
		},
		{ //ListCategory SUCCESS
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/posts/{CATEGORY}", nil)
				req := mux.SetURLVars(r, map[string]string{
					"CATEGORY": "music",
				})
				return req
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().GetCategory("music"),
				mockCommentsRepo.EXPECT().GetByID(gomock.Any()),
			},
			ReturnMockFunc: [][]interface{}{
				{testPosts, nil},
				{testComment, nil},
			},
			HandlerFunc: postsTestHandler.ListCategory,
			ExpectResult: Result{
				Body: []byte(`[{"author":{"username":"rvasily","id":"1"},"category":"music","comments":[{"id":"5ebaf9f23c04c17c56f51245","author":{"username":"rvasily","id":"1"},"body":"something","created":"2020-05-12T22:33:02+03:00"}],"created":"2020-05-12T22:33:02+03:00","id":"5ebaf9ee3c04c17c56f51244","score":1,"text":"Something 1","title":"Lorem","type":"text","upvotePercentage":100,"views":10,"votes":[{"user":"1","vote":1}]}]`),
				Code: http.StatusOK,
			},
		},
		{ //ListCategory. Post repo error
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/posts/{CATEGORY}", nil)
				req := mux.SetURLVars(r, map[string]string{
					"CATEGORY": "music",
				})
				return req
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().GetCategory("music"),
			},
			ReturnMockFunc: [][]interface{}{
				{nil, fmt.Errorf("Internal error")},
			},
			HandlerFunc: postsTestHandler.ListCategory,
			ExpectResult: Result{
				Body: []byte("DB error\n"),
				Code: http.StatusInternalServerError,
			},
		},
		{ //ListCategory. Comment repo error
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/posts/{CATEGORY}", nil)
				req := mux.SetURLVars(r, map[string]string{
					"CATEGORY": "music",
				})
				return req
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().GetCategory("music"),
				mockCommentsRepo.EXPECT().GetByID(gomock.Any()),
			},
			ReturnMockFunc: [][]interface{}{
				{testPosts, nil},
				{nil, fmt.Errorf("Internal error")},
			},
			HandlerFunc: postsTestHandler.ListCategory,
			ExpectResult: Result{
				Body: []byte("DB error\n"),
				Code: http.StatusInternalServerError,
			},
		},
		{ //ListByUserLogin SUCCESS
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/user/{USER_LOGIN}", nil)
				req := mux.SetURLVars(r, map[string]string{
					"USER_LOGIN": "rvasily",
				})
				return req
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().GetByUserLogin("rvasily"),
				mockCommentsRepo.EXPECT().GetByID(gomock.Any()),
			},
			ReturnMockFunc: [][]interface{}{
				{testPosts, nil},
				{testComment, nil},
			},
			HandlerFunc: postsTestHandler.ListByUserLogin,
			ExpectResult: Result{
				Body: []byte(`[{"author":{"username":"rvasily","id":"1"},"category":"music","comments":[{"id":"5ebaf9f23c04c17c56f51245","author":{"username":"rvasily","id":"1"},"body":"something","created":"2020-05-12T22:33:02+03:00"}],"created":"2020-05-12T22:33:02+03:00","id":"5ebaf9ee3c04c17c56f51244","score":1,"text":"Something 1","title":"Lorem","type":"text","upvotePercentage":100,"views":10,"votes":[{"user":"1","vote":1}]}]`),
				Code: http.StatusOK,
			},
		},
		{ //ListByUserLogin. Post repo error
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/user/{USER_LOGIN}", nil)
				req := mux.SetURLVars(r, map[string]string{
					"USER_LOGIN": "rvasily",
				})
				return req
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().GetByUserLogin("rvasily"),
			},
			ReturnMockFunc: [][]interface{}{
				{nil, fmt.Errorf("Internal error")},
			},
			HandlerFunc: postsTestHandler.ListByUserLogin,
			ExpectResult: Result{
				Body: []byte("DB error\n"),
				Code: http.StatusInternalServerError,
			},
		},
		{ //ListByUserLogin. Post repo error
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/user/{USER_LOGIN}", nil)
				req := mux.SetURLVars(r, map[string]string{
					"USER_LOGIN": "rvasily",
				})
				return req
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().GetByUserLogin("rvasily"),
				mockCommentsRepo.EXPECT().GetByID(gomock.Any()),
			},
			ReturnMockFunc: [][]interface{}{
				{testPosts, nil},
				{nil, fmt.Errorf("Internal error")},
			},
			HandlerFunc: postsTestHandler.ListByUserLogin,
			ExpectResult: Result{
				Body: []byte("DB error\n"),
				Code: http.StatusInternalServerError,
			},
		},
		{ //ListByID SUCCESS
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/post/{ID}", nil)
				req := mux.SetURLVars(r, map[string]string{
					"ID": testPost1.ID.Hex(),
				})
				return req
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().GetByID(testPost1.ID),
				mockPostsRepo.EXPECT().UpViews(testPost1.ID),
				mockCommentsRepo.EXPECT().GetByID(gomock.Any()),
			},
			ReturnMockFunc: [][]interface{}{
				{testPost1, nil},
				{nil},
				{testComment, nil},
			},
			HandlerFunc: postsTestHandler.ListByID,
			ExpectResult: Result{
				Body: []byte(`{"author":{"username":"rvasily","id":"1"},"category":"music","comments":[{"id":"5ebaf9f23c04c17c56f51245","author":{"username":"rvasily","id":"1"},"body":"something","created":"2020-05-12T22:33:02+03:00"}],"created":"2020-05-12T22:33:02+03:00","id":"5ebaf9ee3c04c17c56f51244","score":1,"text":"Something 1","title":"Lorem","type":"text","upvotePercentage":100,"views":10,"votes":[{"user":"1","vote":1}]}`),
				Code: http.StatusOK,
			},
		},
		{ //ListByID Error bad request ID
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/post/{ID}", nil)
				req := mux.SetURLVars(r, map[string]string{
					"ID": "Something",
				})
				return req
			}(),
			ExpectMockFunc: []*gomock.Call{},
			ReturnMockFunc: [][]interface{}{},
			HandlerFunc:    postsTestHandler.ListByID,
			ExpectResult: Result{
				Body: []byte("Bad ID\n"),
				Code: http.StatusInternalServerError,
			},
		},
		{ //ListByID. Error post repo
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/post/{ID}", nil)
				req := mux.SetURLVars(r, map[string]string{
					"ID": testPost1.ID.Hex(),
				})
				return req
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().GetByID(testPost1.ID),
			},
			ReturnMockFunc: [][]interface{}{
				{nil, fmt.Errorf("Internal error")},
			},
			HandlerFunc: postsTestHandler.ListByID,
			ExpectResult: Result{
				Body: []byte("DB err\n"),
				Code: http.StatusInternalServerError,
			},
		},
		{ //ListByID. Error UpViews
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/post/{ID}", nil)
				req := mux.SetURLVars(r, map[string]string{
					"ID": testPost1.ID.Hex(),
				})
				return req
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().GetByID(testPost1.ID),
				mockPostsRepo.EXPECT().UpViews(testPost1.ID),
			},
			ReturnMockFunc: [][]interface{}{
				{testPost1, nil},
				{fmt.Errorf("Internal error")},
			},
			HandlerFunc: postsTestHandler.ListByID,
			ExpectResult: Result{
				Body: []byte("DB err\n"),
				Code: http.StatusInternalServerError,
			},
		},
		{ //ListByID. Error comment repo
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/post/{ID}", nil)
				req := mux.SetURLVars(r, map[string]string{
					"ID": testPost1.ID.Hex(),
				})
				return req
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().GetByID(testPost1.ID),
				mockPostsRepo.EXPECT().UpViews(testPost1.ID),
				mockCommentsRepo.EXPECT().GetByID(gomock.Any()),
			},
			ReturnMockFunc: [][]interface{}{
				{testPost1, nil},
				{nil},
				{nil, fmt.Errorf("Internal error")},
			},
			HandlerFunc: postsTestHandler.ListByID,
			ExpectResult: Result{
				Body: []byte("DB err\n"),
				Code: http.StatusInternalServerError,
			},
		},
		{ //Add SUCCESS
			Request: func() *http.Request {
				reqBody := NewPostRequest{
					Category: "music",
					Text:     "Something",
					Link:     "",
					Title:    "Something",
					Type:     "text",
				}
				bodyJSON, _ := json.Marshal(reqBody)
				body := bytes.NewReader(bodyJSON)
				r := httptest.NewRequest("POST", "/api/posts", body)
				token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODk3MjMwODksImlhdCI6MTU4OTcyMjQ4OSwic2Vzc2lvbklkIjo0LCJ1c2VyIjp7ImlkIjoiMSIsInVzZXJuYW1lIjoicnZhc2lseSJ9fQ.mQT_Ftp7s1mBsG2LWx0LIYr3F300KT16Dr5Qk0-OMAQ"
				r.Header.Set("Authorization", token)
				sess := &session.Session{
					ID:   1,
					User: testUser,
				}
				ctx := context.WithValue(r.Context(), session.SessionKey, sess)
				req := r.WithContext(ctx)
				return req
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().Add(
					testUser,
					"music",
					"Something",
					"text",
					"Something",
					"",
				),
				mockCommentsRepo.EXPECT().GetByID(gomock.Any()),
			},
			ReturnMockFunc: [][]interface{}{
				{testPost1, nil},
				{testComment, nil},
			},
			HandlerFunc: postsTestHandler.Add,
			ExpectResult: Result{
				Body: []byte(`{"author":{"username":"rvasily","id":"1"},"category":"music","comments":[{"id":"5ebaf9f23c04c17c56f51245","author":{"username":"rvasily","id":"1"},"body":"something","created":"2020-05-12T22:33:02+03:00"}],"created":"2020-05-12T22:33:02+03:00","id":"5ebaf9ee3c04c17c56f51244","score":1,"text":"Something 1","title":"Lorem","type":"text","upvotePercentage":100,"views":10,"votes":[{"user":"1","vote":1}]}`),
				Code: http.StatusOK,
			},
		},
		{ //Add. Session error
			Request: func() *http.Request {
				reqBody := NewPostRequest{
					Category: "music",
					Text:     "Something",
					Link:     "",
					Title:    "Something",
					Type:     "text",
				}
				bodyJSON, _ := json.Marshal(reqBody)
				body := bytes.NewReader(bodyJSON)
				r := httptest.NewRequest("POST", "/api/posts", body)
				return r
			}(),
			ExpectMockFunc: []*gomock.Call{},
			ReturnMockFunc: [][]interface{}{},
			HandlerFunc:    postsTestHandler.Add,
			ExpectResult: Result{
				Body: []byte("Bad auth\n"),
				Code: http.StatusBadRequest,
			},
		},
		{ //Add. Body error
			Request: func() *http.Request {
				r := httptest.NewRequest("POST", "/api/posts", nil)
				token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODk3MjMwODksImlhdCI6MTU4OTcyMjQ4OSwic2Vzc2lvbklkIjo0LCJ1c2VyIjp7ImlkIjoiMSIsInVzZXJuYW1lIjoicnZhc2lseSJ9fQ.mQT_Ftp7s1mBsG2LWx0LIYr3F300KT16Dr5Qk0-OMAQ"
				r.Header.Set("Authorization", token)
				sess := &session.Session{
					ID:   1,
					User: testUser,
				}
				ctx := context.WithValue(r.Context(), session.SessionKey, sess)
				req := r.WithContext(ctx)
				return req
			}(),
			ExpectMockFunc: []*gomock.Call{},
			ReturnMockFunc: [][]interface{}{},
			HandlerFunc:    postsTestHandler.Add,
			ExpectResult: Result{
				Body: []byte("Internal error\n"),
				Code: http.StatusBadRequest,
			},
		},
		{ //Add. Post repo error
			Request: func() *http.Request {
				reqBody := NewPostRequest{
					Category: "music",
					Text:     "Something",
					Link:     "",
					Title:    "Something",
					Type:     "text",
				}
				bodyJSON, _ := json.Marshal(reqBody)
				body := bytes.NewReader(bodyJSON)
				r := httptest.NewRequest("POST", "/api/posts", body)
				token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODk3MjMwODksImlhdCI6MTU4OTcyMjQ4OSwic2Vzc2lvbklkIjo0LCJ1c2VyIjp7ImlkIjoiMSIsInVzZXJuYW1lIjoicnZhc2lseSJ9fQ.mQT_Ftp7s1mBsG2LWx0LIYr3F300KT16Dr5Qk0-OMAQ"
				r.Header.Set("Authorization", token)
				sess := &session.Session{
					ID:   1,
					User: testUser,
				}
				ctx := context.WithValue(r.Context(), session.SessionKey, sess)
				req := r.WithContext(ctx)
				return req
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().Add(
					testUser,
					"music",
					"Something",
					"text",
					"Something",
					"",
				),
				mockCommentsRepo.EXPECT().GetByID(gomock.Any()),
			},
			ReturnMockFunc: [][]interface{}{
				{testPost1, nil},
				//{testComment, nil},
				{nil, fmt.Errorf("Internal error")},
			},
			HandlerFunc: postsTestHandler.Add,
			ExpectResult: Result{
				Body: []byte("BD error\n"),
				Code: http.StatusInternalServerError,
			},
		},
		{ //Add. Comment repo error
			Request: func() *http.Request {
				reqBody := NewPostRequest{
					Category: "music",
					Text:     "Something",
					Link:     "",
					Title:    "Something",
					Type:     "text",
				}
				bodyJSON, _ := json.Marshal(reqBody)
				body := bytes.NewReader(bodyJSON)
				r := httptest.NewRequest("POST", "/api/posts", body)
				token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODk3MjMwODksImlhdCI6MTU4OTcyMjQ4OSwic2Vzc2lvbklkIjo0LCJ1c2VyIjp7ImlkIjoiMSIsInVzZXJuYW1lIjoicnZhc2lseSJ9fQ.mQT_Ftp7s1mBsG2LWx0LIYr3F300KT16Dr5Qk0-OMAQ"
				r.Header.Set("Authorization", token)
				sess := &session.Session{
					ID:   1,
					User: testUser,
				}
				ctx := context.WithValue(r.Context(), session.SessionKey, sess)
				req := r.WithContext(ctx)
				return req
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().Add(
					testUser,
					"music",
					"Something",
					"text",
					"Something",
					"",
				),
			},
			ReturnMockFunc: [][]interface{}{
				{nil, fmt.Errorf("Internal error")},
			},
			HandlerFunc: postsTestHandler.Add,
			ExpectResult: Result{
				Body: []byte("BD error\n"),
				Code: http.StatusInternalServerError,
			},
		},
		{ //Upvote SUCCESS
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/post/{POST_ID}/upvote", nil)
				token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODk3MjMwODksImlhdCI6MTU4OTcyMjQ4OSwic2Vzc2lvbklkIjo0LCJ1c2VyIjp7ImlkIjoiMSIsInVzZXJuYW1lIjoicnZhc2lseSJ9fQ.mQT_Ftp7s1mBsG2LWx0LIYr3F300KT16Dr5Qk0-OMAQ"
				r.Header.Set("Authorization", token)
				sess := &session.Session{
					ID:   1,
					User: testUser,
				}
				ctx := context.WithValue(r.Context(), session.SessionKey, sess)
				req := r.WithContext(ctx)
				reqID := mux.SetURLVars(req, map[string]string{
					"POST_ID": testPost1.ID.Hex(),
				})
				return reqID
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().Upvote(gomock.Any(), testPost1.ID),
				mockCommentsRepo.EXPECT().GetByID(gomock.Any()),
			},
			ReturnMockFunc: [][]interface{}{
				{testPost1, nil},
				{testComment, nil},
			},
			HandlerFunc: postsTestHandler.Upvote,
			ExpectResult: Result{
				Body: []byte(`{"author":{"username":"rvasily","id":"1"},"category":"music","comments":[{"id":"5ebaf9f23c04c17c56f51245","author":{"username":"rvasily","id":"1"},"body":"something","created":"2020-05-12T22:33:02+03:00"}],"created":"2020-05-12T22:33:02+03:00","id":"5ebaf9ee3c04c17c56f51244","score":1,"text":"Something 1","title":"Lorem","type":"text","upvotePercentage":100,"views":10,"votes":[{"user":"1","vote":1}]}`),
				Code: http.StatusOK,
			},
		},
		{ //Upvote session error
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/post/{POST_ID}/upvote", nil)
				reqID := mux.SetURLVars(r, map[string]string{
					"POST_ID": testPost1.ID.Hex(),
				})
				return reqID
			}(),
			ExpectMockFunc: []*gomock.Call{},
			ReturnMockFunc: [][]interface{}{},
			HandlerFunc:    postsTestHandler.Upvote,
			ExpectResult: Result{
				Body: []byte("Bad auth\n"),
				Code: http.StatusBadRequest,
			},
		},
		{ //Upvote bad post id
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/post/{POST_ID}/upvote", nil)
				token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODk3MjMwODksImlhdCI6MTU4OTcyMjQ4OSwic2Vzc2lvbklkIjo0LCJ1c2VyIjp7ImlkIjoiMSIsInVzZXJuYW1lIjoicnZhc2lseSJ9fQ.mQT_Ftp7s1mBsG2LWx0LIYr3F300KT16Dr5Qk0-OMAQ"
				r.Header.Set("Authorization", token)
				sess := &session.Session{
					ID:   1,
					User: testUser,
				}
				ctx := context.WithValue(r.Context(), session.SessionKey, sess)
				req := r.WithContext(ctx)
				return req
			}(),
			ExpectMockFunc: []*gomock.Call{},
			ReturnMockFunc: [][]interface{}{},
			HandlerFunc:    postsTestHandler.Upvote,
			ExpectResult: Result{
				Body: []byte("Bad id\n"),
				Code: http.StatusBadRequest,
			},
		},
		{ //Upvote post repo error
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/post/{POST_ID}/upvote", nil)
				token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODk3MjMwODksImlhdCI6MTU4OTcyMjQ4OSwic2Vzc2lvbklkIjo0LCJ1c2VyIjp7ImlkIjoiMSIsInVzZXJuYW1lIjoicnZhc2lseSJ9fQ.mQT_Ftp7s1mBsG2LWx0LIYr3F300KT16Dr5Qk0-OMAQ"
				r.Header.Set("Authorization", token)
				sess := &session.Session{
					ID:   1,
					User: testUser,
				}
				ctx := context.WithValue(r.Context(), session.SessionKey, sess)
				req := r.WithContext(ctx)
				reqID := mux.SetURLVars(req, map[string]string{
					"POST_ID": testPost1.ID.Hex(),
				})
				return reqID
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().Upvote(gomock.Any(), testPost1.ID),
			},
			ReturnMockFunc: [][]interface{}{
				{nil, fmt.Errorf("Internal error")},
			},
			HandlerFunc: postsTestHandler.Upvote,
			ExpectResult: Result{
				Body: []byte("Bad upvote\n"),
				Code: http.StatusInternalServerError,
			},
		},
		{ //Upvote comment repo error
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/post/{POST_ID}/upvote", nil)
				token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODk3MjMwODksImlhdCI6MTU4OTcyMjQ4OSwic2Vzc2lvbklkIjo0LCJ1c2VyIjp7ImlkIjoiMSIsInVzZXJuYW1lIjoicnZhc2lseSJ9fQ.mQT_Ftp7s1mBsG2LWx0LIYr3F300KT16Dr5Qk0-OMAQ"
				r.Header.Set("Authorization", token)
				sess := &session.Session{
					ID:   1,
					User: testUser,
				}
				ctx := context.WithValue(r.Context(), session.SessionKey, sess)
				req := r.WithContext(ctx)
				reqID := mux.SetURLVars(req, map[string]string{
					"POST_ID": testPost1.ID.Hex(),
				})
				return reqID
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().Upvote(gomock.Any(), testPost1.ID),
				mockCommentsRepo.EXPECT().GetByID(gomock.Any()),
			},
			ReturnMockFunc: [][]interface{}{
				{testPost1, nil},
				{nil, fmt.Errorf("Internal error")},
			},
			HandlerFunc: postsTestHandler.Upvote,
			ExpectResult: Result{
				Body: []byte("Internal error\n"),
				Code: http.StatusInternalServerError,
			},
		},
		{ //Downvote SUCCESS
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/post/{POST_ID}/upvote", nil)
				token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODk3MjMwODksImlhdCI6MTU4OTcyMjQ4OSwic2Vzc2lvbklkIjo0LCJ1c2VyIjp7ImlkIjoiMSIsInVzZXJuYW1lIjoicnZhc2lseSJ9fQ.mQT_Ftp7s1mBsG2LWx0LIYr3F300KT16Dr5Qk0-OMAQ"
				r.Header.Set("Authorization", token)
				sess := &session.Session{
					ID:   1,
					User: testUser,
				}
				ctx := context.WithValue(r.Context(), session.SessionKey, sess)
				req := r.WithContext(ctx)
				reqID := mux.SetURLVars(req, map[string]string{
					"POST_ID": testPost1.ID.Hex(),
				})
				return reqID
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().Downvote(gomock.Any(), testPost1.ID),
				mockCommentsRepo.EXPECT().GetByID(gomock.Any()),
			},
			ReturnMockFunc: [][]interface{}{
				{testPost1, nil},
				{testComment, nil},
			},
			HandlerFunc: postsTestHandler.Downvote,
			ExpectResult: Result{
				Body: []byte(`{"author":{"username":"rvasily","id":"1"},"category":"music","comments":[{"id":"5ebaf9f23c04c17c56f51245","author":{"username":"rvasily","id":"1"},"body":"something","created":"2020-05-12T22:33:02+03:00"}],"created":"2020-05-12T22:33:02+03:00","id":"5ebaf9ee3c04c17c56f51244","score":1,"text":"Something 1","title":"Lorem","type":"text","upvotePercentage":100,"views":10,"votes":[{"user":"1","vote":1}]}`),
				Code: http.StatusOK,
			},
		},
		{ //Upvote session error
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/post/{POST_ID}/upvote", nil)
				reqID := mux.SetURLVars(r, map[string]string{
					"POST_ID": testPost1.ID.Hex(),
				})
				return reqID
			}(),
			ExpectMockFunc: []*gomock.Call{},
			ReturnMockFunc: [][]interface{}{},
			HandlerFunc:    postsTestHandler.Downvote,
			ExpectResult: Result{
				Body: []byte("Bad auth\n"),
				Code: http.StatusBadRequest,
			},
		},
		{ //Upvote bad post id
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/post/{POST_ID}/upvote", nil)
				token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODk3MjMwODksImlhdCI6MTU4OTcyMjQ4OSwic2Vzc2lvbklkIjo0LCJ1c2VyIjp7ImlkIjoiMSIsInVzZXJuYW1lIjoicnZhc2lseSJ9fQ.mQT_Ftp7s1mBsG2LWx0LIYr3F300KT16Dr5Qk0-OMAQ"
				r.Header.Set("Authorization", token)
				sess := &session.Session{
					ID:   1,
					User: testUser,
				}
				ctx := context.WithValue(r.Context(), session.SessionKey, sess)
				req := r.WithContext(ctx)
				return req
			}(),
			ExpectMockFunc: []*gomock.Call{},
			ReturnMockFunc: [][]interface{}{},
			HandlerFunc:    postsTestHandler.Downvote,
			ExpectResult: Result{
				Body: []byte("Bad id\n"),
				Code: http.StatusBadRequest,
			},
		},
		{ //Upvote post repo error
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/post/{POST_ID}/upvote", nil)
				token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODk3MjMwODksImlhdCI6MTU4OTcyMjQ4OSwic2Vzc2lvbklkIjo0LCJ1c2VyIjp7ImlkIjoiMSIsInVzZXJuYW1lIjoicnZhc2lseSJ9fQ.mQT_Ftp7s1mBsG2LWx0LIYr3F300KT16Dr5Qk0-OMAQ"
				r.Header.Set("Authorization", token)
				sess := &session.Session{
					ID:   1,
					User: testUser,
				}
				ctx := context.WithValue(r.Context(), session.SessionKey, sess)
				req := r.WithContext(ctx)
				reqID := mux.SetURLVars(req, map[string]string{
					"POST_ID": testPost1.ID.Hex(),
				})
				return reqID
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().Downvote(gomock.Any(), testPost1.ID),
			},
			ReturnMockFunc: [][]interface{}{
				{nil, fmt.Errorf("Internal error")},
			},
			HandlerFunc: postsTestHandler.Downvote,
			ExpectResult: Result{
				Body: []byte("Bad downvote\n"),
				Code: http.StatusInternalServerError,
			},
		},
		{ //Upvote comment repo error
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/post/{POST_ID}/upvote", nil)
				token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODk3MjMwODksImlhdCI6MTU4OTcyMjQ4OSwic2Vzc2lvbklkIjo0LCJ1c2VyIjp7ImlkIjoiMSIsInVzZXJuYW1lIjoicnZhc2lseSJ9fQ.mQT_Ftp7s1mBsG2LWx0LIYr3F300KT16Dr5Qk0-OMAQ"
				r.Header.Set("Authorization", token)
				sess := &session.Session{
					ID:   1,
					User: testUser,
				}
				ctx := context.WithValue(r.Context(), session.SessionKey, sess)
				req := r.WithContext(ctx)
				reqID := mux.SetURLVars(req, map[string]string{
					"POST_ID": testPost1.ID.Hex(),
				})
				return reqID
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().Downvote(gomock.Any(), testPost1.ID),
				mockCommentsRepo.EXPECT().GetByID(gomock.Any()),
			},
			ReturnMockFunc: [][]interface{}{
				{testPost1, nil},
				{nil, fmt.Errorf("Internal error")},
			},
			HandlerFunc: postsTestHandler.Downvote,
			ExpectResult: Result{
				Body: []byte("Internal error\n"),
				Code: http.StatusInternalServerError,
			},
		},
		{ //Delete SUCCESS
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/post/{POST_ID}/upvote", nil)
				token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODk3MjMwODksImlhdCI6MTU4OTcyMjQ4OSwic2Vzc2lvbklkIjo0LCJ1c2VyIjp7ImlkIjoiMSIsInVzZXJuYW1lIjoicnZhc2lseSJ9fQ.mQT_Ftp7s1mBsG2LWx0LIYr3F300KT16Dr5Qk0-OMAQ"
				r.Header.Set("Authorization", token)
				sess := &session.Session{
					ID:   1,
					User: testUser,
				}
				ctx := context.WithValue(r.Context(), session.SessionKey, sess)
				req := r.WithContext(ctx)
				reqID := mux.SetURLVars(req, map[string]string{
					"POST_ID": testPost1.ID.Hex(),
				})
				return reqID
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().Delete(testPost1.ID),
			},
			ReturnMockFunc: [][]interface{}{
				{true, nil},
			},
			HandlerFunc: postsTestHandler.Delete,
			ExpectResult: Result{
				Body: []byte("{\"message\": \"success\"}"),
				Code: http.StatusOK,
			},
		},
		{ //Delete bad auth
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/post/{POST_ID}/upvote", nil)
				reqID := mux.SetURLVars(r, map[string]string{
					"POST_ID": testPost1.ID.Hex(),
				})
				return reqID
			}(),
			ExpectMockFunc: []*gomock.Call{},
			ReturnMockFunc: [][]interface{}{},
			HandlerFunc:    postsTestHandler.Delete,
			ExpectResult: Result{
				Body: []byte("Bad auth\n"),
				Code: http.StatusBadRequest,
			},
		},
		{ //Delete bad request
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/post/{POST_ID}/upvote", nil)
				token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODk3MjMwODksImlhdCI6MTU4OTcyMjQ4OSwic2Vzc2lvbklkIjo0LCJ1c2VyIjp7ImlkIjoiMSIsInVzZXJuYW1lIjoicnZhc2lseSJ9fQ.mQT_Ftp7s1mBsG2LWx0LIYr3F300KT16Dr5Qk0-OMAQ"
				r.Header.Set("Authorization", token)
				sess := &session.Session{
					ID:   1,
					User: testUser,
				}
				ctx := context.WithValue(r.Context(), session.SessionKey, sess)
				req := r.WithContext(ctx)
				return req
			}(),
			ExpectMockFunc: []*gomock.Call{},
			ReturnMockFunc: [][]interface{}{},
			HandlerFunc:    postsTestHandler.Delete,
			ExpectResult: Result{
				Body: []byte("Bad id\n"),
				Code: http.StatusBadRequest,
			},
		},
		{ //Delete post repo error false
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/post/{POST_ID}/upvote", nil)
				token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODk3MjMwODksImlhdCI6MTU4OTcyMjQ4OSwic2Vzc2lvbklkIjo0LCJ1c2VyIjp7ImlkIjoiMSIsInVzZXJuYW1lIjoicnZhc2lseSJ9fQ.mQT_Ftp7s1mBsG2LWx0LIYr3F300KT16Dr5Qk0-OMAQ"
				r.Header.Set("Authorization", token)
				sess := &session.Session{
					ID:   1,
					User: testUser,
				}
				ctx := context.WithValue(r.Context(), session.SessionKey, sess)
				req := r.WithContext(ctx)
				reqID := mux.SetURLVars(req, map[string]string{
					"POST_ID": testPost1.ID.Hex(),
				})
				return reqID
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().Delete(testPost1.ID),
			},
			ReturnMockFunc: [][]interface{}{
				{false, nil},
			},
			HandlerFunc: postsTestHandler.Delete,
			ExpectResult: Result{
				Body: []byte("{\"message\": \"failure\"}"),
				Code: http.StatusOK,
			},
		},
		{ //Delete post repo error
			Request: func() *http.Request {
				r := httptest.NewRequest("GET", "/api/post/{POST_ID}/upvote", nil)
				token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODk3MjMwODksImlhdCI6MTU4OTcyMjQ4OSwic2Vzc2lvbklkIjo0LCJ1c2VyIjp7ImlkIjoiMSIsInVzZXJuYW1lIjoicnZhc2lseSJ9fQ.mQT_Ftp7s1mBsG2LWx0LIYr3F300KT16Dr5Qk0-OMAQ"
				r.Header.Set("Authorization", token)
				sess := &session.Session{
					ID:   1,
					User: testUser,
				}
				ctx := context.WithValue(r.Context(), session.SessionKey, sess)
				req := r.WithContext(ctx)
				reqID := mux.SetURLVars(req, map[string]string{
					"POST_ID": testPost1.ID.Hex(),
				})
				return reqID
			}(),
			ExpectMockFunc: []*gomock.Call{
				mockPostsRepo.EXPECT().Delete(testPost1.ID),
			},
			ReturnMockFunc: [][]interface{}{
				{false, fmt.Errorf("Internal error")},
			},
			HandlerFunc: postsTestHandler.Delete,
			ExpectResult: Result{
				Body: []byte("Delete error\n"),
				Code: http.StatusInternalServerError,
			},
		},
	}

	for iTestCase, testCase := range testCases {
		w := httptest.NewRecorder()
		for i, mockFunc := range testCase.ExpectMockFunc {
			mockFunc.Return(testCase.ReturnMockFunc[i]...)
		}
		testCase.HandlerFunc(w, testCase.Request)

		resp := w.Result()

		body, _ := ioutil.ReadAll(resp.Body)
		code := resp.StatusCode

		if assert.Equal(t, body, testCase.ExpectResult.Body) &&
			assert.Equal(t, code, testCase.ExpectResult.Code) {
			fmt.Printf("CASE %d SUCCESS\n", iTestCase)
		}
	}
}

func TestInitError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostsRepo := NewMockPostsRepositoryInterface(ctrl)
	mockCommentsRepo := NewMockCommentsRepositoryInterface(ctrl)
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	postsTestHandler := &PostsHandler{
		Tmpl:        template.New("Something"),
		Logger:      logger,
		PostsRepo:   mockPostsRepo,
		CommentRepo: mockCommentsRepo,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	postsTestHandler.Init(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	code := resp.StatusCode
	if assert.Equal(t, body, []byte("Template error\n")) &&
		assert.Equal(t, code, http.StatusInternalServerError) {
		fmt.Printf("CASE Error Init SUCCESS\n")
	}
}
