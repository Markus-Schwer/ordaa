package matrix

import (
	"errors"
	"fmt"
	"strings"

	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
	"gorm.io/gorm"
	"maunium.net/go/mautrix/event"
)

type CommandHandler = func(*MatrixBoundary, *gorm.DB, *event.Event, string) error

var handlers = map[string]CommandHandler{
	"help":           handleHelp,
	"register":       handleRegister,
	"set_public_key": handleSetPublicKey,
	"new_order":      handleNewOrder,
}

func handleUnrecognizedCommand(m *MatrixBoundary, _ *gorm.DB, evt *event.Event, message string) error {
	m.reply(evt.RoomID, evt.ID, fmt.Sprintf("Command not recognized: %s", message), false)
	return nil
}

func handleRegister(m *MatrixBoundary, tx *gorm.DB, evt *event.Event, message string) error {
	username := evt.Sender.String()
	matrixUser, err := m.repo.GetMatrixUserByUsername(tx, username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user, err := m.repo.CreateUser(tx, &entity.User{Name: username})
		if err != nil {
			return fmt.Errorf("could not create user for sender '%s': %w", username, err)
		}
		matrixUser, err = m.repo.CreateMatrixUser(tx, &entity.MatrixUser{UserUuid: user.Uuid, Username: username})
		if err != nil {
			return fmt.Errorf("could not create matrix user for sender '%s': %w", username, err)
		}

		m.reply(evt.RoomID, evt.ID, fmt.Sprintf("successfully registered user: %s", user.Name), false)
		return nil
	} else if matrixUser != nil {
		return errors.New(fmt.Sprintf("user '%s' is already registered", username))
	} else {
		return fmt.Errorf("error occured while registering user '%s': %w", username, err)
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

	user, err := m.getUserByUsername(tx, username)
	if err != nil {
		return err
	}

	user.PublicKey = publicKey
	_, err = m.repo.UpdateUser(tx, user.Uuid, user)
	if err != nil {
		return fmt.Errorf("could not set public key for user '%s': %w", user.Name, err)
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

	menuName := strings.TrimPrefix(message, "new_order ")
	menu, err := m.repo.GetMenuByName(tx, menuName)
	if err != nil {
		return fmt.Errorf("could not get menu '%s': %w", menuName, err)
	}

	_, err = m.repo.CreateOrder(tx, &entity.Order{Initiator: initiator.Uuid, MenuUuid: menu.Uuid})
	if err != nil {
		return fmt.Errorf("could not create order: %w", err)
	}

	m.reply(evt.RoomID, evt.ID, fmt.Sprintf("started new order for %s", menuName), false)
	return nil
}
