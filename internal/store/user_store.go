package store

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}

	p.plainText = &plainTextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			{
				return false, nil
			}

		default:
			{
				return false, err
			}
		}
	}

	return true, nil
}

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
	Bio          string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserStore interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

func (u *PostgresUserStore) CreateUser(user *User) error {
	query := `
		INSERT INTO users (username, email, password_hash, bio)
		values ($1, $2, $3, $4)
		RETURNING ID;
	`

	err := u.db.QueryRow(query, user.Username, user.Email, user.PasswordHash.hash, user.Bio).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *PostgresUserStore) GetUserByEmail(email string) (*User, error) {
	var user = &User{
		PasswordHash: password{},
	}
	query := `
		SELECT id, username, email, password_hash, bio FROM users WHERE email = $1
	`

	err := u.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Bio,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *PostgresUserStore) UpdateUser(user *User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, bio = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING updated_at
	`

	result, err := u.db.Exec(query, user.Username, user.Email, user.Bio, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
