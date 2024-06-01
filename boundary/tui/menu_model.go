package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
	"gorm.io/gorm"
)

type MenuModel struct {
	*LayoutInfo
	s        lipgloss.Style
	viewport viewport.Model
	it       textinput.Model
	menu     *entity.Menu
	ready    bool
}

func NewMenuModel(info *LayoutInfo) *MenuModel {
	vp := viewport.New(info.width-2, info.height-info.headerOffset-info.footerOffset)
	it := textinput.New()
	it.Placeholder = "Search"
	it.Width = 20
	return &MenuModel{
		LayoutInfo: info,
		s:          info.txtStyle,
		it:         it,
		viewport:   vp,
		ready:      false,
	}
}

func (m *MenuModel) Init() tea.Cmd {
	return func() tea.Msg {
		var menu *entity.Menu
		m.repo.Transaction(func(tx *gorm.DB) error {
			menus, err := m.repo.GetAllMenus(tx)
			if err != nil {
				log.Ctx(m.ctx).Error().Err(err).Msg("allMenus error getting menus")
				return err
			}
			if len(menus) == 0 {
				err := fmt.Errorf("there are no menus")
				log.Ctx(m.ctx).Error().Err(err).Msg("there are no menus")
				return err
			}
			menu = &menus[0]
			return nil
		})
		return menu
	}
}

func (m *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	log.Ctx(m.ctx).Debug().Msgf("menu model update: %v", msg)
	switch msg := msg.(type) {
	case *entity.Menu:
		m.menu = msg
		m.ready = true
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "down":
			m.viewport, cmd = m.viewport.Update(msg)
		}
	}
	m.it, cmd = m.it.Update(msg)
	return m, cmd
}

func (m *MenuModel) View() string {
	if !m.ready {
		return "Loading..."
	}
	rows := make([][]string, 0)
	for _, item := range m.menu.Items {
		if strings.Contains(strings.ToLower(item.Name), strings.ToLower(m.it.Value())) ||
			strings.Contains(strings.ToLower(item.ShortName), strings.ToLower(m.it.Value())) {
			rows = append(rows, []string{item.ShortName, item.Name, fmt.Sprintf("%4.2f€", float64(item.Price)/100.0)})
		}
	}
	table := table.New().
		BorderStyle(m.txtStyle).
		Width(m.width-2).
		Headers("SHORT", "NAME", "PRICE").
		Rows(rows...)
	search := m.it.View() + "\n\n"
	m.viewport.SetContent(m.s.Render(table.Render()))
	m.viewport.Height = m.height - m.headerOffset - m.footerOffset - lipgloss.Height(search) - 2
	m.viewport.YPosition = m.headerOffset
	m.it.Focus()
	return search + m.viewport.View()
}
