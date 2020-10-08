package handlers

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"

	"reddit/pkg/session"
	"reddit/pkg/user"

	"go.uber.org/zap"
)

type SessionManagerInterface interface {
	Check(*http.Request) (*session.Session, error)
	Create(http.ResponseWriter, *user.User) (int64, error)
}

type UserRepositoryInterface interface {
	Authorize(string, string) (*user.User, error)
	Add(string, string) (int64, error)
	GetByID(int64) (*user.User, error)
}

type UserHandler struct {
	Tmpl     *template.Template
	Logger   *zap.SugaredLogger
	UserRepo UserRepositoryInterface
	Sessions SessionManagerInterface
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	dataRequest := new(LoginRequest)
	body, errReadBody := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, dataRequest)
	if errReadBody != nil || err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		h.Logger.Errorf("Bad JSON. Error: %v, %v", errReadBody, err)
		return
	}
	u, err := h.UserRepo.Authorize(dataRequest.Username, dataRequest.Password)
	if err == user.ErrNoUser {
		http.Error(w, `no user`, http.StatusBadRequest)
		h.Logger.Errorf("Error: %v", err)
		return
	}
	if err == user.ErrBadPass {
		http.Error(w, `bad pass`, http.StatusBadRequest)
		h.Logger.Errorf("Error: %v", err)
		return
	}
	sessID, errSess := h.Sessions.Create(w, u)
	if errSess == nil {
		h.Logger.Infof("created session sessionID: %v", sessID)
	} else {
		http.Error(w, `bad pass`, http.StatusInternalServerError)
		h.Logger.Errorf("Error: %v", errSess)
	}
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	dataRequest := new(LoginRequest)
	body, errReadBody := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, dataRequest)
	if errReadBody != nil || err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		h.Logger.Errorf("Bad JSON. Error: %v, %v", errReadBody, err)
		return
	}
	userID, err := h.UserRepo.Add(dataRequest.Username, dataRequest.Password)
	if err == user.ErrAlreadyExisting {
		ans, _ := json.Marshal(map[string][]map[string]string{
			"errors": {
				map[string]string{
					"location": "body",
					"param":    "username",
					"value":    dataRequest.Username,
					"msg":      "already exists",
				},
			},
		})
		http.Error(w, string(ans), http.StatusUnprocessableEntity)
		return
	}
	newUser := &user.User{
		ID:       userID,
		Username: dataRequest.Username,
	}
	sessID, errSess := h.Sessions.Create(w, newUser)
	if errSess == nil {
		h.Logger.Infof("created session sessionID: %v", sessID)
	} else {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		h.Logger.Infof("Can't created session")
	}
}
