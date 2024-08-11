package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
	"gorm.io/gorm"
)

type MenuModel struct {
	*LayoutInfo
	s     lipgloss.Style
	it    textinput.Model
	t     table.Model
	menu  *entity.Menu
	ready bool
}

func NewMenuModel(info *LayoutInfo) *MenuModel {
	columns := []table.Column{
		{Title: "ID", Width: 6},
		{Title: "NAME", Width: 20},
		{Title: "PRICE", Width: 10},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	it := textinput.New()
	it.Placeholder = "Search"
	it.Width = 20
	return &MenuModel{
		LayoutInfo: info,
		s:          info.quitStyle,
		it:         it,
		t:          t,
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
	switch msg := msg.(type) {
	case *entity.Menu:
		m.menu = msg
		if m.menu != nil {
			m.ready = true
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.AlphaNum), msg.Type == tea.KeyBackspace, msg.Type == tea.KeySpace:
			m.it, cmd = m.it.Update(msg)
		case key.Matches(msg, DefaultKeyMap.Up), key.Matches(msg, DefaultKeyMap.Down):
			m.t, cmd = m.t.Update(msg)
		default:
		}
	}
	return m, cmd
}

func (m *MenuModel) View() string {
	if !m.ready {
		return "Loading..."
	}
	m.it.Focus()
	m.t.Focus()
	rows := make([]table.Row, 0)
	for _, v := range m.menu.Items {
		s := strings.ToLower(m.it.Value())
		if !(strings.Contains(strings.ToLower(v.ShortName), s) || strings.Contains(strings.ToLower(v.Name), s)) {
			continue
		}
		rows = append(rows, table.Row{v.ShortName, v.Name, fmt.Sprintf("%4.2f€", float64(v.Price)/100.0)})
	}
	m.t.SetRows(rows)
	m.t.SetHeight(m.height-m.headerOffset-m.footerOffset-2-1)

	// use different styles here
	return fmt.Sprintf("%s\n%s", m.txtStyle.Render(m.it.View()), m.quitStyle.Render(m.t.View()))
}
