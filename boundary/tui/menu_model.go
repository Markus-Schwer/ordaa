package tui

import (
	"errors"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type MenuModel struct {
	*LayoutInfo
	s        lipgloss.Style
	table    *table.Table
	viewport viewport.Model
}

func NewMenuModel(info *LayoutInfo) *MenuModel {
	t := table.New().
		BorderStyle(info.txtStyle).
		Width(info.width).
		Headers("SHORT", "NAME", "PRICE")
	vp := viewport.New(info.width, info.height-info.headerOffset-info.footerOffset)
	return &MenuModel{
		LayoutInfo: info,
		s:          info.txtStyle,
		table:      t,
		viewport:   vp,
	}
}

func (m *MenuModel) Init() tea.Cmd {
	return nil
}

func (m *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *MenuModel) View() string {
	rows := make([][]string, 0)
	m.repo.Transaction(func(tx *gorm.DB) error {
		menus, err := m.repo.GetAllMenus(tx)
		if err != nil {
			log.Ctx(m.ctx).Error().Err(err).Msg("allMenus error getting menus")
			return err
		}
		if len(menus) == 0 {
			err := errors.New("no menus found")
			log.Ctx(m.ctx).Error().Err(err).Msg("there are no menus")
			return err
		}
		menu := menus[0]
		for _, item := range menu.Items {
			rows = append(rows, []string{item.ShortName, item.Name, string(item.Price)})
		}
		return nil
	})
	m.table = table.New().
		BorderStyle(m.txtStyle).
		Width(m.width).
		Headers("SHORT", "NAME", "PRICE").
		Rows(rows...)
	c := m.s.Render(m.table.Render())
	m.viewport.SetContent(c)
	m.viewport.Height = m.height - m.headerOffset - m.footerOffset
	m.viewport.YPosition = m.headerOffset
	return m.viewport.View()
}
