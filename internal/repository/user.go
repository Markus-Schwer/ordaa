package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"

	"github.com/Markus-Schwer/ordaa/internal/entity"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrCannotGetAllUsers = errors.New("could not get all users from db")
	ErrGettingUser       = errors.New("failed to get user")
	ErrUserNotFound      = errors.New("user not found")
	ErrCreatingUser      = errors.New("could not create user")
	ErrUpdatingUser      = errors.New("could not update users")
	ErrDeletingUser      = errors.New("could not delete user")
	ErrSettingPublicKey  = errors.New("setting public key for user")
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) RegisterMatrixUser(ctx context.Context, username string) (*entity.User, error) {
	tx := r.DB.Begin()

	matrixUser, err := r.GetMatrixUserByUsername(ctx, username)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		_ = tx.Rollback()
		return nil, err
	}

	if matrixUser != nil {
		_ = tx.Rollback()
		return nil, ErrUserAlreadyExists
	}

	user := &entity.User{Name: username}
	if err = tx.Create(user).Error; err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrCreatingUser, err)
	}

	if err = tx.Create(&entity.MatrixUser{Username: username, UserUUID: user.UUID}).Error; err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrCreatingUser, err)
	}

	_ = tx.Commit()

	return user, nil
}

func (r *UserRepository) SetPublicKey(ctx context.Context, userUUID *uuid.UUID, publicKey string) error {
	tx := r.DB.Begin()

	user, err := r.GetUser(ctx, userUUID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	sshUser, err := r.GetSSHUser(ctx, user.UUID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	sshUser.PublicKey = publicKey

	if err = tx.Save(sshUser).Error; err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%w: %w", ErrSettingPublicKey, err)
	}

	_ = tx.Commit()

	return nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	users := []entity.User{}

	err := r.DB.Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCannotGetAllUsers, err)
	}

	return users, nil
}

func (r *UserRepository) GetUser(ctx context.Context, userUUID *uuid.UUID) (*entity.User, error) {
	var user entity.User

	err := r.DB.First(&user, userUUID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingUser, err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByName(ctx context.Context, name string) (*entity.User, error) {
	var user entity.User

	err := r.DB.Where(&entity.User{Name: name}).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingUser, err)
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	tx := r.DB.Begin()

	err := tx.Create(&user).Error
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrCreatingUser, err)
	}

	_ = tx.Commit()

	return user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, userUUID *uuid.UUID, user *entity.User) (*entity.User, error) {
	tx := r.DB.Begin()

	foundUser, err := r.GetUser(ctx, userUUID)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrUpdatingUser, err)
	}

	foundUser.Name = user.Name

	err = tx.Save(&foundUser).Error
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrUpdatingUser, err)
	}

	_ = tx.Commit()

	return foundUser, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userUUID *uuid.UUID) error {
	tx := r.DB.Begin()

	err := tx.Delete(&entity.User{UUID: userUUID}).Error
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%w: %w", ErrDeletingUser, err)
	}

	// TODO: check if the user is a password or matrix user
	if err = tx.Where(&entity.MatrixUser{UserUUID: userUUID}).Delete(&entity.MatrixUser{}).Error; err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%w: %w", ErrDeletingUser, err)
	}

	if err = tx.Where(&entity.PasswordUser{UserUUID: userUUID}).Delete(&entity.PasswordUser{}).Error; err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%w: %w", ErrDeletingUser, err)
	}

	_ = tx.Commit()

	return nil
}

func (r *UserRepository) GetAllMatrixUsers(ctx context.Context) ([]entity.MatrixUser, error) {
	matrixUsers := []entity.MatrixUser{}

	err := r.DB.Find(&matrixUsers).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCannotGetAllUsers, err)
	}

	return matrixUsers, nil
}

func (r *UserRepository) GetMatrixUser(ctx context.Context, matrixUserUUID *uuid.UUID) (*entity.MatrixUser, error) {
	var matrixUser entity.MatrixUser

	err := r.DB.Where(&entity.MatrixUser{UUID: matrixUserUUID}).First(&matrixUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingUser, err)
	}

	return &matrixUser, nil
}

func (r *UserRepository) GetMatrixUserByUsername(ctx context.Context, username string) (*entity.MatrixUser, error) {
	var matrixUser entity.MatrixUser

	err := r.DB.Where(&entity.MatrixUser{Username: username}).First(&matrixUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingUser, err)
	}

	return &matrixUser, nil
}

func (r *UserRepository) CreateMatrixUser(ctx context.Context, matrixUser *entity.MatrixUser) (*entity.MatrixUser, error) {
	tx := r.DB.Begin()

	err := tx.Create(&matrixUser).Error
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrCreatingUser, err)
	}

	_ = tx.Commit()

	return matrixUser, nil
}

func (r *UserRepository) UpdateMatrixUser(
	ctx context.Context,
	matrixUserUUID *uuid.UUID,
	matrixUser *entity.MatrixUser,
) (*entity.MatrixUser, error) {
	tx := r.DB.Begin()

	existingMatrixUser, err := r.GetMatrixUser(ctx, matrixUserUUID)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrUpdatingUser, err)
	}

	existingMatrixUser.Username = matrixUser.Username
	existingMatrixUser.UserUUID = matrixUser.UserUUID

	err = tx.Save(&existingMatrixUser).Error
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrUpdatingUser, err)
	}

	_ = tx.Commit()

	return existingMatrixUser, nil
}

func (r *UserRepository) DeleteMatrixUser(ctx context.Context, userUUID *uuid.UUID) error {
	tx := r.DB.Begin()

	err := tx.Where(&entity.MatrixUser{UserUUID: userUUID}).Delete(&entity.MatrixUser{}).Error
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%w: %w", ErrDeletingUser, err)
	}

	_ = tx.Commit()

	return nil
}

func (r *UserRepository) GetAllPasswordUsers(ctx context.Context) ([]entity.PasswordUser, error) {
	passwordUsers := []entity.PasswordUser{}

	err := r.DB.Find(&passwordUsers).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCannotGetAllUsers, err)
	}

	return passwordUsers, nil
}

func (r *UserRepository) FindPasswordUser(ctx context.Context, username string) (*entity.PasswordUser, error) {
	var passwordUser entity.PasswordUser

	err := r.DB.Where(&entity.PasswordUser{Username: username}).Find(&passwordUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingUser, err)
	}

	return &passwordUser, nil
}

func (r *UserRepository) GetPasswordUser(ctx context.Context, passwordUserUUID *uuid.UUID) (*entity.PasswordUser, error) {
	var passwordUser entity.PasswordUser

	err := r.DB.Where(&entity.PasswordUser{UUID: passwordUserUUID}).First(&passwordUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingUser, err)
	}

	return &passwordUser, nil
}

func (r *UserRepository) CreatePasswordUser(ctx context.Context, passwordUser *entity.PasswordUser) (*entity.PasswordUser, error) {
	tx := r.DB.Begin()

	err := tx.Create(&passwordUser).Error
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrCreatingUser, err)
	}

	_ = tx.Commit()

	return passwordUser, nil
}

func (r *UserRepository) UpdatePasswordUser(
	ctx context.Context,
	passwordUserUUID *uuid.UUID,
	passwordUser *entity.PasswordUser,
) (*entity.PasswordUser, error) {
	tx := r.DB.Begin()

	existingPasswordUser, err := r.GetPasswordUser(ctx, passwordUserUUID)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrUpdatingUser, err)
	}

	existingPasswordUser.Username = passwordUser.Username
	existingPasswordUser.UserUUID = passwordUser.UserUUID
	existingPasswordUser.Password = passwordUser.Password

	err = tx.Save(existingPasswordUser).Error
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrUpdatingUser, err)
	}

	_ = tx.Commit()

	return existingPasswordUser, nil
}

func (r *UserRepository) DeletePasswordUser(ctx context.Context, userUUID *uuid.UUID) error {
	tx := r.DB.Begin()

	if err := tx.Where(&entity.PasswordUser{UserUUID: userUUID}).Delete(&entity.PasswordUser{}).Error; err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%w: %w", ErrDeletingUser, err)
	}

	_ = tx.Commit()

	return nil
}

func (r *UserRepository) GetAllSSHUsers(ctx context.Context) ([]entity.SSHUser, error) {
	sshUsers := []entity.SSHUser{}

	err := r.DB.Find(&sshUsers).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCannotGetAllUsers, err)
	}

	return sshUsers, nil
}

func (r *UserRepository) GetSSHUser(ctx context.Context, sshUserUUID *uuid.UUID) (*entity.SSHUser, error) {
	var sshUser entity.SSHUser

	err := r.DB.Where(&entity.SSHUser{UUID: sshUserUUID}).First(&sshUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingUser, err)
	}

	return &sshUser, nil
}

func (r *UserRepository) GetSSHUserByPublicKey(ctx context.Context, publicKey string) (*entity.SSHUser, error) {
	var sshUser entity.SSHUser

	err := r.DB.Where(&entity.SSHUser{PublicKey: publicKey}).First(&sshUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
	} else if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGettingUser, err)
	}

	return &sshUser, nil
}

func (r *UserRepository) CreateSSHUser(ctx context.Context, sshUser *entity.SSHUser) (*entity.SSHUser, error) {
	tx := r.DB.Begin()

	if err := tx.Create(&sshUser).Error; err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrCreatingUser, err)
	}

	_ = tx.Commit()

	return sshUser, nil
}

func (r *UserRepository) UpdateSSHUser(ctx context.Context, sshUserUUID *uuid.UUID, sshUser *entity.SSHUser) (*entity.SSHUser, error) {
	tx := r.DB.Begin()

	existingSSHUser, err := r.GetSSHUser(ctx, sshUserUUID)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrUpdatingUser, err)
	}

	existingSSHUser.PublicKey = sshUser.PublicKey
	existingSSHUser.UserUUID = sshUser.UserUUID

	err = tx.Save(&existingSSHUser).Error
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("%w: %w", ErrUpdatingUser, err)
	}

	_ = tx.Commit()

	return existingSSHUser, nil
}

func (r *UserRepository) DeleteSSHUser(ctx context.Context, userUUID *uuid.UUID) error {
	tx := r.DB.Begin()

	err := tx.Where(&entity.SSHUser{UserUUID: userUUID}).Delete(&entity.SSHUser{}).Error
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("%w: %w", ErrDeletingUser, err)
	}

	_ = tx.Commit()

	return nil
}
