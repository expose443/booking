package repository

import (
	"database/sql"

	"github.com/with-insomnia/Hotel/internal/model"
)

type UserQuery interface {
	CreateUser(user *model.User) error
	GetUserIdByToken(token string) (int, error)
	GetUserByUserId(userID int) (model.User, error)
	GetUserByEmail(email string) (model.User, error)
}

type userQuery struct {
	db *sql.DB
}

func (u *userQuery) CreateUser(user *model.User) error {
	_, err := u.db.Exec("INSERT INTO users(email, first_name, last_name, password, number, role) values($1, $2, $3, $4, $5, $6)", user.Email, user.FirstName, user.LastName, user.Password, user.Number, "user")
	if err != nil {
		return err
	}
	return nil
}

func (u *userQuery) GetUserIdByToken(token string) (int, error) {
	row := u.db.QueryRow("SELECT user_id FROM sessions WHERE token=$1", token)
	var userID int
	if err := row.Scan(&userID); err != nil {
		return -1, err
	}
	return userID, nil
}

func (u *userQuery) GetUserByUserId(userID int) (model.User, error) {
	row := u.db.QueryRow("SELECT user_id, email, password, first_name, last_name, number, role FROM users WHERE user_id = $1", userID)
	var user model.User
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Number, &user.Role); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u *userQuery) GetUserByEmail(email string) (model.User, error) {
	row := u.db.QueryRow("SELECT user_id,email,password, first_name, last_name, number, role FROM users WHERE email = $1", email)
	var user model.User
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Number, &user.Role); err != nil {
		return model.User{}, err
	}
	return user, nil
}
