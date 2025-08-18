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
	db     *sqlx.DB
	logger *zerolog.Logger
}

const (
	stmtGetUserById = `SELECT id, username, password_hash, email, is_active, is_admin 
                       FROM users WHERE id = $1`

	stmtGetUserByEmail = `SELECT id, username, password_hash, email, is_active, is_admin 
                         FROM users WHERE LOWER(email) = LOWER($1)`

	stmtCreateUser = `INSERT INTO users (id, username, password_hash, email, is_active, is_admin, create_at, updated_at)
                      VALUES (:id, :username, :password_hash, :email, :is_active, :is_admin, :create_at, :updated_at)`

	stmtUpdateUser = `UPDATE users SET username = :username, password_hash = :password_hash, 
                      		email = :email, is_active = :is_active, is_admin = :is_admin, updated_at = :updated_at
                      WHERE id = :id`

	stmtGetUserSessions = `SELECT id, user_id, expires_at, ip, user_agent, created_at
                           FROM sessions WHERE user_id = $1`

	stmtGetSessionByID = `SELECT id, user_id, expires_at, ip, user_agent, created_at
                           FROM sessions WHERE id = $1`

	stmtCreateSession = `INSERT INTO sessions (id, user_id, expires_at, ip, user_agent, created_at)
                         VALUES (:id, :user_id, :expires_at, :ip, :user_agent, :created_at)`
)

func newRepository(db *sqlx.DB, logger *zerolog.Logger) (*repository, error) {		
	if db == nil {
		return nil, errors.New("не указан db")
	}
	if logger == nil {
		return nil, errors.New("не указан logger")
	}
	
	return &repository{
		db:     db,
		logger: logger,
	}, nil
}

func (r *repository) getById(ctx context.Context, id string) (*User, error) {
	if id == "" {
		return nil, errors.New("не указан id")
	}

	user := &User{}

	err := r.db.GetContext(ctx, user, stmtGetUserById, id)

	if err != nil {
		errorMsg := fmt.Sprintf("ошибка получения пользователя с id %s", id)
		r.logger.Error().Err(err).Msg(errorMsg)
		return nil, fmt.Errorf("%s: %w", errorMsg, err)
	}
	
	return user, nil
}

func (r *repository) getByEmail(ctx context.Context, email string) (*User, error) {
	if email == "" {
		return nil, errors.New("не указан email")
	}

	user := &User{}

	err := r.db.GetContext(ctx, user, stmtGetUserByEmail, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("пользователь с email %s не найден", email)
		}
		errorMsg := fmt.Sprintf("ошибка получения пользователя с email %s", email)
		r.logger.Error().Err(err).Msg(errorMsg)
		return nil, fmt.Errorf(errorMsg)
	}
	
	return user, nil
}

func (r *repository) userCreate(ctx context.Context, user *User) error {
	_, err := r.db.NamedExecContext(ctx, stmtCreateUser, user)
	if err != nil {
		errorMsg := fmt.Sprintf("ошибка создания пользователя: %+v", user)
		r.logger.Error().Err(err).Msg(errorMsg)
		return errors.New(errorMsg)
	}
	return nil
}

func (r *repository) userUpdate(ctx context.Context, user *User) error {
	_, err := r.db.NamedExecContext(ctx, stmtUpdateUser, user)
	if err != nil {
		errorMsg := fmt.Sprintf("ошибка редактирования пользователя: %+v", user)
		r.logger.Error().Err(err).Msg(errorMsg)
		return errors.New(errorMsg)
	}
	return nil
}

func (r *repository) getUserSessions(ctx context.Context, userId string) ([]session, error) {
	if userId == "" {
		return nil, errors.New("не указан id пользователя")
	}

	var sessions []session

	err := r.db.SelectContext(ctx, &sessions, stmtGetUserSessions, userId)
	if err != nil {
		errorMsg := fmt.Sprintf("ошибка получения сессий пользователя: %s", userId)
		r.logger.Error().Err(err).Msg(errorMsg)
		return nil, errors.New(errorMsg)
	}

	return sessions, nil
}

func (r *repository) getSessionByID(ctx context.Context, sessionId string) (*session, error) {
	session := &session{}
	if err := r.db.GetContext(ctx, session, stmtGetSessionByID, sessionId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("сессия с id %s не найдена", sessionId)
		} else {
			errorMsg := fmt.Sprintf("ошибка получения сессии с id %s", sessionId)
			r.logger.Error().Err(err).Msg(errorMsg)
			return nil, errors.New(errorMsg)
		}
	}
	return session, nil
}

func (r *repository) createSession(ctx context.Context, session *session) error {
	if session == nil {
        return errors.New("не указана сессия")
    }

    existingSession, err := r.getSessionByID(ctx, session.ID)
    if err != nil && !errors.Is(err, sql.ErrNoRows) {
        return fmt.Errorf("ошибка проверки сессии: %w", err)
    }

    if existingSession != nil {
        if existingSession.UserID != session.UserID {
            return errors.New("пользователь не соответствует сессии")
        }
        return nil
    }

    _, err = r.db.NamedExecContext(ctx, stmtCreateSession, session)
    if err != nil {
        errorMsg := fmt.Sprintf("ошибка создания сессии для пользователя %s", session.UserID)
        r.logger.Error().Err(err).Msg(errorMsg)
        return fmt.Errorf("%s: %w", errorMsg, err)
    }
    
    return nil
}