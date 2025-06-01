package model

import (
	"github.com/gofrs/uuid"
)

type User struct {
	UUID *uuid.UUID
	Name string
}

type MatrixUser struct {
	UUID     *uuid.UUID
	UserUUID *uuid.UUID
	Username string
}

type PasswordUser struct {
	UUID     *uuid.UUID
	UserUUID *uuid.UUID
	Username string
	Password string
}

type SSHUser struct {
	UUID      *uuid.UUID
	UserUUID  *uuid.UUID
	PublicKey string
}
