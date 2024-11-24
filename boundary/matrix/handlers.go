package matrix

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/ssh"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
	"gorm.io/gorm"
	"maunium.net/go/mautrix/event"
)

type CommandHandler = func(*MatrixBoundary, *gorm.DB, *event.Event, string) error

var handlers = map[string]CommandHandler{
	"help":           handleHelp,
	"register":       handleRegister,
	"set_public_key": handleSetPublicKey,
	"start":          handleNewOrder,
	"add":            handleNewOrderItem,
	"paid":           handlePaid,
	"toggle_paid":    handleMarkPaid,
}
var (
	ErrParsingPublicKey = errors.New("cannot parse public key")
)

func handleUnrecognizedCommand(m *MatrixBoundary, _ *gorm.DB, evt *event.Event, message string) error {
	m.reply(evt.RoomID, evt.ID, fmt.Sprintf("Command not recognized: %s", message), false)
	return nil
}

func handleRegister(m *MatrixBoundary, tx *gorm.DB, evt *event.Event, message string) error {
	username := evt.Sender.String()
	matrixUser, err := m.repo.GetMatrixUserByUsername(tx, username)
	if errors.Is(err, entity.ErrOrderNotFound) {
		user, err := m.repo.CreateUser(tx, &entity.User{Name: username})
		if err != nil {
			msg := fmt.Sprintf("could not create user for sender '%s'", username)
			log.Ctx(m.ctx).Warn().Err(err).Msg(msg)
			return errors.New(msg)
		}
		matrixUser, err = m.repo.CreateMatrixUser(tx, &entity.MatrixUser{UserUuid: user.Uuid, Username: username})
		if err != nil {
			msg := fmt.Sprintf("could not create matrix user for sender '%s'", username)
			log.Ctx(m.ctx).Warn().Err(err).Msg(msg)
			return errors.New(msg)
		}

		m.reply(evt.RoomID, evt.ID, fmt.Sprintf("successfully registered user: %s", user.Name), false)
		return nil
	} else if matrixUser != nil {
		msg := fmt.Sprintf("user '%s' is already registered", username)
		log.Ctx(m.ctx).Warn().Msg(msg)
		return errors.New(msg)
	} else {
		msg := fmt.Sprintf("error occured while registering user '%s'", username)
		log.Ctx(m.ctx).Warn().Err(err).Msg(msg)
		return errors.New(msg)
	}
}

func handleHelp(m *MatrixBoundary, tx *gorm.DB, evt *event.Event, message string) error {
	m.reply(evt.RoomID, evt.ID, "Hello world", false)
	return nil
}

func handleSetPublicKey(m *MatrixBoundary, tx *gorm.DB, evt *event.Event, message string) error {
	username := evt.Sender.String()
	publicKey := strings.TrimPrefix(message, "set_public_key ")
	if publicKey == "" {
		return errors.New("public key must not be empty")
	}
	publicKeySegments := strings.Split(publicKey, " ")
	if len(publicKeySegments) != 2 {
		return errors.New("public key must be in the format 'key_type base64_encoded_key'")
	}
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeySegments[1])
	if err != nil {
		log.Ctx(m.ctx).Warn().Err(err).Msg("could not parse public key")
		return ErrParsingPublicKey
	}
	_, err = ssh.ParsePublicKey(publicKeyBytes)
	if err != nil {
		log.Ctx(m.ctx).Warn().Err(err).Msg("could not parse public key")
		return ErrParsingPublicKey
	}

	user, err := m.getUserByUsername(tx, username)
	if err != nil {
		return err
	}

	sshUser, err := m.repo.GetSshUser(tx, user.Uuid)
	if err != nil {
		msg := fmt.Sprintf("could not get ssh user for username '%s'", username)
		log.Ctx(m.ctx).Warn().Err(err).Msg(msg)
		return errors.New(msg)
	}

	sshUser.PublicKey = publicKey
	_, err = m.repo.UpdateSshUser(tx, user.Uuid, sshUser)
	if err != nil {
		msg := fmt.Sprintf("could not set public key for user '%s'", user.Name)
		log.Ctx(m.ctx).Warn().Err(err).Msg(msg)
		return errors.New(msg)
	}

	m.reply(evt.RoomID, evt.ID, fmt.Sprintf("successfully set ssh public key for user: %s", username), false)
	return nil
}

func handleNewOrder(m *MatrixBoundary, tx *gorm.DB, evt *event.Event, message string) error {
	username := evt.Sender.String()

	initiator, err := m.getUserByUsername(tx, username)
	if err != nil {
		return err
	}

	menuName := strings.TrimPrefix(message, "start ")
	menu, err := m.repo.GetMenuByName(tx, menuName)
	if err != nil {
		msg := fmt.Sprintf("could not get menu '%s'", menuName)
		log.Ctx(m.ctx).Warn().Err(err).Msg(msg)
		return errors.New(msg)
	}

	_, err = m.repo.CreateOrder(tx, &entity.Order{Initiator: initiator.Uuid, MenuUuid: menu.Uuid})
	if errors.Is(err, entity.ErrActiveOrderForMenuAlreadyExists) {
		msg := fmt.Sprintf("there is already an active order for menu '%s'", menuName)
		log.Ctx(m.ctx).Warn().Err(err).Msg(msg)
		return errors.New(msg)
	} else if err != nil {
		msg := "could not create order"
		log.Ctx(m.ctx).Warn().Err(err).Msg(msg)
		return errors.New(msg)
	}

	m.reply(evt.RoomID, evt.ID, fmt.Sprintf("started new order for %s", menuName), false)
	return nil
}

func handleNewOrderItem(m *MatrixBoundary, tx *gorm.DB, evt *event.Event, message string) error {
	username := evt.Sender.String()

	user, err := m.getUserByUsername(tx, username)
	if err != nil {
		return err
	}

	message = strings.TrimPrefix(message, "add ")
	args := strings.Split(message, " ")
	if len(args) != 2 {
		log.Ctx(m.ctx).Warn().Msgf("message '%s' wasn't formatted correctly", message)
		return errors.New("message must be in the format 'add [menu_name] [short_name]'")
	}
	menuName := args[0]

	order, err := m.repo.GetActiveOrderByMenuName(tx, menuName)
	if err != nil {
		msg := fmt.Sprintf("there is no active order for menu '%s'", menuName)
		log.Ctx(m.ctx).Warn().Err(err).Msg(msg)
		return errors.New(msg)
	}

	shortName := args[1]
	menuItem, err := m.repo.GetMenuItemByShortName(tx, order.MenuUuid, shortName)
	if err != nil {
		msg := fmt.Sprintf("could not get menu item '%s'", shortName)
		log.Ctx(m.ctx).Warn().Err(err).Msg(msg)
		return errors.New(msg)
	}

	_, err = m.repo.CreateOrderItem(tx, order.Uuid, &entity.OrderItem{OrderUuid: order.Uuid, User: user.Uuid, MenuItemUuid: menuItem.Uuid})
	if errors.Is(err, entity.ErrOrderNotOpen) {
		msg := "order is not open"
		log.Ctx(m.ctx).Warn().Err(err).Msg(msg)
		return errors.New(msg)
	} else if err != nil {
		msg := "could not create order item"
		log.Ctx(m.ctx).Warn().Err(err).Msg(msg)
		return errors.New(msg)
	}

	m.reply(evt.RoomID, evt.ID, fmt.Sprintf("added '%s' to order '%s'", menuItem.Name, menuName), false)

	return nil
}

func handlePaid(m *MatrixBoundary, tx *gorm.DB, evt *event.Event, message string) error {
	username := evt.Sender.String()

	user, err := m.getUserByUsername(tx, username)
	if err != nil {
		return err
	}

	message = strings.TrimPrefix(message, "paid ")
	args := strings.Split(message, " ")
	if len(args) != 1 {
		log.Ctx(m.ctx).Warn().Msgf("message '%s' wasn't formatted correctly", message)
		return errors.New("message must be in the format 'paid [menu_name]'")
	}
	menuName := args[0]

	order, err := m.repo.GetActiveOrderByMenuName(tx, menuName)
	if err != nil {
		msg := fmt.Sprintf("there is no active order for menu '%s'", menuName)
		log.Ctx(m.ctx).Warn().Err(err).Msg(msg)
		return errors.New(msg)
	}

	order.SugarPerson = user.Uuid

	_, err = m.repo.UpdateOrder(tx, order.Uuid, user.Uuid, order)
	if errors.Is(err, entity.ErrSugarPersonChangeForbidden) {
		msg := "The sugar person cannot be changed after it has been set"
		log.Ctx(m.ctx).Warn().Err(err).Msg(msg)
		return errors.New(msg)
	} else if err != nil {
		msg := "could not update order"
		log.Ctx(m.ctx).Warn().Err(err).Msg(msg)
		return errors.New(msg)
	}

	m.reply(evt.RoomID, evt.ID, "You are now the sugar person. This cannot be undone!", false)

	return nil
}

func handleMarkPaid(m *MatrixBoundary, tx *gorm.DB, evt *event.Event, message string) error {
	currentUsername := evt.Sender.String()

	currentUser, err := m.getUserByUsername(tx, currentUsername)
	if err != nil {
		return err
	}

	message = strings.TrimPrefix(message, "toggle_paid ")
	args := strings.Split(message, " ")
	if len(args) != 2 {
		log.Ctx(m.ctx).Warn().Msgf("message '%s' wasn't formatted correctly", message)
		return errors.New("message must be in the format 'toggle_paid [menu_name] [username]'")
	}
	menuName := args[0]

	order, err := m.repo.GetActiveOrderByMenuName(tx, menuName)
	if err != nil {
		msg := fmt.Sprintf("there is no active order for menu '%s'", menuName)
		log.Ctx(m.ctx).Warn().Err(err).Msg(msg)
		return errors.New(msg)
	}

	usernameParam := args[1]
	user, err := m.repo.GetUserByName(tx, usernameParam)
	if err != nil {
		msg := fmt.Sprintf("could not get user '%s'", usernameParam)
		log.Ctx(m.ctx).Warn().Err(err).Msg(msg)
		return errors.New(msg)
	}

	orderItems, err := m.repo.GetAllOrderItemsForOrderAndUser(tx, order.Uuid, user.Uuid)
	if err != nil {
		msg := "could not get order items for user"
		log.Ctx(m.ctx).Warn().Err(err).Msg(msg)
		return errors.New(msg)
	}

	log.Ctx(m.ctx).Info().Msgf("found %d order items for user '%s' in order '%s'", len(orderItems), user.Name, menuName)

	allPaid := true
	for _, oi := range orderItems {
		if !oi.Paid {
			allPaid = false
		}
	}

	var paidStatusStr string
	if !allPaid {
		paidStatusStr = "paid"
	} else {
		paidStatusStr = "not paid"
	}

	for _, existingOrderItem := range orderItems {
		existingOrderItem.Paid = !allPaid
		_, err = m.repo.UpdateOrderItem(tx, existingOrderItem.Uuid, currentUser.Uuid, &existingOrderItem)
		if err != nil {
			msg := fmt.Sprintf("could not mark order items as %s", paidStatusStr)
			log.Ctx(m.ctx).Warn().Err(err).Msg(msg)
			return errors.New(msg)
		}
	}

	m.reply(evt.RoomID, evt.ID, fmt.Sprintf("successfully marked all items of user '%s' in order '%s' as %s", user.Name, menuName, paidStatusStr), false)

	return nil
}
