package user

import (
	"fmt"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestAdd(t *testing.T) {
	assert := assert.New(t)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var repo *UserRepo
	repo = NewUserRepo(db)

	testUser := &User{
		ID:       int64(1),
		Username: "rvasily",
		password: "lovelove",
	}

	mock.
		ExpectExec("INSERT INTO users").
		WithArgs(testUser.Username, testUser.password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	userID, err := repo.Add(testUser.Username, testUser.password)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if !assert.Equal(userID, testUser.ID) {
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error
	mock.
		ExpectExec("INSERT INTO users").
		WithArgs(testUser.Username, testUser.password).
		WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Something wrong"})
	userID, err = repo.Add(testUser.Username, testUser.password)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if !assert.Equal(err, ErrAlreadyExisting) {
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	// result error
	mock.
		ExpectExec(`INSERT INTO users`).
		WithArgs(testUser.Username, testUser.password).
		WillReturnError(fmt.Errorf("bad_result"))

	_, err = repo.Add(testUser.Username, testUser.password)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetByID(t *testing.T) {
	assert := assert.New(t)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var repo *UserRepo
	repo = NewUserRepo(db)

	testUser := &User{
		ID:       int64(1),
		Username: "rvasily",
		password: "lovelove",
	}

	rows := sqlmock.
		NewRows([]string{"username", "password"}).
		AddRow(testUser.Username, testUser.password)

	mock.
		ExpectQuery("SELECT username, password FROM users WHERE").
		WithArgs(testUser.ID).
		WillReturnRows(rows)

	user, err := repo.GetByID(testUser.ID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if !assert.Equal(testUser, user) {
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error
	mock.
		ExpectQuery("SELECT username, password FROM users WHERE").
		WithArgs(testUser.ID).
		WillReturnError(fmt.Errorf("bad_result"))

	user, err = repo.GetByID(testUser.ID)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAuthorize(t *testing.T) {
	assert := assert.New(t)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var repo *UserRepo
	repo = NewUserRepo(db)

	testUser := &User{
		ID:       int64(1),
		Username: "rvasily",
		password: "lovelove",
	}

	rows := sqlmock.
		NewRows([]string{"id", "password"}).
		AddRow(testUser.ID, testUser.password)

	mock.
		ExpectQuery("SELECT `id`, `password` FROM users WHERE").
		WithArgs(testUser.Username).
		WillReturnRows(rows)

	user, err := repo.Authorize(testUser.Username, testUser.password)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if !assert.Equal(testUser, user) {
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error
	mock.
		ExpectQuery("SELECT `id`, `password` FROM users WHERE").
		WithArgs(testUser.Username).
		WillReturnError(fmt.Errorf("bad_result"))

	user, err = repo.Authorize(testUser.Username, testUser.password)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	//bad pass

	rows = sqlmock.
		NewRows([]string{"id", "password"}).
		AddRow(testUser.ID, testUser.password+"7")

	mock.
		ExpectQuery("SELECT `id`, `password` FROM users WHERE").
		WithArgs(testUser.Username).
		WillReturnRows(rows)

	user, err = repo.Authorize(testUser.Username, testUser.password)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if !assert.Equal(err, ErrBadPass) {
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
