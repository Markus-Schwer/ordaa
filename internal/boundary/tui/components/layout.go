package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/gofrs/uuid"
	"github.com/Markus-Schwer/ordaa/internal/entity"
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
