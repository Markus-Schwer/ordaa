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

type MatrixUser struct {
	Uuid     uuid.UUID `db:"uuid" json:"uuid"`
	UserUuid uuid.UUID `db:"user_uuid" json:"user_uuid"`
	Username string    `db:"username" json:"username"`
}

type NewMatrixUser struct {
	UserUuid uuid.UUID `db:"user_uuid" json:"user_uuid"`
	Username string    `db:"username" json:"username"`
}

type PasswordUser struct {
	Uuid     uuid.UUID `db:"uuid" json:"uuid"`
	UserUuid uuid.UUID `db:"user_uuid" json:"user_uuid"`
	Username string    `db:"username" json:"username"`
	Password string    `db:"password" json:"-"`
}

type NewPasswordUser struct {
	UserUuid uuid.UUID `db:"user_uuid" json:"-"`
	Username string    `db:"username" json:"username"`
	Password string    `db:"password" json:"password"`
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
	_, err := tx.Exec("UPDATE users SET name = $2 WHERE uuid = $1", userUuid, user.Name)
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

func (*Repository) GetAllMatrixUsers(tx *sqlx.Tx) ([]MatrixUser, error) {
	matrixUsers := []MatrixUser{}
	err := tx.Select(&matrixUsers, "SELECT * FROM matrix_users")
	if err != nil {
		return nil, fmt.Errorf("could not get all users from db: %w", err)
	}

	return matrixUsers, nil
}

func (repo *Repository) GetMatrixUser(tx *sqlx.Tx, matrixUserUuid uuid.UUID) (*MatrixUser, error) {
	var matrixUser MatrixUser
	if err := tx.Get(&matrixUser, "SELECT * FROM matrix_users WHERE uuid = $1", matrixUserUuid); err != nil {
		return nil, fmt.Errorf("failed to get user %s: %w", matrixUserUuid, err)
	}

	return &matrixUser, nil
}

func (repo *Repository) CreateMatrixUser(tx *sqlx.Tx, matrixUser *NewMatrixUser) (*MatrixUser, error) {
	var createdMatrixUser MatrixUser
	err := tx.Get(&createdMatrixUser, "INSERT INTO matrix_users (user_uuid, username) VALUES ($1, $2) RETURNING uuid, name", matrixUser.UserUuid, matrixUser.Username)
	if err != nil {
		return nil, fmt.Errorf("could not create user %s: %w", matrixUser.Username, err)
	}
	return &createdMatrixUser, nil
}

func (repo *Repository) UpdateMatrixUser(tx *sqlx.Tx, matrixUserUuid uuid.UUID, matrixUser *NewMatrixUser) (*MatrixUser, error) {
	_, err := tx.Exec("UPDATE users SET username = $2, user_uuid = $3 WHERE uuid = $1", matrixUserUuid, matrixUser.Username, matrixUser.UserUuid)
	if err != nil {
		return nil, fmt.Errorf("could not update user %s: %w", matrixUserUuid, err)
	}

	return repo.GetMatrixUser(tx, matrixUserUuid)
}

func (repo *Repository) DeleteMatrixUser(tx *sqlx.Tx, matrixUserUuid uuid.UUID) error {
	_, err := tx.Exec("DELETE FROM matrix_users WHERE uuid = $1", matrixUserUuid)
	if err != nil {
		return fmt.Errorf("could not delete user %s: %w", matrixUserUuid, err)
	}

	return nil
}

func (*Repository) GetAllPasswordUsers(tx *sqlx.Tx) ([]PasswordUser, error) {
	passwordUsers := []PasswordUser{}
	err := tx.Select(&passwordUsers, "SELECT * FROM password_users")
	if err != nil {
		return nil, fmt.Errorf("could not get all users from db: %w", err)
	}

	return passwordUsers, nil
}

func (repo *Repository) FindPasswordUser(tx *sqlx.Tx, username string) (*PasswordUser, error) {
	var passwordUser PasswordUser
	if err := tx.Get(&passwordUser, "SELECT * FROM password_users WHERE username = $1", username); err != nil {
		return nil, fmt.Errorf("failed to get user %s: %w", username, err)
	}

	return &passwordUser, nil
}

func (repo *Repository) GetPasswordUser(tx *sqlx.Tx, passwordUserUuid uuid.UUID) (*PasswordUser, error) {
	var passwordUser PasswordUser
	if err := tx.Get(&passwordUser, "SELECT * FROM password_users WHERE uuid = $1", passwordUserUuid); err != nil {
		return nil, fmt.Errorf("failed to get user %s: %w", passwordUserUuid, err)
	}

	return &passwordUser, nil
}

func (repo *Repository) CreatePasswordUser(tx *sqlx.Tx, passwordUser *NewPasswordUser) (*PasswordUser, error) {
	var createdPasswordUser PasswordUser
	err := tx.Get(&createdPasswordUser, "INSERT INTO password_users (user_uuid, username, password) VALUES ($1, $2, $3) RETURNING *", passwordUser.UserUuid, passwordUser.Username, passwordUser.Password)
	if err != nil {
		return nil, fmt.Errorf("could not create user %s: %w", passwordUser.Username, err)
	}
	return &createdPasswordUser, nil
}

func (repo *Repository) UpdatePasswordUser(tx *sqlx.Tx, passwordUserUuid uuid.UUID, passwordUser *NewPasswordUser) (*PasswordUser, error) {
	_, err := tx.Exec("UPDATE users SET username = $2, password = $3, user_uuid = $4 WHERE uuid = $1", passwordUserUuid, passwordUser.UserUuid, passwordUser.Username, passwordUser.Password)
	if err != nil {
		return nil, fmt.Errorf("could not update user %s: %w", passwordUserUuid, err)
	}

	return repo.GetPasswordUser(tx, passwordUserUuid)
}

func (repo *Repository) DeletePasswordUser(tx *sqlx.Tx, passwordUserUuid uuid.UUID) error {
	_, err := tx.Exec("DELETE FROM password_users WHERE uuid = $1", passwordUserUuid)
	if err != nil {
		return fmt.Errorf("could not delete user %s: %w", passwordUserUuid, err)
	}

	return nil
}
