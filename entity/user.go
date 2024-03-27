package entity

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Uuid uuid.UUID `db:"uuid" json:"uuid"`
	Name string    `db:"name" json:"name"`
}

type NewUser struct {
	Name string `json:"name"`
}

func (*Repository) GetAllUsers(tx *sqlx.Tx) ([]User, error) {
	users := []User{}
	err := tx.Select(&users, "SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("could not get all users from db: %w", err)
	}

	return users, nil
}

func (repo *Repository) GetUser(tx *sqlx.Tx, userUuid uuid.UUID) (*User, error) {
	var user User
	if err := tx.Get(&user, "SELECT * FROM users WHERE uuid = $1", userUuid); err != nil {
		return nil, fmt.Errorf("failed to get user %s: %w", userUuid, err)
	}

	return &user, nil
}

func (repo *Repository) CreateUser(tx *sqlx.Tx, user *NewUser) (*User, error) {
	var createdUser User
	err := tx.Get(&createdUser, "INSERT INTO users (name) VALUES ($1) RETURNING uuid, name", user.Name)
	if err != nil {
		return nil, fmt.Errorf("could not create user %s: %w", user.Name, err)
	}
	return &createdUser, nil
}

func (repo *Repository) UpdateUser(tx *sqlx.Tx, userUuid uuid.UUID, user *NewUser) (*User, error) {
	_, err := tx.Exec("UPDATE users SET name = $2, url = $3 WHERE uuid = $1", userUuid, user.Name)
	if err != nil {
		return nil, fmt.Errorf("could not update user %s: %w", userUuid, err)
	}

	return repo.GetUser(tx, userUuid)
}

func (repo *Repository) DeleteUser(tx *sqlx.Tx, userUuid uuid.UUID) error {
	_, err := tx.Exec("DELETE FROM users WHERE uuid = $1", userUuid)
	if err != nil {
		return fmt.Errorf("could not delete user %s: %w", userUuid, err)
	}

	return nil
}
