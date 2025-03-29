package matrix

import (
	"context"
	"errors"
	"strings"
	"testing"

	"gitlab.com/sfz.aalen/hackwerk/ordaa/entity"
	"gitlab.com/sfz.aalen/hackwerk/ordaa/utils/ptr"
	"gorm.io/gorm"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHelpHandler(t *testing.T) {
	ctx := context.Background()
	repo := &entity.RepositoryMock{}
	m := &MatrixBoundaryMock{}

	message := "help"
	evt := &event.Event{}

	m.replyFunc = func(room id.RoomID, evt id.EventID, content string, asHtml bool) id.EventID {
		return evt
	}

	handler := handlers[message]

	assert.NoError(t, handler(ctx, m, repo, nil, evt, message))
	calls := m.replyCalls()
	assert.Equal(t, 1, len(calls))
	assert.Equal(t, "Hello world", calls[0].Content)
}

func TestRegisterHandler(t *testing.T) {
	type testCase struct {
		user           *entity.User
		matrixUser     *entity.MatrixUser
		matrixUsername string
		expectedErr    error
		expectedReply  *string
	}

	testCases := []testCase{
		{
			user:           ptr.To(entity.User{Uuid: ptr.To(uuid.Must(uuid.FromString("bdfc1ef0-2388-4e39-a482-a85914e04140"))), Name: "@test:matrix.org"}),
			matrixUser:     ptr.To(entity.MatrixUser{UserUuid: ptr.To(uuid.Must(uuid.FromString("bdfc1ef0-2388-4e39-a482-a85914e04140"))), Username: "@test:matrix.org"}),
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("user '@test:matrix.org' is already registered"),
			expectedReply:  nil,
		},
		{
			user:           nil,
			matrixUser:     nil,
			matrixUsername: "@test:matrix.org",
			expectedErr:    nil,
			expectedReply:  ptr.To("successfully registered user: @test:matrix.org"),
		},
	}

	for _, c := range testCases {
		ctx := context.Background()
		repo := &entity.RepositoryMock{}
		m := &MatrixBoundaryMock{}

		repo.GetMatrixUserByUsernameFunc = func(tx *gorm.DB, username string) (*entity.MatrixUser, error) {
			if c.matrixUser != nil && username == c.matrixUser.Username {
				return c.matrixUser, nil
			}
			return nil, entity.ErrUserNotFound
		}

		repo.CreateUserFunc = func(tx *gorm.DB, userParam *entity.User) (*entity.User, error) {
			return userParam, nil
		}

		repo.CreateMatrixUserFunc = func(tx *gorm.DB, matrixUserParam *entity.MatrixUser) (*entity.MatrixUser, error) {
			return matrixUserParam, nil
		}

		m.replyFunc = func(room id.RoomID, evt id.EventID, content string, asHtml bool) id.EventID {
			return evt
		}

		message := "register"
		evt := &event.Event{Sender: id.UserID(c.matrixUsername)}

		handler := handlers[message]

		err := handler(ctx, m, repo, nil, evt, message)
		if c.expectedErr != nil {
			assert.Equal(t, c.expectedErr, err)
		} else {
			assert.NoError(t, err)
		}

		if c.expectedReply != nil {
			calls := m.replyCalls()
			assert.Equal(t, 1, len(calls))
			assert.Equal(t, *c.expectedReply, calls[0].Content)
		}
	}
}

func TestSetPublicKeyHandler(t *testing.T) {
	type testCase struct {
		message        string
		user           *entity.User
		sshUser        *entity.SshUser
		matrixUsername string
		expectedErr    error
		expectedReply  *string
	}

	testCases := []testCase{
		{
			message:        "set_public_key quatsch",
			user:           ptr.To(entity.User{Uuid: ptr.To(uuid.Must(uuid.FromString("bdfc1ef0-2388-4e39-a482-a85914e04140"))), Name: "@test:matrix.org"}),
			sshUser:        ptr.To(entity.SshUser{UserUuid: ptr.To(uuid.Must(uuid.FromString("bdfc1ef0-2388-4e39-a482-a85914e04140"))), PublicKey: ""}),
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("public key must be in the format 'key_type base64_encoded_key'"),
			expectedReply:  nil,
		},
		{
			message:        "set_public_key quatsch quatsch",
			user:           ptr.To(entity.User{Uuid: ptr.To(uuid.Must(uuid.FromString("bdfc1ef0-2388-4e39-a482-a85914e04140"))), Name: "@test:matrix.org"}),
			sshUser:        ptr.To(entity.SshUser{UserUuid: ptr.To(uuid.Must(uuid.FromString("bdfc1ef0-2388-4e39-a482-a85914e04140"))), PublicKey: ""}),
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("cannot parse public key"),
			expectedReply:  nil,
		},
		{
			message:        "set_public_key ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMIBClfyMvCvvRbGZNcdOQce7V6dbqUhbLTzLgE0n8sw test",
			user:           ptr.To(entity.User{Uuid: ptr.To(uuid.Must(uuid.FromString("bdfc1ef0-2388-4e39-a482-a85914e04140"))), Name: "@test:matrix.org"}),
			sshUser:        ptr.To(entity.SshUser{UserUuid: ptr.To(uuid.Must(uuid.FromString("bdfc1ef0-2388-4e39-a482-a85914e04140"))), PublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMIBClfyMvCvvRbGZNcdOQce7V6dbqUhbLTzLgE0n8sw test"}),
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("public key must be in the format 'key_type base64_encoded_key'"),
			expectedReply:  nil,
		},
		{
			message:        "set_public_key ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMIBClfyMvCvvRbGZNcdOQce7V6dbqUhbLTzLgE0n8sw",
			user:           ptr.To(entity.User{Uuid: ptr.To(uuid.Must(uuid.FromString("bdfc1ef0-2388-4e39-a482-a85914e04140"))), Name: "@test:matrix.org"}),
			sshUser:        ptr.To(entity.SshUser{UserUuid: ptr.To(uuid.Must(uuid.FromString("bdfc1ef0-2388-4e39-a482-a85914e04140"))), PublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMIBClfyMvCvvRbGZNcdOQce7V6dbqUhbLTzLgE0n8sw"}),
			matrixUsername: "@test:matrix.org",
			expectedErr:    nil,
			expectedReply:  ptr.To("successfully set ssh public key for user: @test:matrix.org"),
		},
		{
			message:        "set_public_key ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMIBClfyMvCvvRbGZNcdOQce7V6dbqUhbLTzLgE0n8sw",
			user:           ptr.To(entity.User{Uuid: ptr.To(uuid.Must(uuid.FromString("bdfc1ef0-2388-4e39-a482-a85914e04140"))), Name: "@test:matrix.org"}),
			sshUser:        nil,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("could not get ssh user for username '@test:matrix.org'"),
			expectedReply:  nil,
		},
	}

	for _, c := range testCases {
		ctx := context.Background()
		repo := &entity.RepositoryMock{}
		m := &MatrixBoundaryMock{}

		m.getUserByUsernameFunc = func(tx *gorm.DB, username string) (*entity.User, error) {
			if username == c.user.Name {
				return c.user, nil
			}
			return nil, entity.ErrUserNotFound
		}

		repo.GetSshUserFunc = func(tx *gorm.DB, uuidMoqParam *uuid.UUID) (*entity.SshUser, error) {
			if c.sshUser != nil && *uuidMoqParam == *c.sshUser.UserUuid {
				return c.sshUser, nil
			}
			return nil, entity.ErrUserNotFound
		}

		repo.UpdateSshUserFunc = func(tx *gorm.DB, uuidMoqParam *uuid.UUID, sshUser *entity.SshUser) (*entity.SshUser, error) {
			if c.sshUser != nil && *uuidMoqParam == *c.sshUser.UserUuid {
				return c.sshUser, nil
			}
			return nil, entity.ErrUserNotFound
		}

		m.replyFunc = func(room id.RoomID, evt id.EventID, content string, asHtml bool) id.EventID {
			return evt
		}

		evt := &event.Event{Sender: id.UserID(c.matrixUsername)}

		handler := handlers["set_public_key"]

		err := handler(ctx, m, repo, nil, evt, c.message)
		if c.expectedErr != nil {
			assert.Equal(t, c.expectedErr, err)
		} else {
			assert.NoError(t, err)
		}

		if c.expectedReply != nil {
			calls := m.replyCalls()
			assert.Equal(t, 1, len(calls))
			assert.Equal(t, *c.expectedReply, calls[0].Content)
		}
	}
}

func TestNewOrderHandler(t *testing.T) {
	user := &entity.User{Uuid: ptr.To(uuid.Must(uuid.FromString("bdfc1ef0-2388-4e39-a482-a85914e04140"))), Name: "@test:matrix.org"}

	menus := map[string]*entity.Menu{
		"sangam": {
			Uuid: ptr.To(uuid.Must(uuid.FromString("1fe21e51-82fc-4372-bba3-7c66783884dc"))),
			Name: "sangam",
		},
		"sangam2": {
			Uuid: ptr.To(uuid.Must(uuid.FromString("555c90b1-95d7-4698-aaa5-db3ba0884f1b"))),
			Name: "sangam2",
		},
	}

	orders := []entity.Order{
		{
			MenuUuid: menus["sangam"].Uuid,
			State: entity.Open,
		},
		{
			MenuUuid: menus["sangam2"].Uuid,
			State: entity.Delivered,
		},
	}

	type testCase struct {
		message        string
		user           *entity.User
		matrixUsername string
		expectedErr    error
		expectedReply  *string
	}

	testCases := []testCase{
		{
			message:        "start",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("menu must be specified"),
			expectedReply:  nil,
		},
		{
			message:        "start ",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("menu must be specified"),
			expectedReply:  nil,
		},
		{
			message:        "start sangam",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("there is already an active order for menu 'sangam'"),
			expectedReply:  nil,
		},
		{
			message:        "start sangam2",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    nil,
			expectedReply:  ptr.To("started new order for sangam2"),
		},
	}

	for _, c := range testCases {
		ctx := context.Background()
		repo := &entity.RepositoryMock{}
		m := &MatrixBoundaryMock{}

		m.getUserByUsernameFunc = func(tx *gorm.DB, username string) (*entity.User, error) {
			if username == c.user.Name {
				return c.user, nil
			}
			return nil, entity.ErrUserNotFound
		}

		repo.GetMenuByNameFunc = func(tx *gorm.DB, name string) (*entity.Menu, error) {
			if menu := menus[name]; menu != nil {
				return menu, nil
			}
			return nil, entity.ErrMenuNotFound
		}

		repo.CreateOrderFunc = func(tx *gorm.DB, orderMoqParam *entity.Order) (*entity.Order, error) {
			for _, order := range orders {
				if *order.MenuUuid == *orderMoqParam.MenuUuid && order.State != entity.Delivered {
					return nil, entity.ErrActiveOrderForMenuAlreadyExists
				}
			}
			return orderMoqParam, nil
		}

		m.replyFunc = func(room id.RoomID, evt id.EventID, content string, asHtml bool) id.EventID {
			return evt
		}

		evt := &event.Event{Sender: id.UserID(c.matrixUsername)}

		handler := handlers["start"]

		err := handler(ctx, m, repo, nil, evt, c.message)
		if c.expectedErr != nil {
			assert.Equal(t, c.expectedErr, err)
		} else {
			assert.NoError(t, err)
		}

		if c.expectedReply != nil {
			calls := m.replyCalls()
			assert.Equal(t, 1, len(calls))
			assert.Equal(t, *c.expectedReply, calls[0].Content)
		}
	}
}

func TestNewOrderItemHandler(t *testing.T) {
	user := &entity.User{Uuid: ptr.To(uuid.Must(uuid.FromString("bdfc1ef0-2388-4e39-a482-a85914e04140"))), Name: "@test:matrix.org"}

	menus := map[string]*entity.Menu{
		"sangam": {
			Uuid: ptr.To(uuid.Must(uuid.FromString("1fe21e51-82fc-4372-bba3-7c66783884dc"))),
			Name: "sangam",
			Items: []entity.MenuItem{
				{
					ShortName: "62",
					Name: "Butter Chicken",
					MenuUuid: ptr.To(uuid.Must(uuid.FromString("555c90b1-95d7-4698-aaa5-db3ba0884f1b"))),
				},
			},
		},
		"sangam2": {
			Uuid: ptr.To(uuid.Must(uuid.FromString("555c90b1-95d7-4698-aaa5-db3ba0884f1b"))),
			Name: "sangam2",
		},
	}

	orders := []entity.Order{
		{
			MenuUuid: menus["sangam"].Uuid,
			State: entity.Open,
		},
		{
			MenuUuid: menus["sangam2"].Uuid,
			State: entity.Delivered,
		},
	}

	type testCase struct {
		message        string
		user           *entity.User
		matrixUsername string
		expectedErr    error
		expectedReply  *string
	}

	testCases := []testCase{
		{
			message:        "add",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("message must be in the format 'add [menu_name] [short_name]'"),
			expectedReply:  nil,
		},
		{
			message:        "add ",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("message must be in the format 'add [menu_name] [short_name]'"),
			expectedReply:  nil,
		},
		{
			message:        "add sangam",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("message must be in the format 'add [menu_name] [short_name]'"),
			expectedReply:  nil,
		},
		{
			message:        "add sangam2 62",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("there is no active order for menu 'sangam2'"),
			expectedReply:  nil,
		},
		{
			message:        "add sangam 62",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    nil,
			expectedReply:  ptr.To("added 'Butter Chicken' to order 'sangam'"),
		},
		{
			message:        "add sangam 63",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("could not get menu item '63'"),
			expectedReply:  nil,
		},
	}

	for _, c := range testCases {
		ctx := context.Background()
		repo := &entity.RepositoryMock{}
		m := &MatrixBoundaryMock{}

		m.getUserByUsernameFunc = func(tx *gorm.DB, username string) (*entity.User, error) {
			if username == c.user.Name {
				return c.user, nil
			}
			return nil, entity.ErrUserNotFound
		}

		repo.GetMenuByNameFunc = func(tx *gorm.DB, name string) (*entity.Menu, error) {
			if menu := menus[name]; menu != nil {
				return menu, nil
			}
			return nil, entity.ErrMenuNotFound
		}

		repo.GetActiveOrderByMenuNameFunc = func(tx *gorm.DB, name string) (*entity.Order, error) {
			var foundMenuUuid *uuid.UUID
			for _, menu := range menus {
				if menu.Name == name {
					foundMenuUuid = menu.Uuid
				}
			}

			if foundMenuUuid == nil {
				return nil, entity.ErrMenuNotFound
			}

			for _, order := range orders {
				if *order.MenuUuid == *foundMenuUuid && order.State != entity.Delivered {
					return &order, nil
				}
			}
			return nil, entity.ErrOrderNotFound
		}

		repo.GetMenuItemByShortNameFunc = func(tx *gorm.DB, menuUuid *uuid.UUID, shortName string) (*entity.MenuItem, error) {
			var foundMenu *entity.Menu
			for _, menu := range menus {
				if *menu.Uuid == *menuUuid {
					foundMenu = menu
				}
			}

			if foundMenu == nil {
				return nil, entity.ErrMenuNotFound
			}

			for _, menuItem := range foundMenu.Items {
				if menuItem.ShortName == shortName {
					return &menuItem, nil
				}
			}

			return nil, entity.ErrMenuItemNotFound
		}

		repo.CreateOrderItemFunc = func(tx *gorm.DB, uuidMoqParam *uuid.UUID, orderItem *entity.OrderItem) (*entity.OrderItem, error) {
			return orderItem, nil
		}

		m.replyFunc = func(room id.RoomID, evt id.EventID, content string, asHtml bool) id.EventID {
			return evt
		}

		evt := &event.Event{Sender: id.UserID(c.matrixUsername)}

		handler := handlers["add"]

		err := handler(ctx, m, repo, nil, evt, c.message)
		if c.expectedErr != nil {
			assert.Equal(t, c.expectedErr, err)
		} else {
			assert.NoError(t, err)
		}

		if c.expectedReply != nil {
			calls := m.replyCalls()
			assert.Equal(t, 1, len(calls))
			assert.Equal(t, *c.expectedReply, calls[0].Content)
		}
	}
}

func TestMarkPaidHandler(t *testing.T) {
	user := &entity.User{Uuid: ptr.To(uuid.Must(uuid.FromString("bdfc1ef0-2388-4e39-a482-a85914e04140"))), Name: "@test:matrix.org"}
	user1 := &entity.User{Uuid: ptr.To(uuid.Must(uuid.FromString("567e17cd-3569-44e0-90bd-af17d27f8d80"))), Name: "@test1:matrix.org"}

	menus := map[string]*entity.Menu{
		"sangam_open": {
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			Name: "sangam_open",
		},
		"sangam_ordered": {
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			Name: "sangam_ordered",
		},
		"sangam_finalized": {
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			Name: "sangam_finalized",
		},
		"sangam_delivered": {
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			Name: "sangam_delivered",
		},
		"sangam_finalized_sp_test": {
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			Name: "sangam_finalized_sp_test",
		},
	}

	orders := []entity.Order{
		{
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			MenuUuid: menus["sangam_open"].Uuid,
			State: entity.Open,
		},
		{
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			MenuUuid: menus["sangam_ordered"].Uuid,
			State: entity.Ordered,
		},
		{
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			MenuUuid: menus["sangam_finalized"].Uuid,
			State: entity.Finalized,
		},
		{
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			MenuUuid: menus["sangam_delivered"].Uuid,
			State: entity.Delivered,
		},
		{
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			MenuUuid: menus["sangam_finalized_sp_test"].Uuid,
			State: entity.Finalized,
			SugarPerson: user.Uuid,
		},
	}

	type testCase struct {
		message        string
		user           *entity.User
		matrixUsername string
		items          []entity.OrderItem
		expectedErr    error
		expectedReply  *string
	}

	testCases := []testCase{
		{
			message:        "toggle_paid",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("message must be in the format 'toggle_paid [menu_name] [username]'"),
			expectedReply:  nil,
		},
		{
			message:        "toggle_paid sangam_open @test:matrix.org",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("no order items for user '@test:matrix.org' in order 'sangam_open' found"),
			expectedReply:  nil,
		},
		{
			message:        "toggle_paid sangam_open @test:matrix.org",
			user:           user,
			matrixUsername: "@test:matrix.org",
			items:          []entity.OrderItem{
				{
					Uuid: ptr.To(uuid.Must(uuid.NewV4())),
					User: user.Uuid,
					OrderUuid: orders[2].Uuid,
					Paid: false,
				},
			},
			expectedErr:    errors.New("could not mark order items as paid, because sugar person has not been set yet"),
			expectedReply:  nil,
		},
		{
			message:        "toggle_paid sangam_finalized @test:matrix.org",
			user:           user,
			matrixUsername: "@test:matrix.org",
			items:          []entity.OrderItem{
				{
					Uuid: ptr.To(uuid.Must(uuid.NewV4())),
					User: user.Uuid,
					OrderUuid: orders[2].Uuid,
				},
			},
			expectedErr:    errors.New("could not mark order items as paid, because sugar person has not been set yet"),
			expectedReply:  nil,
		},
		{
			message:        "toggle_paid sangam_finalized_sp_test @test:matrix.org",
			user:           user,
			matrixUsername: "@test:matrix.org",
			items:          []entity.OrderItem{
				{
					Uuid: ptr.To(uuid.Must(uuid.NewV4())),
					User: user1.Uuid,
					OrderUuid: orders[4].Uuid,
					Paid: false,
				},
				{
					Uuid: ptr.To(uuid.Must(uuid.NewV4())),
					User: user.Uuid,
					OrderUuid: orders[4].Uuid,
					Paid: false,
				},
			},
			expectedErr:    nil,
			expectedReply:  ptr.To("successfully marked all items of user '@test:matrix.org' in order 'sangam_finalized_sp_test' as paid"),
		},
		{
			message:        "toggle_paid sangam_finalized_sp_test @test1:matrix.org",
			user:           user,
			matrixUsername: "@test1:matrix.org",
			items:          []entity.OrderItem{
				{
					Uuid: ptr.To(uuid.Must(uuid.NewV4())),
					User: user1.Uuid,
					OrderUuid: orders[4].Uuid,
					Paid: false,
				},
				{
					Uuid: ptr.To(uuid.Must(uuid.NewV4())),
					User: user1.Uuid,
					OrderUuid: orders[4].Uuid,
					Paid: false,
				},
			},
			expectedErr:    nil,
			expectedReply:  ptr.To("successfully marked all items of user '@test1:matrix.org' in order 'sangam_finalized_sp_test' as paid"),
		},
	}

	for _, c := range testCases {
		ctx := context.Background()
		repo := &entity.RepositoryMock{}
		m := &MatrixBoundaryMock{}

		m.getUserByUsernameFunc = func(tx *gorm.DB, username string) (*entity.User, error) {
			if username == c.user.Name {
				return c.user, nil
			}
			if username == user1.Name {
				return user1, nil
			}
			return nil, entity.ErrUserNotFound
		}

		repo.GetMenuByNameFunc = func(tx *gorm.DB, name string) (*entity.Menu, error) {
			if menu := menus[name]; menu != nil {
				return menu, nil
			}
			return nil, entity.ErrMenuNotFound
		}

		repo.GetActiveOrderByMenuNameFunc = func(tx *gorm.DB, name string) (*entity.Order, error) {
			var foundMenuUuid *uuid.UUID
			for _, menu := range menus {
				if menu.Name == name {
					foundMenuUuid = menu.Uuid
				}
			}

			if foundMenuUuid == nil {
				return nil, entity.ErrMenuNotFound
			}

			for _, order := range orders {
				if *order.MenuUuid == *foundMenuUuid && order.State != entity.Delivered {
					return &order, nil
				}
			}
			return nil, entity.ErrOrderNotFound
		}

		repo.UpdateOrderFunc = func(tx *gorm.DB, uuidMoqParam, currentUser *uuid.UUID, orderMoqParam *entity.Order) (*entity.Order, error) {
			var foundOrder *entity.Order
			for _, order := range orders {
				if *order.Uuid == *orderMoqParam.Uuid {
					foundOrder = &order
				}
			}

			if foundOrder == nil {
				return nil, entity.ErrOrderNotFound
			}

			if foundOrder.SugarPerson != nil && *foundOrder.SugarPerson != *orderMoqParam.SugarPerson {
				return nil, entity.ErrSugarPersonChangeForbidden
			}

			return orderMoqParam, nil
		}

		repo.GetAllOrderItemsForOrderAndUserFunc = func(tx *gorm.DB, orderUuid, userUuid *uuid.UUID) ([]entity.OrderItem, error) {
			var foundItems []entity.OrderItem
			for _, item := range c.items {
				if *item.User == *userUuid {
					foundItems = append(foundItems, item)
				}
			}

			return foundItems, nil
		}

		repo.UpdateOrderItemFunc = func(tx *gorm.DB, orderItemUuid, userUuid *uuid.UUID, orderItem *entity.OrderItem) (*entity.OrderItem, error) {
			var foundItem *entity.OrderItem
			for _, item := range c.items {
				if *item.Uuid == *orderItem.Uuid {
					foundItem = &item
				}
			}

			if foundItem == nil {
				return nil, entity.ErrOrderItemNotFound
			}

			var foundOrder *entity.Order
			for _, order := range orders {
				if *order.Uuid == *foundItem.OrderUuid {
					foundOrder = &order
				}
			}

			if foundOrder == nil {
				return nil, entity.ErrOrderNotFound
			}

			if foundOrder.SugarPerson == nil {
				return nil, entity.ErrSugarPersonNotSet
			}

			return orderItem, nil
		}

		m.replyFunc = func(room id.RoomID, evt id.EventID, content string, asHtml bool) id.EventID {
			return evt
		}

		evt := &event.Event{Sender: id.UserID(c.matrixUsername)}

		handler := handlers[strings.Split(c.message, " ")[0]]

		err := handler(ctx, m, repo, nil, evt, c.message)
		if c.expectedErr != nil {
			assert.Equal(t, c.expectedErr, err)
		} else {
			assert.NoError(t, err)
		}

		if c.expectedReply != nil {
			calls := m.replyCalls()
			assert.Equal(t, 1, len(calls))
			assert.Equal(t, *c.expectedReply, calls[0].Content)
		}
	}
}

func TestPaidHandler(t *testing.T) {
	user := &entity.User{Uuid: ptr.To(uuid.Must(uuid.FromString("bdfc1ef0-2388-4e39-a482-a85914e04140"))), Name: "@test:matrix.org"}
	user1 := &entity.User{Uuid: ptr.To(uuid.Must(uuid.FromString("567e17cd-3569-44e0-90bd-af17d27f8d80"))), Name: "@test1:matrix.org"}

	menus := map[string]*entity.Menu{
		"sangam_open": {
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			Name: "sangam_open",
		},
		"sangam_ordered": {
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			Name: "sangam_ordered",
		},
		"sangam_finalized": {
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			Name: "sangam_finalized",
		},
		"sangam_delivered": {
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			Name: "sangam_delivered",
		},
		"sangam_finalized_sp_test": {
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			Name: "sangam_finalized_sp_test",
		},
	}

	orders := []entity.Order{
		{
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			MenuUuid: menus["sangam_open"].Uuid,
			State: entity.Open,
		},
		{
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			MenuUuid: menus["sangam_ordered"].Uuid,
			State: entity.Ordered,
		},
		{
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			MenuUuid: menus["sangam_finalized"].Uuid,
			State: entity.Finalized,
		},
		{
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			MenuUuid: menus["sangam_delivered"].Uuid,
			State: entity.Delivered,
		},
		{
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			MenuUuid: menus["sangam_finalized_sp_test"].Uuid,
			State: entity.Finalized,
			SugarPerson: user.Uuid,
		},
	}

	type testCase struct {
		message        string
		user           *entity.User
		matrixUsername string
		expectedErr    error
		expectedReply  *string
	}

	testCases := []testCase{
		{
			message:        "paid",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("message must be in the format 'paid [menu_name]'"),
			expectedReply:  nil,
		},
		{
			message:        "paid sangam_open",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    nil,
			expectedReply:  ptr.To("You are now the sugar person. This cannot be undone!"),
		},
		{
			message:        "paid sangam_finalized",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    nil,
			expectedReply:  ptr.To("You are now the sugar person. This cannot be undone!"),
		},
		{
			message:        "paid sangam_ordered",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    nil,
			expectedReply:  ptr.To("You are now the sugar person. This cannot be undone!"),
		},
		{
			message:        "paid sangam_delivered",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("there is no active order for menu 'sangam_delivered'"),
			expectedReply:  nil,
		},
		{
			message:        "paid sangam_finalized_sp_test",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    nil,
			expectedReply:  ptr.To("You are now the sugar person. This cannot be undone!"),
		},
		{
			message:        "paid sangam_finalized_sp_test",
			user:           user1,
			matrixUsername: "@test1:matrix.org",
			expectedErr:    errors.New("The sugar person cannot be changed after it has been set"),
			expectedReply:  nil,
		},
	}

	for _, c := range testCases {
		ctx := context.Background()
		repo := &entity.RepositoryMock{}
		m := &MatrixBoundaryMock{}

		m.getUserByUsernameFunc = func(tx *gorm.DB, username string) (*entity.User, error) {
			if username == c.user.Name {
				return c.user, nil
			}
			return nil, entity.ErrUserNotFound
		}

		repo.GetMenuByNameFunc = func(tx *gorm.DB, name string) (*entity.Menu, error) {
			if menu := menus[name]; menu != nil {
				return menu, nil
			}
			return nil, entity.ErrMenuNotFound
		}

		repo.GetActiveOrderByMenuNameFunc = func(tx *gorm.DB, name string) (*entity.Order, error) {
			var foundMenuUuid *uuid.UUID
			for _, menu := range menus {
				if menu.Name == name {
					foundMenuUuid = menu.Uuid
				}
			}

			if foundMenuUuid == nil {
				return nil, entity.ErrMenuNotFound
			}

			for _, order := range orders {
				if *order.MenuUuid == *foundMenuUuid && order.State != entity.Delivered {
					return &order, nil
				}
			}
			return nil, entity.ErrOrderNotFound
		}

		repo.UpdateOrderFunc = func(tx *gorm.DB, uuidMoqParam, currentUser *uuid.UUID, orderMoqParam *entity.Order) (*entity.Order, error) {
			var foundOrder *entity.Order
			for _, order := range orders {
				if *order.Uuid == *orderMoqParam.Uuid {
					foundOrder = &order
				}
			}

			if foundOrder == nil {
				return nil, entity.ErrOrderNotFound
			}

			if foundOrder.SugarPerson != nil && *foundOrder.SugarPerson != *orderMoqParam.SugarPerson {
				return nil, entity.ErrSugarPersonChangeForbidden
			}

			return orderMoqParam, nil
		}

		m.replyFunc = func(room id.RoomID, evt id.EventID, content string, asHtml bool) id.EventID {
			return evt
		}

		evt := &event.Event{Sender: id.UserID(c.matrixUsername)}

		handler := handlers[strings.Split(c.message, " ")[0]]

		err := handler(ctx, m, repo, nil, evt, c.message)
		if c.expectedErr != nil {
			assert.Equal(t, c.expectedErr, err)
		} else {
			assert.NoError(t, err)
		}

		if c.expectedReply != nil {
			calls := m.replyCalls()
			assert.Equal(t, 1, len(calls))
			assert.Equal(t, *c.expectedReply, calls[0].Content)
		}
	}
}

func TestStateTransitionHandler(t *testing.T) {
	user := &entity.User{Uuid: ptr.To(uuid.Must(uuid.FromString("bdfc1ef0-2388-4e39-a482-a85914e04140"))), Name: "@test:matrix.org"}

	menus := map[string]*entity.Menu{
		"sangam_open": {
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			Name: "sangam_open",
		},
		"sangam_ordered": {
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			Name: "sangam_ordered",
		},
		"sangam_finalized": {
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			Name: "sangam_finalized",
		},
		"sangam_delivered": {
			Uuid: ptr.To(uuid.Must(uuid.NewV4())),
			Name: "sangam_delivered",
		},
	}

	orders := []entity.Order{
		{
			MenuUuid: menus["sangam_open"].Uuid,
			State: entity.Open,
		},
		{
			MenuUuid: menus["sangam_ordered"].Uuid,
			State: entity.Ordered,
		},
		{
			MenuUuid: menus["sangam_finalized"].Uuid,
			State: entity.Finalized,
		},
		{
			MenuUuid: menus["sangam_delivered"].Uuid,
			State: entity.Delivered,
		},
	}

	type testCase struct {
		message        string
		user           *entity.User
		matrixUsername string
		expectedErr    error
		expectedReply  *string
	}

	testCases := []testCase{
		{
			message:        "finalize",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("menu must be specified"),
			expectedReply:  nil,
		},
		{
			message:        "finalize sangam_open",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    nil,
			expectedReply:  ptr.To("successfully set state of order 'sangam_open' to finalized"),
		},
		{
			message:        "finalize sangam_delivered",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("there is no active order for menu 'sangam_delivered'"),
			expectedReply:  nil,
		},
		{
			message:        "finalize sangam_ordered",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    nil,
			expectedReply:  ptr.To("successfully set state of order 'sangam_ordered' to finalized"),
		},
		{
			message:        "re-open",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("menu must be specified"),
			expectedReply:  nil,
		},
		{
			message:        "re-open sangam_open",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("order 'sangam_open' is already in state open"),
			expectedReply:  nil,
		},
		{
			message:        "re-open sangam_delivered",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("there is no active order for menu 'sangam_delivered'"),
			expectedReply:  nil,
		},
		{
			message:        "re-open sangam_ordered",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    nil,
			expectedReply:  ptr.To("successfully set state of order 'sangam_ordered' to open"),
		},
		{
			message:        "ordered",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("menu must be specified"),
			expectedReply:  nil,
		},
		{
			message:        "ordered sangam_open",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    nil,
			expectedReply:  ptr.To("successfully set state of order 'sangam_open' to ordered"),
		},
		{
			message:        "ordered sangam_delivered",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("there is no active order for menu 'sangam_delivered'"),
			expectedReply:  nil,
		},
		{
			message:        "ordered sangam_ordered",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("order 'sangam_ordered' is already in state ordered"),
			expectedReply:  nil,
		},
		{
			message:        "delivered",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("menu must be specified"),
			expectedReply:  nil,
		},
		{
			message:        "delivered sangam_open",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    nil,
			expectedReply:  ptr.To("successfully set state of order 'sangam_open' to delivered"),
		},
		{
			message:        "delivered sangam_delivered",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    errors.New("there is no active order for menu 'sangam_delivered'"),
			expectedReply:  nil,
		},
		{
			message:        "delivered sangam_ordered",
			user:           user,
			matrixUsername: "@test:matrix.org",
			expectedErr:    nil,
			expectedReply:  ptr.To("successfully set state of order 'sangam_ordered' to delivered"),
		},
	}

	for _, c := range testCases {
		ctx := context.Background()
		repo := &entity.RepositoryMock{}
		m := &MatrixBoundaryMock{}

		m.getUserByUsernameFunc = func(tx *gorm.DB, username string) (*entity.User, error) {
			if username == c.user.Name {
				return c.user, nil
			}
			return nil, entity.ErrUserNotFound
		}

		repo.GetMenuByNameFunc = func(tx *gorm.DB, name string) (*entity.Menu, error) {
			if menu := menus[name]; menu != nil {
				return menu, nil
			}
			return nil, entity.ErrMenuNotFound
		}

		repo.GetActiveOrderByMenuNameFunc = func(tx *gorm.DB, name string) (*entity.Order, error) {
			var foundMenuUuid *uuid.UUID
			for _, menu := range menus {
				if menu.Name == name {
					foundMenuUuid = menu.Uuid
				}
			}

			if foundMenuUuid == nil {
				return nil, entity.ErrMenuNotFound
			}

			for _, order := range orders {
				if *order.MenuUuid == *foundMenuUuid && order.State != entity.Delivered {
					return &order, nil
				}
			}
			return nil, entity.ErrOrderNotFound
		}

		repo.UpdateOrderFunc = func(tx *gorm.DB, uuidMoqParam, currentUser *uuid.UUID, order *entity.Order) (*entity.Order, error) {
			return order, nil
		}

		m.replyFunc = func(room id.RoomID, evt id.EventID, content string, asHtml bool) id.EventID {
			return evt
		}

		evt := &event.Event{Sender: id.UserID(c.matrixUsername)}

		handler := handlers[strings.Split(c.message, " ")[0]]

		err := handler(ctx, m, repo, nil, evt, c.message)
		if c.expectedErr != nil {
			assert.Equal(t, c.expectedErr, err)
		} else {
			assert.NoError(t, err)
		}

		if c.expectedReply != nil {
			calls := m.replyCalls()
			assert.Equal(t, 1, len(calls))
			assert.Equal(t, *c.expectedReply, calls[0].Content)
		}
	}
}
