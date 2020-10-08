package session

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reddit/pkg/user"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var ExampleTokenSecret = []byte("It's top secret")

var (
	ErrPayload          = errors.New("Payload error")
	ErrBadJSON          = errors.New("JSON error")
	ErrJWTBadSignMethod = errors.New("Bad sign method")
	ErrJWTConwertToken  = errors.New("Error convert token")
)

var validTime int64 = 600 // seconds

type SessionsManager struct {
	DB *sql.DB
}

func NewSessionsMem(db *sql.DB) *SessionsManager {
	return &SessionsManager{DB: db}
}

func (sm *SessionsManager) Check(r *http.Request) (*Session, error) {
	inToken := r.Header.Get("Authorization")
	if inToken != "" {
		inToken = strings.Split(inToken, " ")[1]
	} else {
		return nil, ErrNoAuth
	}

	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, ErrJWTBadSignMethod
		}
		return ExampleTokenSecret, nil
	}
	token, err := jwt.Parse(inToken, hashSecretGetter)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("%v, %v", ErrBadJSON, err)
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrPayload
	}
	sessionID := int64(payload["sessionId"].(float64))
	var userIDFromDB, expTimeDB int64
	err = sm.DB.
		QueryRow("SELECT `user_id`, `exp_time` FROM sessions WHERE id = ?", sessionID).
		Scan(&userIDFromDB, &expTimeDB)
	if err != nil {
		return nil, err
	}
	userIDStr := payload["user"].(map[string]interface{})["id"].(string)
	userName := payload["user"].(map[string]interface{})["username"].(string)
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Bad user. Err: %v", err)
	}
	if userID != userIDFromDB {
		return nil, fmt.Errorf("Bad user")
	}
	if expTimeDB < time.Now().Unix() {
		return nil, fmt.Errorf("Bad time")
	}

	sess := &Session{
		ID: sessionID,
		User: &user.User{
			ID:       userID,
			Username: userName,
		},
	}

	return sess, nil
}

func (sm *SessionsManager) Create(w http.ResponseWriter, userS *user.User) (int64, error) {
	createTime := time.Now().Unix()
	_, err := sm.DB.Exec(
		"INSERT INTO sessions (`user_id`, `create_time`, `exp_time`) VALUES (?, ?, ?)",
		userS.ID,
		createTime,
		createTime+validTime,
	)
	if err != nil {
		return 0, err
	}
	var sessID int64
	err = sm.DB.
		QueryRow("SELECT id FROM sessions WHERE user_id = ? AND create_time = ? AND exp_time = ?",
			userS.ID,
			createTime,
			createTime+validTime).
		Scan(&sessID)
	if err != nil {
		return 0, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]string{
			"username": userS.Username,
			"id":       strconv.FormatInt(userS.ID, 10),
		},
		"sessionId": sessID,
		"iat":       createTime,
		"exp":       createTime + validTime,
	})
	tokenString, err := token.SignedString(ExampleTokenSecret)
	if err != nil {
		log.Println("Error convert token", err)
		return 0, ErrJWTConwertToken
	}
	resp, _ := json.Marshal(map[string]interface{}{
		"token": tokenString,
	})
	w.Write(resp)
	return sessID, nil
}
