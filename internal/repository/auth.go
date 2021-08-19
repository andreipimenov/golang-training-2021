package repository

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"

	"github.com/google/uuid"

	"github.com/andreipimenov/golang-training-2021/internal/model"
)

type user struct {
	id    uuid.UUID
	name  string
	email sql.NullString
}

type Auth struct {
	*sql.DB
}

func NewAuth(db *sql.DB) *Auth {
	return &Auth{db}
}

func (db *Auth) GetUser(username string, password string) (*model.User, error) {
	var u user
	err := db.QueryRow("SELECT id, name, email FROM users WHERE username = $1 AND password = $2", username, hash(password)).Scan(&u.id, &u.name, &u.email)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		ID:   u.id,
		Name: u.name,
	}
	if u.email.Valid {
		user.Email = u.email.String
	} else {
		user.Email = "default@email.com"
	}
	return user, nil
}

func (db *Auth) SaveToken(userID uuid.UUID, token string) error {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO tokens (id, token, user_id) VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE SET token=EXCLUDED.token", tokenID, token, userID)
	if err != nil {
		return err
	}
	return nil
}

func (db *Auth) GetTokenByUserID(userID uuid.UUID) (string, error) {
	var token string
	err := db.QueryRow("SELECT token FROM tokens WHERE user_id = $1", userID).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}

func hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
