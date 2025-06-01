package entity

import (
	"fmt"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type User struct {
	UUID *uuid.UUID `gorm:"column:uuid;primaryKey" json:"uuid"`
	Name string     `gorm:"column:name" json:"name"`
}

type MatrixUser struct {
	UUID     *uuid.UUID `gorm:"column:uuid;primaryKey" json:"uuid"`
	UserUUID *uuid.UUID `gorm:"column:user_uuid" json:"user_uuid"`
	Username string     `gorm:"column:username" json:"username"`
}

type PasswordUser struct {
	UUID     *uuid.UUID `gorm:"column:uuid;primaryKey" json:"uuid"`
	UserUUID *uuid.UUID `gorm:"column:user_uuid" json:"user_uuid"`
	Username string     `gorm:"column:username" json:"username" validate:"required"`
	Password string     `gorm:"column:password" json:"password" validate:"required"`
}

type SSHUser struct {
	UUID      *uuid.UUID `gorm:"column:uuid;primaryKey" json:"uuid"`
	UserUUID  *uuid.UUID `gorm:"column:user_uuid" json:"user_uuid"`
	PublicKey string     `gorm:"column:public_key" json:"public_key"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	newUUID, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotCreatUUID, err)
	}

	user.UUID = &newUUID

	return nil
}

func (matrixUser *MatrixUser) BeforeCreate(tx *gorm.DB) (err error) {
	newUUID, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotCreatUUID, err)
	}

	matrixUser.UUID = &newUUID

	return nil
}

func (passwordUser *PasswordUser) BeforeCreate(tx *gorm.DB) (err error) {
	newUUID, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotCreatUUID, err)
	}

	passwordUser.UUID = &newUUID

	return nil
}
