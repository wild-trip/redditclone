package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http/httptest"
	"reddit/pkg/user"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAuthorize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepositoryInterface(ctrl)
	mockSessionManager := NewMockSessionManagerInterface(ctrl)
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	//Success
	userTestHandler := &UserHandler{
		Tmpl:     template.Must(template.ParseFiles("../../template/index.html")),
		UserRepo: mockRepo,
		Logger:   logger,
		Sessions: mockSessionManager,
	}
	testUser := &user.User{
		ID:       1,
		Username: "rvasily",
	}
	password := "lovelove"
	loginRequest := LoginRequest{
		Username: testUser.Username,
		Password: password,
	}

	testRequest, _ := json.Marshal(loginRequest)
	r := httptest.NewRequest("POST", "/api/login", bytes.NewReader(testRequest))
	w := httptest.NewRecorder()

	mockRepo.EXPECT().Authorize(testUser.Username, password).Return(testUser, nil)

	mockSessionManager.EXPECT().Create(w, testUser).Return(int64(1), nil)

	userTestHandler.Login(w, r)
	resp := w.Result()

	body, _ := ioutil.ReadAll(resp.Body)
	code := resp.StatusCode

	assert.Equal(t, body, []byte(""))
	assert.Equal(t, code, 200)

	//Error bad JSON
	testRequest = []byte("JSON")
	r = httptest.NewRequest("POST", "/api/login", bytes.NewReader(testRequest))
	w = httptest.NewRecorder()

	userTestHandler.Login(w, r)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)
	code = resp.StatusCode

	assert.Equal(t, body, []byte("Bad request\n"))
	assert.Equal(t, code, 400)

	//Error no user
	testRequest, _ = json.Marshal(loginRequest)
	r = httptest.NewRequest("POST", "/api/login", bytes.NewReader(testRequest))
	w = httptest.NewRecorder()

	mockRepo.EXPECT().Authorize(testUser.Username, password).Return(nil, user.ErrNoUser)

	userTestHandler.Login(w, r)
	resp = w.Result()

	body, _ = ioutil.ReadAll(resp.Body)
	code = resp.StatusCode

	assert.Equal(t, body, []byte("no user\n"))
	assert.Equal(t, code, 400)

	//Error bad pass
	testRequest, _ = json.Marshal(loginRequest)
	r = httptest.NewRequest("POST", "/api/login", bytes.NewReader(testRequest))
	w = httptest.NewRecorder()

	mockRepo.EXPECT().Authorize(testUser.Username, password).Return(nil, user.ErrBadPass)

	userTestHandler.Login(w, r)
	resp = w.Result()

	body, _ = ioutil.ReadAll(resp.Body)
	code = resp.StatusCode

	assert.Equal(t, body, []byte("bad pass\n"))
	assert.Equal(t, code, 400)

	//Error session
	testRequest, _ = json.Marshal(loginRequest)
	r = httptest.NewRequest("POST", "/api/login", bytes.NewReader(testRequest))
	w = httptest.NewRecorder()

	mockRepo.EXPECT().Authorize(testUser.Username, password).Return(testUser, nil)

	mockSessionManager.EXPECT().Create(w, testUser).Return(int64(0), fmt.Errorf("Internal error"))

	userTestHandler.Login(w, r)
	resp = w.Result()

	body, _ = ioutil.ReadAll(resp.Body)
	code = resp.StatusCode

	assert.Equal(t, body, []byte("bad pass\n"))
	assert.Equal(t, code, 500)
}

func TestSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepositoryInterface(ctrl)
	mockSessionManager := NewMockSessionManagerInterface(ctrl)
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	//Success
	userTestHandler := &UserHandler{
		Tmpl:     template.Must(template.ParseFiles("../../template/index.html")),
		UserRepo: mockRepo,
		Logger:   logger,
		Sessions: mockSessionManager,
	}
	testUser := &user.User{
		ID:       1,
		Username: "rvasily",
	}
	password := "lovelove"
	loginRequest := LoginRequest{
		Username: testUser.Username,
		Password: password,
	}

	testRequest, _ := json.Marshal(loginRequest)
	r := httptest.NewRequest("POST", "/api/register", bytes.NewReader(testRequest))
	w := httptest.NewRecorder()

	mockRepo.EXPECT().Add(testUser.Username, password).Return(int64(1), nil)

	mockSessionManager.EXPECT().Create(w, testUser).Return(int64(1), nil)

	userTestHandler.SignUp(w, r)
	resp := w.Result()

	body, _ := ioutil.ReadAll(resp.Body)
	code := resp.StatusCode

	assert.Equal(t, body, []byte(""))
	assert.Equal(t, code, 200)

	//Error bad JSON
	testRequest = []byte("JSON")
	r = httptest.NewRequest("POST", "/api/register", bytes.NewReader(testRequest))
	w = httptest.NewRecorder()

	userTestHandler.SignUp(w, r)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)
	code = resp.StatusCode

	assert.Equal(t, body, []byte("Bad request\n"))
	assert.Equal(t, code, 400)

	//Error already existing
	testRequest, _ = json.Marshal(loginRequest)
	r = httptest.NewRequest("POST", "/api/register", bytes.NewReader(testRequest))
	w = httptest.NewRecorder()

	mockRepo.EXPECT().Add(testUser.Username, password).Return(int64(0), user.ErrAlreadyExisting)

	userTestHandler.SignUp(w, r)
	resp = w.Result()

	body, _ = ioutil.ReadAll(resp.Body)
	code = resp.StatusCode

	exp, _ := json.Marshal(map[string][]map[string]string{
		"errors": {
			map[string]string{
				"location": "body",
				"param":    "username",
				"value":    loginRequest.Username,
				"msg":      "already exists",
			},
		},
	})
	exp = append(exp, '\n')
	assert.Equal(t, body, exp)
	assert.Equal(t, code, 422)

	//Error session
	testRequest, _ = json.Marshal(loginRequest)
	r = httptest.NewRequest("POST", "/api/login", bytes.NewReader(testRequest))
	w = httptest.NewRecorder()

	mockRepo.EXPECT().Add(testUser.Username, password).Return(int64(1), nil)

	mockSessionManager.EXPECT().Create(w, gomock.Any()).Return(int64(0), fmt.Errorf("Internal error"))

	userTestHandler.SignUp(w, r)
	resp = w.Result()

	body, _ = ioutil.ReadAll(resp.Body)
	code = resp.StatusCode

	assert.Equal(t, body, []byte("Internal error\n"))
	assert.Equal(t, code, 500)
}
