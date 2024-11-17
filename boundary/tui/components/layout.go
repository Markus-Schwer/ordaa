package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/gofrs/uuid"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
)

type ComponentKey string

type LayoutEvent struct {
	Component ComponentKey
	Uuid      *uuid.UUID
}

type Layout interface {
	BoxStyle() lipgloss.Style
	ContentWidth() int
	ContentHeight() int
	Repository() entity.Repository
}
