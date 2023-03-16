package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/with-insomnia/Hotel/internal/model"
	"github.com/with-insomnia/Hotel/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(user *model.User) (model.Session, error)
	Register(user *model.User) error
	Logout(token string) error
}

type authService struct {
	sessionQuery repository.SessionQuery
	userQuery    repository.UserQuery
}

func NewAuthService(dao repository.DAO) AuthService {
	return &authService{
		dao.NewSessionQuery(),
		dao.NewUserQuery(),
	}
}

func (a *authService) Login(user *model.User) (model.Session, error) {
	userDB, err := a.userQuery.GetUserByEmail(user.Email)
	if err != nil {
		log.Printf("user %s sign in was failed\n", user.Email)
		return model.Session{}, errors.New("wrong password or email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password))
	if err != nil {
		log.Printf("user %s sign in was failed\n", user.Email)
		return model.Session{}, errors.New("wrong password or email")
	}

	sessionDb, err := a.sessionQuery.GetSessionByUserID(int(userDB.ID))
	if err != nil {
		log.Printf("session for user_id %d is not found\n", userDB.ID)
	} else {
		err := a.sessionQuery.DeleteSession(sessionDb.Token)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("session for user_id %d is deleted\n", user.ID)
		}
	}
	sessionToken := uuid.NewString()
	expiry := time.Now().Add(10 * time.Minute)
	fmt.Println(expiry)
	session := model.Session{
		UserId: userDB.ID,
		Token:  sessionToken,
		Expiry: expiry,
	}

	err = a.sessionQuery.CreateSession(session)
	if err != nil {
		return model.Session{}, fmt.Errorf("session for user %d was failed\nerror: %w", user.ID, err)
	}
	log.Printf("user %s sign in was successfully\n", user.Email)
	return session, nil
}

func (a *authService) Register(user *model.User) error {
	_, emailExist := a.userQuery.GetUserByEmail(user.Email)

	if emailExist == nil {
		return errors.New("user exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 0)
	if err != nil {
		return errors.New("password hash was fail")
	}

	newUser := model.User{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Number:    user.Number,
		Password:  string(passwordHash),
	}

	err = a.userQuery.CreateUser(&newUser)
	if err != nil {
		return err
	}
	return nil
}

func (a *authService) Logout(token string) error {
	return a.sessionQuery.DeleteSession(token)
}
