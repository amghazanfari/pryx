package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/amghazanfari/pryx/utils"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB            *sql.DB
	BytesPerToken int
	Duration      time.Duration
}

func (service *PasswordResetService) Create(email string) (*PasswordReset, error) {
	email = strings.ToLower(email)
	bytesPerToken := service.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	duration := service.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}
	token, err := utils.String(bytesPerToken)
	if err != nil {
		return nil, err
	}

	tokenHash := service.hash(token)
	fmt.Printf("token hash created: %s\n", tokenHash)

	row := service.DB.QueryRow(`
	SELECT id
	FROM users
	WHERE email = $1
	`, email)

	var userID int
	err = row.Scan(&userID)
	if err != nil {
		return nil, err
	}

	passwordReset := PasswordReset{
		UserID:    userID,
		Token:     token,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(time.Duration(duration)),
	}

	row = service.DB.QueryRow(`
	INSERT INTO password_reset (user_id, token_hash, expires_at)
	VALUES ($1, $2, $3) 
	ON CONFLICT (user_id) DO UPDATE
	SET user_id = $1, token_hash = $2, expires_at = $3
	RETURNING id
	`, passwordReset.UserID, passwordReset.TokenHash, passwordReset.ExpiresAt)

	err = row.Scan(&passwordReset.ID)

	if err != nil {
		return nil, fmt.Errorf("reset password: %w", err)
	}
	return &passwordReset, nil
}

func (service *PasswordResetService) Consume(token string) (*User, error) {
	tokenHash := service.hash(token)

	row := service.DB.QueryRow(`
	SELECT user_id, expires_at FROM password_reset 
	WHERE token_hash = $1
	`, tokenHash)

	var userID int
	var expiresAt time.Time
	err := row.Scan(&userID, &expiresAt)
	if err != nil {
		return nil, fmt.Errorf("consume token: %w", err)
	}

	if time.Now().Compare(expiresAt) > 0 {
		return nil, fmt.Errorf("the token is expired")
	}

	err = service.delete(token)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}

	user := User{
		ID: userID,
	}
	return &user, nil
}

func (service *PasswordResetService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

func (service *PasswordResetService) delete(token string) error {
	tokenHash := service.hash(token)

	_, err := service.DB.Exec(`
	DELETE FROM password_reset 
	WHERE token_hash = $1
	`, tokenHash)

	return err
}
