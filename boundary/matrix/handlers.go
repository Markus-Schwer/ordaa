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
	"help":                   handleHelp,
	"register":               handleRegister,
	"set_public_key":         handleSetPublicKey,
}

func handleUnrecognizedCommand(m *MatrixBoundary, tx *gorm.DB, evt *event.Event, message string) error {
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

		m.reply(evt.RoomID, evt.ID, fmt.Sprintf("successfully registered user: %s", username), false)
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

	matrixUser, err := m.repo.GetMatrixUserByUsername(tx, username)
	if err != nil {
		return fmt.Errorf("could not get matrix user of sender '%s' for message '%s': %w", username, message, err)
	}

	user, err := m.repo.GetUser(tx, matrixUser.UserUuid)
	if err != nil {
		return fmt.Errorf("could not get user of sender '%s' for message '%s': %w", username, message, err)
	}

	user.PublicKey = publicKey
	_, err = m.repo.UpdateUser(tx, user.Uuid, user)
	if err != nil {
		return fmt.Errorf("could not set public key for user '%s': %w", username, err)
	}

	m.reply(evt.RoomID, evt.ID, fmt.Sprintf("successfully set ssh public key for user: %s", username), false)
	return nil
}
