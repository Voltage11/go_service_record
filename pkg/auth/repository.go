package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type repository struct {
	db *sqlx.DB
	logger *zerolog.Logger
}

func newRepository(db *sqlx.DB, logger *zerolog.Logger) *repository {
	return &repository{
		db: db,
		logger: logger,
	}
}

func (r *repository) getUserByEmail(ctx context.Context, email string) (*User, error) {
	if email == "" {
		return nil, errors.New("пустой email")
	}
	stmt := `SELECT id, username, password_hash, email, is_active, is_admin FROM users WHERE email = $1`
	user := &User{}
	err := r.db.GetContext(ctx, user, stmt, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return user, nil
}

func (r *repository) getSessionById(ctx context.Context, id string) (*session, error) {
	if id == "" {
		return nil, errors.New("пустой id сессии")
	}
	stmt := `SELECT id, user_id, expires_at, ip, user_agent, created_at FROM sessions WHERE id = $1`
	session := &session{}
	err := r.db.GetContext(ctx, session, stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			r.logger.Error().Err(err).Msg(fmt.Sprintf("ошибка при получении сессии: %s", id))
			return nil, err
		}
	}
	return session, nil
}

//func (r *repository) updateSe

func (r *repository) addSession(ctx context.Context, session *session) error {
	if session == nil {
		return errors.New("пустая сессия")
	}
	// Проверяем, что сессия не существует
	existsSession, err := r.getSessionById(ctx, session.Id)
	if err != nil {
		return err
	}
	if existsSession != nil {
		// TODO: Обновляем сессию
		return nil
	}
	
	stmt := `INSET INTO sessions (id, user_id, expires_at, ip, user_agent, created_at) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = r.db.ExecContext(ctx, stmt, session.Id, session.UserId, session.ExpiresAt, session.Ip, session.UserAgent, session.CreatedAt)
	if err != nil {
		r.logger.Error().Err(err).Msg(fmt.Sprintf("ошибка при добавлении сессии: %v", session))
		return err
	}
	
	r.logger.Info().Msg(fmt.Sprintf("сессия добавлена: %v", session))
	return nil
}