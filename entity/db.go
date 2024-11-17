package entity

import (
	"context"

	"github.com/gofrs/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const DatabaseUrlKey = "DATABASE_URL"

type Repository interface {
	GetAllMenus(*gorm.DB) ([]Menu, error)
	GetAllUsers(*gorm.DB) ([]User, error)
	GetAllMatrixUsers(*gorm.DB) ([]MatrixUser, error)
	GetAllPasswordUsers(*gorm.DB) ([]PasswordUser, error)
	GetAllOrders(*gorm.DB) ([]Order, error)
	GetOrder(*gorm.DB, *uuid.UUID) (*Order, error)
	GetAllOrderItems(*gorm.DB, *uuid.UUID) ([]OrderItem, error)
	GetOrderItem(*gorm.DB, *uuid.UUID) (*OrderItem, error)
	CreateOrderItem(*gorm.DB, *uuid.UUID, *OrderItem) (*OrderItem, error)
	GetMenu(*gorm.DB, *uuid.UUID) (*Menu, error)
	GetMenuByName(*gorm.DB, string) (*Menu, error)
	GetMenuItem(*gorm.DB, *uuid.UUID) (*MenuItem, error)
	CreateMenu(*gorm.DB, *Menu) (*Menu, error)
	UpdateMenu(*gorm.DB, *uuid.UUID, *Menu) (*Menu, error)
	CreateMenuItem(*gorm.DB, *MenuItem, *uuid.UUID) (*MenuItem, error)
	DeleteMenuItem(*gorm.DB, *uuid.UUID) error
	DeleteMenu(*gorm.DB, *uuid.UUID) error
	CreateOrder(*gorm.DB, *Order) (*Order, error)
	UpdateOrder(*gorm.DB, *uuid.UUID, *Order) (*Order, error)
	UpdateOrderItem(*gorm.DB, *uuid.UUID, *OrderItem) (*OrderItem, error)
	DeleteOrderItem(*gorm.DB, *uuid.UUID) error
	DeleteOrder(*gorm.DB, *uuid.UUID) error
	GetUser(*gorm.DB, *uuid.UUID) (*User, error)
	GetUserByPublicKey(*gorm.DB, string) (*User, error)
	CreateUser(*gorm.DB, *User) (*User, error)
	UpdateUser(*gorm.DB, *uuid.UUID, *User) (*User, error)
	DeleteUser(*gorm.DB, *uuid.UUID) error
	GetMatrixUser(*gorm.DB, *uuid.UUID) (*MatrixUser, error)
	GetMatrixUserByUsername(*gorm.DB, string) (*MatrixUser, error)
	CreateMatrixUser(*gorm.DB, *MatrixUser) (*MatrixUser, error)
	UpdateMatrixUser(*gorm.DB, *uuid.UUID, *MatrixUser) (*MatrixUser, error)
	DeleteMatrixUser(*gorm.DB, *uuid.UUID) error
	FindPasswordUser(*gorm.DB, string) (*PasswordUser, error)
	GetPasswordUser(*gorm.DB, *uuid.UUID) (*PasswordUser, error)
	CreatePasswordUser(*gorm.DB, *PasswordUser) (*PasswordUser, error)
	UpdatePasswordUser(*gorm.DB, *uuid.UUID, *PasswordUser) (*PasswordUser, error)
	DeletePasswordUser(*gorm.DB, *uuid.UUID) error
	Transaction(func(tx *gorm.DB) error) error
}

type RepositoryImpl struct {
	Db  *gorm.DB
	ctx context.Context
}

func NewRepository(ctx context.Context, databaseUrl string) (*RepositoryImpl, error) {
	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
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
