package entity

import (
	"errors"
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
	UserUuid *uuid.UUID `gorm:"column:user_uuid" json:"user_uuid"`
	Username string     `gorm:"column:username" json:"username"`
}

type PasswordUser struct {
	Uuid     *uuid.UUID `gorm:"column:uuid;primaryKey" json:"uuid"`
	UserUuid *uuid.UUID `gorm:"column:user_uuid" json:"user_uuid"`
	Username string     `gorm:"column:username" json:"username" validate:"required"`
	Password string     `gorm:"column:password" json:"password" validate:"required"`
}

type SshUser struct {
	Uuid      *uuid.UUID `gorm:"column:uuid;primaryKey" json:"uuid"`
	UserUuid  *uuid.UUID `gorm:"column:user_uuid" json:"user_uuid"`
	PublicKey string     `gorm:"column:public_key" json:"public_key"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotCreatUuid, err)
	}

	user.Uuid = &newUuid
	return nil
}

func (matrixUser *MatrixUser) BeforeCreate(tx *gorm.DB) (err error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotCreatUuid, err)
	}

	matrixUser.Uuid = &newUuid
	return nil
}

func (passwordUser *PasswordUser) BeforeCreate(tx *gorm.DB) (err error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotCreatUuid, err)
	}

	passwordUser.Uuid = &newUuid
	return nil
}

func (*RepositoryImpl) GetAllUsers(tx *gorm.DB) ([]User, error) {
	users := []User{}
	err := tx.Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCannotGetAllUsers, err)
	}

	return users, nil
}

func (repo *RepositoryImpl) GetUser(tx *gorm.DB, userUuid *uuid.UUID) (*User, error) {
	var user User
	err := tx.First(&user, userUuid).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingUser, err)
	}

	return &user, nil
}

func (repo *RepositoryImpl) GetUserByName(tx *gorm.DB, name string) (*User, error) {
	var user User
	err := tx.Where(&User{Name: name}).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingUser, err)
	}

	return &user, nil
}

func (repo *RepositoryImpl) CreateUser(tx *gorm.DB, user *User) (*User, error) {
	err := tx.Create(&user).Error
	if err != nil {
		return nil, fmt.Errorf("%: %w", ErrCreatingUser, err)
	}
	return user, nil
}

func (repo *RepositoryImpl) UpdateUser(tx *gorm.DB, userUuid *uuid.UUID, user *User) (*User, error) {
	foundUser, err := repo.GetUser(tx, userUuid)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUpdatingUser, err)
	}
	foundUser.Name = user.Name
	err = tx.Save(&foundUser).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUpdatingUser, err)
	}

	return foundUser, nil
}

func (repo *RepositoryImpl) DeleteUser(tx *gorm.DB, userUuid *uuid.UUID) error {
	err := tx.Delete(&User{Uuid: userUuid}).Error
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDeletingUser, err)
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
		return nil, fmt.Errorf("%w: %w", ErrCannotGetAllUsers, err)
	}

	return matrixUsers, nil
}

func (repo *RepositoryImpl) GetMatrixUser(tx *gorm.DB, matrixUserUuid *uuid.UUID) (*MatrixUser, error) {
	var matrixUser MatrixUser
	err := tx.Where(&MatrixUser{Uuid: matrixUserUuid}).First(&matrixUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingUser, err)
	}

	return &matrixUser, nil
}

func (repo *RepositoryImpl) GetMatrixUserByUsername(tx *gorm.DB, username string) (*MatrixUser, error) {
	var matrixUser MatrixUser
	err := tx.Where(&MatrixUser{Username: username}).First(&matrixUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingUser, err)
	}

	return &matrixUser, nil
}

func (repo *RepositoryImpl) CreateMatrixUser(tx *gorm.DB, matrixUser *MatrixUser) (*MatrixUser, error) {
	err := tx.Create(&matrixUser).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreatingUser, err)
	}
	return matrixUser, nil
}

func (repo *RepositoryImpl) UpdateMatrixUser(tx *gorm.DB, matrixUserUuid *uuid.UUID, matrixUser *MatrixUser) (*MatrixUser, error) {
	existingMatrixUser, err := repo.GetMatrixUser(tx, matrixUserUuid)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUpdatingUser, err)
	}
	existingMatrixUser.Username = matrixUser.Username
	existingMatrixUser.UserUuid = matrixUser.UserUuid
	err = tx.Save(&existingMatrixUser).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUpdatingUser, err)
	}

	return existingMatrixUser, nil
}

func (repo *RepositoryImpl) DeleteMatrixUser(tx *gorm.DB, userUuid *uuid.UUID) error {
	err := tx.Where(&MatrixUser{UserUuid: userUuid}).Delete(&MatrixUser{}).Error
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDeletingUser, err)
	}

	return nil
}

func (*RepositoryImpl) GetAllPasswordUsers(tx *gorm.DB) ([]PasswordUser, error) {
	passwordUsers := []PasswordUser{}
	err := tx.Find(&passwordUsers).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCannotGetAllUsers, err)
	}

	return passwordUsers, nil
}

func (repo *RepositoryImpl) FindPasswordUser(tx *gorm.DB, username string) (*PasswordUser, error) {
	var passwordUser PasswordUser
	err := tx.Where(&PasswordUser{Username: username}).Find(&passwordUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingUser, err)
	}

	return &passwordUser, nil
}

func (repo *RepositoryImpl) GetPasswordUser(tx *gorm.DB, passwordUserUuid *uuid.UUID) (*PasswordUser, error) {
	var passwordUser PasswordUser
	err := tx.Where(&PasswordUser{Uuid: passwordUserUuid}).First(&passwordUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingUser, err)
	}

	return &passwordUser, nil
}

func (repo *RepositoryImpl) CreatePasswordUser(tx *gorm.DB, passwordUser *PasswordUser) (*PasswordUser, error) {
	err := tx.Create(&passwordUser).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreatingUser, err)
	}
	return passwordUser, nil
}

func (repo *RepositoryImpl) UpdatePasswordUser(tx *gorm.DB, passwordUserUuid *uuid.UUID, passwordUser *PasswordUser) (*PasswordUser, error) {
	existingPasswordUser, err := repo.GetPasswordUser(tx, passwordUserUuid)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUpdatingUser, err)
	}

	existingPasswordUser.Username = passwordUser.Username
	existingPasswordUser.UserUuid = passwordUser.UserUuid
	existingPasswordUser.Password = passwordUser.Password
	err = tx.Save(existingPasswordUser).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUpdatingUser, err)
	}

	return existingPasswordUser, nil
}

func (repo *RepositoryImpl) DeletePasswordUser(tx *gorm.DB, userUuid *uuid.UUID) error {
	err := tx.Where(&PasswordUser{UserUuid: userUuid}).Delete(&PasswordUser{}).Error
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDeletingUser, err)
	}

	return nil
}

func (*RepositoryImpl) GetAllSshUsers(tx *gorm.DB) ([]SshUser, error) {
	sshUsers := []SshUser{}
	err := tx.Find(&sshUsers).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCannotGetAllUsers, err)
	}

	return sshUsers, nil
}

func (repo *RepositoryImpl) GetSshUser(tx *gorm.DB, sshUserUuid *uuid.UUID) (*SshUser, error) {
	var sshUser SshUser
	err := tx.Where(&SshUser{Uuid: sshUserUuid}).First(&sshUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingUser, err)
	}

	return &sshUser, nil
}

func (repo *RepositoryImpl) GetSshUserByPublicKey(tx *gorm.DB, publicKey string) (*SshUser, error) {
	var sshUser SshUser
	err := tx.Where(&SshUser{PublicKey: publicKey}).First(&sshUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingUser, err)
	}

	return &sshUser, nil
}

func (repo *RepositoryImpl) CreateSshUser(tx *gorm.DB, sshUser *SshUser) (*SshUser, error) {
	err := tx.Create(&sshUser).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreatingUser, err)
	}
	return sshUser, nil
}

func (repo *RepositoryImpl) UpdateSshUser(tx *gorm.DB, sshUserUuid *uuid.UUID, sshUser *SshUser) (*SshUser, error) {
	existingSshUser, err := repo.GetSshUser(tx, sshUserUuid)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUpdatingUser, err)
	}
	existingSshUser.PublicKey = sshUser.PublicKey
	existingSshUser.UserUuid = sshUser.UserUuid
	err = tx.Save(&existingSshUser).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUpdatingUser, err)
	}

	return existingSshUser, nil
}

func (repo *RepositoryImpl) DeleteSshUser(tx *gorm.DB, userUuid *uuid.UUID) error {
	err := tx.Where(&SshUser{UserUuid: userUuid}).Delete(&SshUser{}).Error
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDeletingUser, err)
	}

	return nil
}
