package entity

import (
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type MockRepository struct {
	users []User
	passwordUsers []PasswordUser
	matrixUsers []MatrixUser
	sshUsers []SshUser
	menus []Menu
	menuItems []MenuItem
	orders []Order
	orderItems []OrderItem
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		users: []User{},
		passwordUsers: []PasswordUser{},
		matrixUsers: []MatrixUser{},
		sshUsers: []SshUser{},
		menus: []Menu{},
		menuItems: []MenuItem{},
		orders: []Order{},
		orderItems: []OrderItem{},
	}
}

func (*MockRepository) Transaction(fn func(*gorm.DB) error) error {
	return fn(nil)
}

func (repo *MockRepository) GetAllMenus(tx *gorm.DB) ([]Menu, error) {
	return repo.menus, nil
}

func (repo *MockRepository) GetMenu(tx *gorm.DB, menuUuid *uuid.UUID) (*Menu, error) {
	for _, m := range repo.menus {
		if *m.Uuid == *menuUuid {
			return &m, nil
		}
	}

	return nil, ErrMenuNotFound
}

func (repo *MockRepository) GetMenuByName(tx *gorm.DB, name string) (*Menu, error) {
	for _, m := range repo.menus {
		if m.Name == name {
			return &m, nil
		}
	}

	return nil, ErrMenuNotFound
}

func (repo *MockRepository) GetMenuItem(tx *gorm.DB, menuItemUuid *uuid.UUID) (*MenuItem, error) {
	for _, mi := range repo.menuItems {
		if *mi.Uuid == *menuItemUuid {
			return &mi, nil
		}
	}

	return nil, ErrMenuItemNotFound
}

func (repo *MockRepository) GetMenuItemByShortName(tx *gorm.DB, menuUuid *uuid.UUID, shortName string) (*MenuItem, error) {
	for _, mi := range repo.menuItems {
		if *mi.MenuUuid == *menuUuid && mi.ShortName == shortName {
			return &mi, nil
		}
	}

	return nil, ErrMenuItemNotFound
}

func (repo *MockRepository) CreateMenu(tx *gorm.DB, menu *Menu) (*Menu, error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	menu.Uuid = &newUuid

	repo.menus = append(repo.menus, *menu)
	return menu, nil
}

func (repo *MockRepository) UpdateMenu(tx *gorm.DB, menuUuid *uuid.UUID, menu *Menu) (*Menu, error) {
	if err := repo.DeleteMenu(tx, menuUuid); err != nil {
		return nil, err
	}
	if _, err := repo.CreateMenu(tx, menu); err != nil {
		return nil, err
	}
	return menu, nil
}

func (repo *MockRepository) CreateMenuItem(tx *gorm.DB, menuItem *MenuItem) (*MenuItem, error) {
	*menuItem.Uuid, _ = uuid.NewV4()
	repo.menuItems = append(repo.menuItems, *menuItem)

	foundIndex := -1
	for i, m := range repo.menus {
		if *m.Uuid == *menuItem.MenuUuid {
			foundIndex = i
		}
	}

	if foundIndex == -1 {
		return nil, errors.New("menu not found")
	}

	repo.menus[foundIndex].Items = append(repo.menus[foundIndex].Items, *menuItem)
	return menuItem, nil
}

func (repo *MockRepository) DeleteMenuItem(tx *gorm.DB, menuItemUuid *uuid.UUID) error {
	var menuUuid *uuid.UUID
	newMenuItems := []MenuItem{}
	for _, mi := range repo.menuItems {
		if *mi.Uuid != *menuItemUuid {
			newMenuItems = append(newMenuItems, mi)
		} else {
			menuUuid = mi.MenuUuid
		}
	}

	for i, m := range repo.menus {
		if *m.Uuid != *menuUuid {
			repo.menus[i].Items = newMenuItems
		}
	}

	repo.menuItems = newMenuItems
	return nil
}

func (repo *MockRepository) DeleteMenu(tx *gorm.DB, menuUuid *uuid.UUID) error {
	newMenus := []Menu{}
	for _, m := range repo.menus {
		if *m.Uuid != *menuUuid {
			newMenus = append(newMenus, m)
		} else {
		}
	}

	newMenuItems := []MenuItem{}
	for _, mi := range repo.menuItems {
		if *mi.Uuid != *menuUuid {
			newMenuItems = append(newMenuItems, mi)
		} else {
		}
	}

	repo.menus = newMenus
	repo.menuItems = newMenuItems
	return nil
}

func (repo *MockRepository) GetAllOrders(tx *gorm.DB) ([]Order, error) {
	return repo.orders, nil
}

func (repo *MockRepository) GetOrder(tx *gorm.DB, uuid *uuid.UUID) (*Order, error) {
	for _, o := range repo.orders {
		if *o.Uuid == *uuid {
			return &o, nil
		}
	}

	return nil, ErrOrderNotFound
}

func (repo *MockRepository) GetActiveOrderByMenu(tx *gorm.DB, menuUuid *uuid.UUID) (*Order, error) {
	for _, o := range repo.orders {
		if *o.MenuUuid == *menuUuid {
			return &o, nil
		}
	}

	return nil, ErrOrderNotFound
}

func (repo *MockRepository) GetActiveOrderByMenuName(tx *gorm.DB, menuName string) (*Order, error) {
	menu, err := repo.GetMenuByName(tx, menuName)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	for _, o := range repo.orders {
		if *o.MenuUuid == *menu.Uuid {
			return &o, nil
		}
	}

	return nil, ErrOrderNotFound
}

func (repo *MockRepository) GetAllOrderItems(tx *gorm.DB, orderUuid *uuid.UUID) ([]OrderItem, error) {
	return repo.orderItems, nil
}

func (repo *MockRepository) GetOrderItem(tx *gorm.DB, uuid *uuid.UUID) (*OrderItem, error) {
	for _, oi := range repo.orderItems {
		if *oi.Uuid == *uuid {
			return &oi, nil
		}
	}

	return nil, ErrOrderItemNotFound
}

func (repo *MockRepository) CreateOrderItem(tx *gorm.DB, orderUuid *uuid.UUID, orderItem *OrderItem) (*OrderItem, error) {
	menuItemUuid := orderItem.MenuItemUuid
	if menuItemUuid == nil {
		return nil, fmt.Errorf("%w: %w", ErrCreatingOrderItem, ErrMenuItemUuidMissing)
	}

	menuItem, err := repo.GetMenuItem(tx, menuItemUuid)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreatingOrderItem, err)
	}

	orderItem.Paid = false
	orderItem.Price = menuItem.Price

	newUuid, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreatingOrderItem, err)
	}
	orderItem.Uuid = &newUuid

	repo.orderItems = append(repo.orderItems, *orderItem)

	foundIndex := -1
	for i, o := range repo.orders {
		if *o.Uuid == *orderUuid {
			foundIndex = i
		}
	}

	if foundIndex == -1 {
		return nil, errors.New("order not found")
	}

	repo.orders[foundIndex].Items = append(repo.orders[foundIndex].Items, *orderItem)
	return orderItem, nil
}

func (repo *MockRepository) CreateOrder(tx *gorm.DB, order *Order) (*Order, error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	order.Uuid = &newUuid

	repo.orders = append(repo.orders, *order)
	return order, nil
}

func (repo *MockRepository) UpdateOrder(tx *gorm.DB, orderUuid *uuid.UUID, order *Order) (*Order, error) {
	if err := repo.DeleteOrder(tx, orderUuid); err != nil {
		return nil, err
	}
	if _, err := repo.CreateOrder(tx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (repo *MockRepository) UpdateOrderItem(tx *gorm.DB, orderItemUuid *uuid.UUID, orderItem *OrderItem) (*OrderItem, error) {
	if err := repo.DeleteOrderItem(tx, orderItemUuid); err != nil {
		return nil, err
	}
	if _, err := repo.CreateOrderItem(tx, orderItem.OrderUuid, orderItem); err != nil {
		return nil, err
	}
	return orderItem, nil
}

func (repo *MockRepository) DeleteOrderItem(tx *gorm.DB, orderItemUuid *uuid.UUID) error {
	var orderUuid *uuid.UUID
	newOrderItems := []OrderItem{}
	for _, oi := range repo.orderItems {
		if oi.Uuid != orderItemUuid {
			newOrderItems = append(newOrderItems, oi)
		} else {
			orderUuid = oi.OrderUuid
		}
	}

	for i, o := range repo.orders {
		if o.Uuid != orderUuid {
			repo.orders[i].Items = newOrderItems
		}
	}

	repo.orderItems = newOrderItems
	return nil
}

func (repo *MockRepository) DeleteOrder(tx *gorm.DB, orderUuid *uuid.UUID) error {
	newOrders := []Order{}
	for _, o := range repo.orders {
		if o.Uuid != orderUuid {
			newOrders = append(newOrders, o)
		} else {
		}
	}

	newOrderItems := []OrderItem{}
	for _, oi := range repo.orderItems {
		if oi.Uuid != orderUuid {
			newOrderItems = append(newOrderItems, oi)
		} else {
		}
	}

	repo.orders = newOrders
	repo.orderItems = newOrderItems
	return nil
}

func (repo *MockRepository) GetAllUsers(tx *gorm.DB) ([]User, error) {
	return repo.users, nil
}

func (repo *MockRepository) GetUser(tx *gorm.DB, userUuid *uuid.UUID) (*User, error) {
	for _, u := range repo.users {
		if *u.Uuid == *userUuid {
			return &u, nil
		}
	}

	return nil, ErrUserNotFound
}

func (repo *MockRepository) CreateUser(tx *gorm.DB, user *User) (*User, error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	user.Uuid = &newUuid

	repo.users = append(repo.users, *user)
	return user, nil
}

func (repo *MockRepository) UpdateUser(tx *gorm.DB, userUuid *uuid.UUID, user *User) (*User, error) {
	if err := repo.DeleteUser(tx, userUuid); err != nil {
		return nil, err
	}
	if _, err := repo.CreateUser(tx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *MockRepository) DeleteUser(tx *gorm.DB, userUuid *uuid.UUID) error {
	newUsers := []User{}
	for _, u := range repo.users {
		if u.Uuid != userUuid {
			newUsers = append(newUsers, u)
		} else {
		}
	}

	repo.users = newUsers
	return nil
}

func (repo *MockRepository) GetAllMatrixUsers(tx *gorm.DB) ([]MatrixUser, error) {
	return repo.matrixUsers, nil
}

func (repo *MockRepository) GetMatrixUser(tx *gorm.DB, matrixUserUuid *uuid.UUID) (*MatrixUser, error) {
	for _, u := range repo.matrixUsers {
		if *u.Uuid == *matrixUserUuid {
			return &u, nil
		}
	}

	return nil, ErrUserNotFound
}

func (repo *MockRepository) GetMatrixUserByUsername(tx *gorm.DB, username string) (*MatrixUser, error) {
	for _, u := range repo.matrixUsers {
		if u.Username == username {
			return &u, nil
		}
	}

	return nil, ErrUserNotFound
}

func (repo *MockRepository) CreateMatrixUser(tx *gorm.DB, matrixUser *MatrixUser) (*MatrixUser, error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	matrixUser.Uuid = &newUuid

	repo.matrixUsers = append(repo.matrixUsers, *matrixUser)
	return matrixUser, nil
}

func (repo *MockRepository) UpdateMatrixUser(tx *gorm.DB, matrixUserUuid *uuid.UUID, matrixUser *MatrixUser) (*MatrixUser, error) {
	if err := repo.DeleteMatrixUser(tx, matrixUserUuid); err != nil {
		return nil, err
	}
	if _, err := repo.CreateMatrixUser(tx, matrixUser); err != nil {
		return nil, err
	}
	return matrixUser, nil
}

func (repo *MockRepository) DeleteMatrixUser(tx *gorm.DB, userUuid *uuid.UUID) error {
	newMatrixUsers := []MatrixUser{}
	for _, u := range repo.matrixUsers {
		if u.UserUuid != userUuid {
			newMatrixUsers = append(newMatrixUsers, u)
		}
	}

	repo.matrixUsers = newMatrixUsers
	return nil
}

func (repo *MockRepository) GetAllPasswordUsers(tx *gorm.DB) ([]PasswordUser, error) {
	return repo.passwordUsers, nil
}

func (repo *MockRepository) FindPasswordUser(tx *gorm.DB, username string) (*PasswordUser, error) {
	for _, u := range repo.passwordUsers {
		if u.Username == username {
			return &u, nil
		}
	}

	return nil, errors.New("password user not found")
}

func (repo *MockRepository) GetPasswordUser(tx *gorm.DB, passwordUserUuid *uuid.UUID) (*PasswordUser, error) {
	for _, u := range repo.passwordUsers {
		if *u.Uuid == *passwordUserUuid {
			return &u, nil
		}
	}

	return nil, ErrUserNotFound
}

func (repo *MockRepository) CreatePasswordUser(tx *gorm.DB, passwordUser *PasswordUser) (*PasswordUser, error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	passwordUser.Uuid = &newUuid

	repo.passwordUsers = append(repo.passwordUsers, *passwordUser)
	return passwordUser, nil
}

func (repo *MockRepository) UpdatePasswordUser(tx *gorm.DB, passwordUserUuid *uuid.UUID, passwordUser *PasswordUser) (*PasswordUser, error) {
	if err := repo.DeletePasswordUser(tx, passwordUserUuid); err != nil {
		return nil, err
	}
	if _, err := repo.CreatePasswordUser(tx, passwordUser); err != nil {
		return nil, err
	}
	return passwordUser, nil
}

func (repo *MockRepository) DeletePasswordUser(tx *gorm.DB, userUuid *uuid.UUID) error {
	newPasswordUsers := []PasswordUser{}
	for _, u := range repo.passwordUsers {
		if u.UserUuid != userUuid {
			newPasswordUsers = append(newPasswordUsers, u)
		}
	}

	repo.passwordUsers = newPasswordUsers
	return nil
}

func (repo *MockRepository) GetAllSshUsers(tx *gorm.DB) ([]SshUser, error) {
	return repo.sshUsers, nil
}

func (repo *MockRepository) GetSshUser(tx *gorm.DB, sshUserUuid *uuid.UUID) (*SshUser, error) {
	for _, u := range repo.sshUsers {
		if *u.Uuid == *sshUserUuid {
			return &u, nil
		}
	}

	return nil, ErrUserNotFound
}

func (repo *MockRepository) GetSshUserByPublicKey(tx *gorm.DB, publicKey string) (*SshUser, error) {
	for _, u := range repo.sshUsers {
		if u.PublicKey == publicKey {
			return &u, nil
		}
	}

	return nil, ErrUserNotFound
}

func (repo *MockRepository) CreateSshUser(tx *gorm.DB, sshUser *SshUser) (*SshUser, error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	sshUser.Uuid = &newUuid

	repo.sshUsers = append(repo.sshUsers, *sshUser)
	return sshUser, nil
}

func (repo *MockRepository) UpdateSshUser(tx *gorm.DB, sshUserUuid *uuid.UUID, sshUser *SshUser) (*SshUser, error) {
	if err := repo.DeleteSshUser(tx, sshUserUuid); err != nil {
		return nil, err
	}
	if _, err := repo.CreateSshUser(tx, sshUser); err != nil {
		return nil, err
	}
	return sshUser, nil
}

func (repo *MockRepository) DeleteSshUser(tx *gorm.DB, userUuid *uuid.UUID) error {
	newSshUsers := []SshUser{}
	for _, u := range repo.sshUsers {
		if u.UserUuid != userUuid {
			newSshUsers = append(newSshUsers, u)
		}
	}

	repo.sshUsers = newSshUsers
	return nil
}

