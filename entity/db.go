package entity

import (
	"context"

	"github.com/gofrs/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const DatabaseUrlKey = "DATABASE_URL"

type Repository interface {
	GetAllMenus(tx *gorm.DB) ([]Menu, error)
	GetAllUsers(tx *gorm.DB) ([]User, error)
	GetAllMatrixUsers(tx *gorm.DB) ([]MatrixUser, error)
	GetAllPasswordUsers(tx *gorm.DB) ([]PasswordUser, error)
	GetAllOrders(tx *gorm.DB) ([]Order, error)
	GetOrder(tx *gorm.DB, uuid *uuid.UUID) (*Order, error)
	GetActiveOrderByMenu(tx *gorm.DB, uuid *uuid.UUID) (*Order, error)
	GetActiveOrderByMenuName(tx *gorm.DB, name string) (*Order, error)
	GetAllOrderItems(tx *gorm.DB, uuid *uuid.UUID) ([]OrderItem, error)
	GetOrderItem(tx *gorm.DB, uuid *uuid.UUID) (*OrderItem, error)
	CreateOrderItem(tx *gorm.DB, uuid *uuid.UUID, orderItem *OrderItem) (*OrderItem, error)
	GetMenu(tx *gorm.DB, uuid *uuid.UUID) (*Menu, error)
	GetMenuByName(tx *gorm.DB, name string) (*Menu, error)
	GetMenuItem(tx *gorm.DB, uuid *uuid.UUID) (*MenuItem, error)
	GetMenuItemByShortName(tx *gorm.DB, menuUuid *uuid.UUID, shortName string) (*MenuItem, error)
	CreateMenu(tx *gorm.DB, menu *Menu) (*Menu, error)
	UpdateMenu(tx *gorm.DB, uuid *uuid.UUID, menu *Menu) (*Menu, error)
	CreateMenuItem(tx *gorm.DB, menuItem *MenuItem) (*MenuItem, error)
	DeleteMenuItem(tx *gorm.DB, uuid *uuid.UUID) error
	DeleteMenu(tx *gorm.DB, uuid *uuid.UUID) error
	CreateOrder(tx *gorm.DB, order *Order) (*Order, error)
	UpdateOrder(tx *gorm.DB, uuid *uuid.UUID, currentUser *uuid.UUID, order *Order) (*Order, error)
	UpdateOrderItem(tx *gorm.DB, orderItemUuid *uuid.UUID, userUuid *uuid.UUID, orderItem *OrderItem) (*OrderItem, error)
	DeleteOrderItem(tx *gorm.DB, uuid *uuid.UUID) error
	GetAllOrderItemsForOrderAndUser(tx *gorm.DB, orderUuid *uuid.UUID, userUuid *uuid.UUID) ([]OrderItem, error)
	DeleteOrder(tx *gorm.DB, uuid *uuid.UUID) error
	GetUser(tx *gorm.DB, uuid *uuid.UUID) (*User, error)
	GetUserByName(tx *gorm.DB, name string) (*User, error)
	CreateUser(tx *gorm.DB, user *User) (*User, error)
	UpdateUser(tx *gorm.DB, uuid *uuid.UUID, user *User) (*User, error)
	DeleteUser(tx *gorm.DB, uuid *uuid.UUID) error
	GetMatrixUser(tx *gorm.DB, uuid *uuid.UUID) (*MatrixUser, error)
	GetMatrixUserByUsername(tx *gorm.DB, username string) (*MatrixUser, error)
	CreateMatrixUser(tx *gorm.DB, matrixUser *MatrixUser) (*MatrixUser, error)
	UpdateMatrixUser(tx *gorm.DB, uuid *uuid.UUID, matrixUser *MatrixUser) (*MatrixUser, error)
	DeleteMatrixUser(tx *gorm.DB, uuid *uuid.UUID) error
	FindPasswordUser(tx *gorm.DB, username string) (*PasswordUser, error)
	GetPasswordUser(tx *gorm.DB, uuid *uuid.UUID) (*PasswordUser, error)
	CreatePasswordUser(tx *gorm.DB, passwordUser *PasswordUser) (*PasswordUser, error)
	UpdatePasswordUser(tx *gorm.DB, uuid *uuid.UUID, passwordUser *PasswordUser) (*PasswordUser, error)
	DeletePasswordUser(tx *gorm.DB, uuid *uuid.UUID) error
	GetAllSshUsers(tx *gorm.DB) ([]SshUser, error)
	GetSshUser(tx *gorm.DB, uuid *uuid.UUID) (*SshUser, error)
	GetSshUserByPublicKey(tx *gorm.DB, publicKey string) (*SshUser, error)
	CreateSshUser(tx *gorm.DB, sshUser *SshUser) (*SshUser, error)
	UpdateSshUser(tx *gorm.DB, uuid *uuid.UUID, sshUser *SshUser) (*SshUser, error)
	DeleteSshUser(tx *gorm.DB, uuid *uuid.UUID) error
	Transaction(func(tx *gorm.DB) error) error
}

type RepositoryImpl struct {
	Db  *gorm.DB
	ctx context.Context
}

func NewRepository(ctx context.Context, databaseUrl string) (*RepositoryImpl, error) {
	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	return &RepositoryImpl{
		Db:  db,
		ctx: ctx,
	}, nil
}

func (repo *RepositoryImpl) Transaction(callback func(tx *gorm.DB) error) error {
	return repo.Db.Transaction(callback)
}
