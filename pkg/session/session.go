package session

import (
	"context"
	"errors"
	"reddit/pkg/user"
)

type Session struct {
	ID   int64
	User *user.User
}

var (
	ErrNoAuth = errors.New("No session found")
)

type sessKey string

var SessionKey sessKey = "sessionKey"

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)
	if !ok || sess == nil {
		return nil, ErrNoAuth
	}
	return sess, nil
}
