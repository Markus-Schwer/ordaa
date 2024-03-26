package entity

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Uuid uuid.UUID `db:"uuid"`
	Name string    `db:"name"`
}

type NewUser struct {
	Name string
}

func (*Repository) GetAllUsers(tx *sqlx.Tx) ([]User, error) {
	var users []User
	if err := tx.Select(&users, "SELECT * FROM users"); err != nil {
		return nil, fmt.Errorf("could not query users: %w", err)
	}

	return users, nil
}

func (*Repository) GetUser(tx *sqlx.Tx, uuid uuid.UUID) (*User, error) {
  var user User
  if err := tx.Get(&user, "SELECT * FROM users WHERE id = $1", uuid); err != nil {
      return nil, fmt.Errorf("could not get user %s: %w", uuid, err)
  }

  return &user, nil
}
