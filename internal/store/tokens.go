package store

import (
	"database/sql"
	"go-server/internal/tokens"
	"time"
)

type TokenStore interface {
	Insert(token *tokens.Token) error
	CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int, scope string) error
}

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		db: db,
	}
}

func (t *PostgresTokenStore) CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error) {
	tokens, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = t.Insert(tokens)
	return tokens, err
}

func (t *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope) 
		VALUES ($1, $2, $3, $4)
	`

	_, err := t.db.Exec(query, token.Hash, token.UserId, token.Expiry, token.Scope)

	return err
}

func (t *PostgresTokenStore) DeleteAllTokensForUser(userID int, scope string) error {
	query := `
		DELETE FROM tokens WHERE user_id = $1 AND scope = $2
	`

	result, err := t.db.Exec(query, userID, scope)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}
