package service

import (
	"github.com/with-insomnia/Hotel/internal/model"
	"github.com/with-insomnia/Hotel/internal/repository"
)

type SessionService interface {
	CreateSession(session *model.Session) error
	GetSessionByToken(token string) (model.Session, error)
	GetSessionByUserID(userId int) (model.Session, error)
	GetAllSessionsTime() ([]model.Session, error)
	DeleteSession(token string) error
}

type sessionService struct {
	repository.SessionQuery
}

func NewSessionService(dao repository.DAO) SessionService {
	return &sessionService{dao.NewSessionQuery()}
}

func (s *sessionService) CreateSession(session *model.Session) error {
	err := s.SessionQuery.CreateSession(*session)
	if err != nil {
		return err
	}
	return nil
}

func (s *sessionService) GetSessionByToken(token string) (model.Session, error) {
	session, err := s.SessionQuery.GetSessionByToken(token)
	if err != nil {
		return model.Session{}, err
	}
	return session, nil
}

func (s *sessionService) GetSessionByUserID(userId int) (model.Session, error) {
	session, err := s.SessionQuery.GetSessionByUserID(userId)
	if err != nil {
		return model.Session{}, err
	}
	return session, nil
}

func (s *sessionService) DeleteSession(token string) error {
	return s.SessionQuery.DeleteSession(token)
}

func (s *sessionService) GetAllSessionsTime() ([]model.Session, error) {
	return s.SessionQuery.GetAllSessionsTime()
}
