package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/amghazanfari/pryx/utils"
)

const (
	MinBytesPerToken = 32
)

type Session struct {
	ID     int
	UserID int
	// Token will only set in creating new session
	Token     string
	TokenHash string
}

type SessionService struct {
	DB            *sql.DB
	BytesPerToken int
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := utils.String(bytesPerToken)
	if err != nil {
		return nil, err
	}

	tokenHash := ss.hash(token)
	fmt.Printf("token hash created: %s\n", tokenHash)

	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: tokenHash,
	}

	row := ss.DB.QueryRow(`
	INSERT INTO sessions (user_id, token_hash)
	VALUES ($1, $2) 
	ON CONFLICT (user_id) DO UPDATE
	SET user_id = $1, token_hash = $2
	RETURNING id
	`, userID, tokenHash)

	err = row.Scan(&session.ID)

	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	tokenHash := ss.hash(token)
	fmt.Printf("looking for user with token_hash: %s\n", tokenHash)

	row := ss.DB.QueryRow(`
	SELECT users.id, users.email, users.password_hash FROM users
	JOIN sessions ON users.id = sessions.user_id
	WHERE token_hash = $1 LIMIT 1
	`, tokenHash)

	var user User

	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)
	fmt.Printf("trying to delete session for user with token_hash: %s\n", tokenHash)

	_, err := ss.DB.Exec(`
	DELETE FROM sessions
	WHERE token_hash = $1
	`, tokenHash)

	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
