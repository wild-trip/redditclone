package middleware

import (
	"context"
	"log"
	"net/http"
	"regexp"

	"reddit/pkg/session"
	"reddit/pkg/user"
)

var authUrls []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile(`^/api/posts$`),
	regexp.MustCompile(`^/api/post/.+/upvote$`),
	regexp.MustCompile(`^/api/post/.+/downvote$`),
	regexp.MustCompile(`^/api/post/.+/downvote$`),
	regexp.MustCompile(`^/api/post/.+$`),
	regexp.MustCompile(`^/api/post/.+/.+$`),
}

func Auth(sm *session.SessionsManager, next http.Handler, userRepo *user.UserRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("auth middleware")
		var isAuth = false
		for _, regexpPath := range authUrls {
			if regexpPath.MatchString(r.URL.String()) {
				isAuth = true
				break
			}
		}
		if !isAuth {
			log.Println("Continue without auth")
			next.ServeHTTP(w, r)
			return
		}

		sess, err := sm.Check(r)
		if err != nil {
			log.Println("no auth. Error:", err)
			//http.Error(w, "", http.StatusInternalServerError)
			next.ServeHTTP(w, r)
			return
		}
		log.Println("authorization done for", sess.User.Username)
		ctx := context.WithValue(r.Context(), session.SessionKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
