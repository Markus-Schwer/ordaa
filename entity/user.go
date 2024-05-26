package entity

import (
	"fmt"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type User struct {
	Uuid *uuid.UUID `gorm:"column:uuid;primaryKey" json:"uuid"`
	Name string     `gorm:"column:name" json:"name"`
}

type MatrixUser struct {
	Uuid     *uuid.UUID `gorm:"column:uuid;primaryKey" json:"uuid"`
	UserUuid *uuid.UUID  `gorm:"column:user_uuid" json:"user_uuid"`
	Username string     `gorm:"column:username" json:"username"`
}

type PasswordUser struct {
	Uuid     *uuid.UUID `gorm:"column:uuid;primaryKey" json:"uuid"`
	UserUuid *uuid.UUID  `gorm:"column:user_uuid" json:"user_uuid"`
	Username string     `gorm:"column:username" json:"username"`
	Password string     `gorm:"column:password" json:"password"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("could not create uuid: %w", err)
	}

	user.Uuid = &newUuid
	return nil
}

func (matrixUser *MatrixUser) BeforeCreate(tx *gorm.DB) (err error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("could not create uuid: %w", err)
	}

	matrixUser.Uuid = &newUuid
	return nil
}

func (passwordUser *PasswordUser) BeforeCreate(tx *gorm.DB) (err error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("could not create uuid: %w", err)
	}

	passwordUser.Uuid = &newUuid
	return nil
}

func (*RepositoryImpl) GetAllUsers(tx *gorm.DB) ([]User, error) {
	users := []User{}
	err := tx.Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("could not get all users from db: %w", err)
	}

	return users, nil
}

func (repo *RepositoryImpl) GetUser(tx *gorm.DB, userUuid *uuid.UUID) (*User, error) {
	var user User
	if err := tx.First(&user, userUuid).Error; err != nil {
		return nil, fmt.Errorf("failed to get user %s: %w", userUuid, err)
	}

	return &user, nil
}

func (repo *RepositoryImpl) CreateUser(tx *gorm.DB, user *User) (*User, error) {
	err := tx.Create(&user).Error
	if err != nil {
		return nil, fmt.Errorf("could not create user %s: %w", user.Name, err)
	}
	return user, nil
}

func (repo *RepositoryImpl) UpdateUser(tx *gorm.DB, userUuid *uuid.UUID, user *User) (*User, error) {
	foundUser, err := repo.GetUser(tx, userUuid)
	if err != nil {
		return nil, fmt.Errorf("could not update user %s: %w", userUuid, err)
	}
	foundUser.Name = user.Name
	err = tx.Save(&foundUser).Error
	if err != nil {
		return nil, fmt.Errorf("could not update user %s: %w", userUuid, err)
	}

	return foundUser, nil
}

func (repo *RepositoryImpl) DeleteUser(tx *gorm.DB, userUuid *uuid.UUID) error {
	err := tx.Delete(&User{Uuid: userUuid}).Error
	if err != nil {
		return fmt.Errorf("could not delete user %s: %w", userUuid, err)
	}

	// TODO: check if the user is a password or matrix user
	repo.DeleteMatrixUser(tx, userUuid)
	repo.DeletePasswordUser(tx, userUuid)

	return nil
}

func (*RepositoryImpl) GetAllMatrixUsers(tx *gorm.DB) ([]MatrixUser, error) {
	matrixUsers := []MatrixUser{}
	err := tx.Find(&matrixUsers).Error
	if err != nil {
		return nil, fmt.Errorf("could not get all users from db: %w", err)
	}

	return matrixUsers, nil
}

func (repo *RepositoryImpl) GetMatrixUser(tx *gorm.DB, matrixUserUuid *uuid.UUID) (*MatrixUser, error) {
	var matrixUser MatrixUser
	if err := tx.Where(&MatrixUser{Uuid: matrixUserUuid}).First(&matrixUser).Error; err != nil {
		return nil, fmt.Errorf("failed to get user %s: %w", matrixUserUuid, err)
	}

	return &matrixUser, nil
}

func (repo *RepositoryImpl) CreateMatrixUser(tx *gorm.DB, matrixUser *MatrixUser) (*MatrixUser, error) {
	err := tx.Create(&matrixUser).Error
	if err != nil {
		return nil, fmt.Errorf("could not create matrix user %s: %w", matrixUser.Username, err)
	}
	return matrixUser, nil
}

func (repo *RepositoryImpl) UpdateMatrixUser(tx *gorm.DB, matrixUserUuid *uuid.UUID, matrixUser *MatrixUser) (*MatrixUser, error) {
	existingMatrixUser, err := repo.GetMatrixUser(tx, matrixUserUuid)
	if err != nil {
		return nil, fmt.Errorf("could not update user %s: %w", matrixUserUuid, err)
	}
	existingMatrixUser.Username = matrixUser.Username
	existingMatrixUser.UserUuid = matrixUser.UserUuid
	err = tx.Save(&existingMatrixUser).Error
	if err != nil {
		return nil, fmt.Errorf("could not update user %s: %w", matrixUserUuid, err)
	}

	return existingMatrixUser, nil
}

func (repo *RepositoryImpl) DeleteMatrixUser(tx *gorm.DB, userUuid *uuid.UUID) error {
	err := tx.Where(&MatrixUser{UserUuid: userUuid}).Delete(&MatrixUser{}).Error
	if err != nil {
		return fmt.Errorf("could not delete user %s: %w", userUuid, err)
	}

	return nil
}

func (*RepositoryImpl) GetAllPasswordUsers(tx *gorm.DB) ([]PasswordUser, error) {
	passwordUsers := []PasswordUser{}
	err := tx.Find(&passwordUsers).Error
	if err != nil {
		return nil, fmt.Errorf("could not get all users from db: %w", err)
	}

	return passwordUsers, nil
}

func (repo *RepositoryImpl) FindPasswordUser(tx *gorm.DB, username string) (*PasswordUser, error) {
	var passwordUser PasswordUser
	if err := tx.Where(&PasswordUser{Username: username}).Find(&passwordUser).Error; err != nil {
		return nil, fmt.Errorf("failed to get user %s: %w", username, err)
	}

	return &passwordUser, nil
}

func (repo *RepositoryImpl) GetPasswordUser(tx *gorm.DB, passwordUserUuid *uuid.UUID) (*PasswordUser, error) {
	var passwordUser PasswordUser
	if err := tx.Where(&PasswordUser{Uuid: passwordUserUuid}).First(&passwordUser).Error; err != nil {
		return nil, fmt.Errorf("failed to get user %s: %w", passwordUserUuid, err)
	}

	return &passwordUser, nil
}

func (repo *RepositoryImpl) CreatePasswordUser(tx *gorm.DB, passwordUser *PasswordUser) (*PasswordUser, error) {
	err := tx.Create(&passwordUser).Error
	if err != nil {
		return nil, fmt.Errorf("could not create password user %s: %w", passwordUser.Username, err)
	}
	return passwordUser, nil
}

func (repo *RepositoryImpl) UpdatePasswordUser(tx *gorm.DB, passwordUserUuid *uuid.UUID, passwordUser *PasswordUser) (*PasswordUser, error) {
	existingPasswordUser, err := repo.GetPasswordUser(tx, passwordUserUuid)
	if err != nil {
		return nil, fmt.Errorf("could not update user %s: %w", passwordUserUuid, err)
	}

	existingPasswordUser.Username = passwordUser.Username
	existingPasswordUser.UserUuid = passwordUser.UserUuid
	existingPasswordUser.Password = passwordUser.Password
	err = tx.Save(existingPasswordUser).Error
	if err != nil {
		return nil, fmt.Errorf("could not update user %s: %w", passwordUserUuid, err)
	}

	return existingPasswordUser, nil
}

func (repo *RepositoryImpl) DeletePasswordUser(tx *gorm.DB, userUuid *uuid.UUID) error {
	err := tx.Where(&PasswordUser{UserUuid: userUuid}).Delete(&PasswordUser{}).Error
	if err != nil {
		return fmt.Errorf("could not delete user %s: %w", userUuid, err)
	}

	return nil
}
