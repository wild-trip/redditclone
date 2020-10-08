package user

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type User struct {
	Username string `json:"username" bson:"username"`
	ID       int64  `json:"id,string" bson:"id"`
	password string
}

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

var (
	ErrNoUser          = errors.New("No user found")
	ErrBadPass         = errors.New("Invald password")
	ErrAlreadyExisting = errors.New("This user name already existing")
)

func (repo *UserRepo) Authorize(login, pass string) (*User, error) {
	var passwordDB string
	var userID int64
	err := repo.DB.
		QueryRow("SELECT `id`, `password` FROM users WHERE username = ?", login).
		Scan(&userID, &passwordDB)
	fmt.Printf("This : %v, %v, %v", err, userID, passwordDB)
	if err != nil {
		return nil, err
	}
	if passwordDB != pass {
		return nil, ErrBadPass
	}
	user := &User{
		ID:       userID,
		Username: login,
		password: pass,
	}
	return user, nil
}

func (repo *UserRepo) Add(login, pass string) (int64, error) {
	result, err := repo.DB.Exec(
		"INSERT INTO users (`username`, `password`) VALUES (?, ?)",
		login,
		pass,
	)
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		if mysqlError.Number == 1062 {
			return 0, ErrAlreadyExisting
		}
	}
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (repo *UserRepo) GetByID(id int64) (*User, error) {
	var password string
	var username string
	err := repo.DB.
		QueryRow("SELECT username, password FROM users WHERE id = ?",
			id).
		Scan(&username, &password)
	if err != nil {
		return nil, fmt.Errorf("BD error: %v", err)
	}
	user := &User{
		ID:       id,
		Username: username,
		password: password,
	}
	return user, nil
}
