package service

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/Markus-Schwer/ordaa/internal/entity"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	GetUser(ctx context.Context, uuid *uuid.UUID) (*entity.User, error)
	GetUserByName(ctx context.Context, name string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateUser(ctx context.Context, uuid *uuid.UUID, user *entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, uuid *uuid.UUID) error

	GetAllMatrixUsers(ctx context.Context) ([]entity.MatrixUser, error)
	GetMatrixUser(ctx context.Context, uuid *uuid.UUID) (*entity.MatrixUser, error)
	GetMatrixUserByUsername(ctx context.Context, username string) (*entity.MatrixUser, error)
	CreateMatrixUser(ctx context.Context, matrixUser *entity.MatrixUser) (*entity.MatrixUser, error)
	UpdateMatrixUser(ctx context.Context, uuid *uuid.UUID, matrixUser *entity.MatrixUser) (*entity.MatrixUser, error)
	DeleteMatrixUser(ctx context.Context, uuid *uuid.UUID) error

	GetAllPasswordUsers(ctx context.Context) ([]entity.PasswordUser, error)
	FindPasswordUser(ctx context.Context, username string) (*entity.PasswordUser, error)
	GetPasswordUser(ctx context.Context, uuid *uuid.UUID) (*entity.PasswordUser, error)
	CreatePasswordUser(ctx context.Context, passwordUser *entity.PasswordUser) (*entity.PasswordUser, error)
	UpdatePasswordUser(ctx context.Context, uuid *uuid.UUID, passwordUser *entity.PasswordUser) (*entity.PasswordUser, error)
	DeletePasswordUser(ctx context.Context, uuid *uuid.UUID) error

	GetAllSSHUsers(ctx context.Context) ([]entity.SSHUser, error)
	GetSSHUser(ctx context.Context, uuid *uuid.UUID) (*entity.SSHUser, error)
	GetSSHUserByPublicKey(ctx context.Context, publicKey string) (*entity.SSHUser, error)
	CreateSSHUser(ctx context.Context, sshUser *entity.SSHUser) (*entity.SSHUser, error)
	UpdateSSHUser(ctx context.Context, uuid *uuid.UUID, sshUser *entity.SSHUser) (*entity.SSHUser, error)
	DeleteSSHUser(ctx context.Context, uuid *uuid.UUID) error

	RegisterMatrixUser(ctx context.Context, username string) (*entity.User, error)
	SetPublicKey(ctx context.Context, userUUID *uuid.UUID, publicKey string) error
}

type UserService struct {
	UserRepository UserRepository
}

func (i *UserService) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	return i.UserRepository.GetAllUsers(ctx)
}

func (i *UserService) GetUser(ctx context.Context, uuid *uuid.UUID) (*entity.User, error) {
	return i.UserRepository.GetUser(ctx, uuid)
}

func (i *UserService) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	return i.UserRepository.CreateUser(ctx, user)
}

func (i *UserService) UpdateUser(ctx context.Context, uuid *uuid.UUID, user *entity.User) (*entity.User, error) {
	return i.UserRepository.UpdateUser(ctx, uuid, user)
}

func (i *UserService) DeleteUser(ctx context.Context, uuid *uuid.UUID) error {
	return i.UserRepository.DeleteUser(ctx, uuid)
}

func (i *UserService) RegisterMatrixUser(ctx context.Context, username string) (*entity.User, error) {
	return i.UserRepository.RegisterMatrixUser(ctx, username)
}

func (i *UserService) SetPublicKey(ctx context.Context, userUUID *uuid.UUID, publicKey string) error {
	return i.UserRepository.SetPublicKey(ctx, userUUID, publicKey)
}

func (i *UserService) GetMatrixUserByUsername(ctx context.Context, username string) (*entity.MatrixUser, error) {
	return i.UserRepository.GetMatrixUserByUsername(ctx, username)
}
