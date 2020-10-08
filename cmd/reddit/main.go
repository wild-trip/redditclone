package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"reddit/pkg/handlers"
	"reddit/pkg/middleware"
	"reddit/pkg/posts"
	"reddit/pkg/session"
	"reddit/pkg/user"

	mgo "gopkg.in/mgo.v2"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	r := mux.NewRouter()
	templates := template.Must(template.ParseFiles("./template/index.html"))
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./template/static"))),
	)
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	//Connect to SQL DB
	dsn := "root:love@tcp(localhost:3306)/golang?charset=utf8&interpolateParams=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Errorf("SQL error: %v", err)
		return
	}
	logger.Infof("SQL connect to DB")
	db.SetMaxOpenConns(10)
	err = db.Ping()
	if err != nil {
		logger.Errorf("Can't connect to SQL: %v", err)
		return
	}
	defer db.Close()
	//Connect to Mondo DB
	sessMongoDB, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer sessMongoDB.Close()
	postsCollection := sessMongoDB.DB("coursera").C("posts")
	commentsCollection := sessMongoDB.DB("coursera").C("comments")
	logger.Infof("MongoDB connect to DB")

	//SQL Database
	sm := session.NewSessionsMem(db)
	userRepo := user.NewUserRepo(db)
	//Mongo DB
	postsRepo := posts.NewRepo(postsCollection)
	commentRepo := posts.NewCommentRepo(commentsCollection)

	userHandler := &handlers.UserHandler{
		Tmpl:     templates,
		UserRepo: userRepo,
		Logger:   logger,
		Sessions: sm,
	}

	handlers := &handlers.PostsHandler{
		Tmpl:        templates,
		Logger:      logger,
		PostsRepo:   postsRepo,
		CommentRepo: commentRepo,
	}

	r.HandleFunc("/api/register", userHandler.SignUp).Methods("POST")
	r.HandleFunc("/api/login", userHandler.Login).Methods("POST")

	r.HandleFunc("/", handlers.Init)
	r.HandleFunc("/api/posts/", handlers.ListAll).Methods("GET")
	r.HandleFunc("/api/posts", handlers.Add).Methods("POST")
	r.HandleFunc("/api/posts/{CATEGORY}", handlers.ListCategory).Methods("GET")
	r.HandleFunc("/api/post/{ID}", handlers.ListByID).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}/upvote", handlers.Upvote).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}/downvote", handlers.Downvote).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}", handlers.Delete).Methods("DELETE")
	r.HandleFunc("/api/user/{USER_LOGIN}", handlers.ListByUserLogin).Methods("GET")

	r.HandleFunc("/api/post/{POST_ID}", handlers.AddComment).Methods("POST")
	r.HandleFunc("/api/post/{POST_ID}/{COMMENT_ID}", handlers.DeleteComment).Methods("DELETE")

	mux := middleware.Auth(sm, r, userRepo)
	mux = middleware.AccessLog(logger, mux)
	mux = middleware.Panic(mux)

	addr := ":8080"
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	http.ListenAndServe(addr, mux)
}
