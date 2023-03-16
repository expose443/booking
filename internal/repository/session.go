package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/with-insomnia/Hotel/internal/model"
)

type SessionQuery interface {
	GetSessionByToken(token string) (model.Session, error)
	GetSessionByUserID(userId int) (model.Session, error)
	GetAllSessionsTime() ([]model.Session, error)
	CreateSession(session model.Session) error
	DeleteSession(token string) error
}

type sessionQuery struct {
	db *sql.DB
}

func (s *sessionQuery) CreateSession(session model.Session) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO sessions(user_id, token, expiry) VALUES($1,$2,$3)", session.UserId, session.Token, session.Expiry)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *sessionQuery) GetSessionByToken(token string) (model.Session, error) {
	row := s.db.QueryRow("SELECT user_id, token, expiry FROM sessions WHERE token = $1", token)
	var session model.Session
	if err := row.Scan(&session.UserId, &session.Token, &session.Expiry); err != nil {
		return model.Session{}, err
	}
	return session, nil
}

func (s *sessionQuery) GetSessionByUserID(userId int) (model.Session, error) {
	row := s.db.QueryRow("SELECT user_id, token, expiry FROM sessions WHERE user_id = $1", userId)
	var session model.Session
	if err := row.Scan(&session.UserId, &session.Token, &session.Expiry); err != nil {
		fmt.Println(err)
		return model.Session{}, err
	}
	return session, nil
}

func (s *sessionQuery) DeleteSession(token string) error {
	res, err := s.db.Exec("DELETE FROM sessions WHERE token = $1", token)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("delete session was failed")
	}
	return nil
}

func (s *sessionQuery) GetAllSessionsTime() ([]model.Session, error) {
	rows, err := s.db.Query("SELECT expiry, token FROM sessions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sessions []model.Session
	for rows.Next() {
		var session model.Session
		if err := rows.Scan(&session.Expiry, &session.Token); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}
