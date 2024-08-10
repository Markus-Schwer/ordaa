package tui

import (
	"fmt"

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
	s            lipgloss.Style
	it           textinput.Model
	t            table.Model
	menu         *entity.Menu
	ready        bool
	searching    bool
	searchString string
}

func NewMenuModel(info *LayoutInfo) *MenuModel {
	columns := []table.Column{
		{Title: "ID", Width: 6},
		{Title: "NAME", Width: 20},
		{Title: "PRICE", Width: 10},
	}
	t := table.New(
		table.WithColumns(columns),
		// table.WithHeight(info.height-info.headerOffset-info.footerOffset),
		table.WithHeight(20),
		table.WithRows([]table.Row{table.Row{"42", "chicken tikka masala", "1590"}}),
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
		LayoutInfo:   info,
		s:            info.quitStyle,
		it:           it,
		t:            t,
		ready:        false,
		searching:    true,
		searchString: "",
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
		// m.menu = msg
		m.ready = true
		rows := make([]table.Row, 0)
		for _, v := range msg.Items {
			rows = append(rows, table.Row{v.ShortName, v.Name, fmt.Sprintf("%4.2f€", float64(v.Price)/100.0)})
		}
		m.t.SetRows(rows)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Up):
			if !m.searching {
				if m.t.Cursor() == 0 {
					m.searching = true
				} else {
					m.t, cmd = m.t.Update(msg)
				}
			}
		case msg.String() == "down":
			// default keymap down can't be used here because it would interpret
			// typing a 'j' as moving downwards
			if m.searching {
				m.searching = false
			} else {
				m.t, cmd = m.t.Update(msg)
			}
		default:
			if m.searching {
				m.it, cmd = m.it.Update(msg)
			} else {
				m.t, cmd = m.t.Update(msg)
			}
		}
	}
	return m, cmd
}

func (m *MenuModel) View() string {
	if !m.ready {
		return "Loading..."
	}
	if m.searching {
		m.it.Focus()
	} else {
		m.t.Focus()
	}
	// rows := make([][]string, 0)
	// for _, item := range m.menu.Items {
	// 	if strings.Contains(strings.ToLower(item.Name), strings.ToLower(m.it.Value())) ||
	// 		strings.Contains(strings.ToLower(item.ShortName), strings.ToLower(m.it.Value())) {
	// 		rows = append(rows, []string{item.ShortName, item.Name, })
	// 	}
	// }
	// table := table.New().
	// 	BorderStyle(m.txtStyle).
	// 	Width(m.width-2).
	// 	Headers("SHORT", "NAME", "PRICE").
	// 	Rows(rows...)
	// search := m.it.View() + "\n\n"
	// m.viewport.SetContent(m.s.Render(table.Render()))
	// m.viewport.Height = m.height - m.headerOffset - m.footerOffset - lipgloss.Height(search) - 2
	// m.viewport.YPosition = m.headerOffset
	// m.it.Focus()
	// return search + m.viewport.View()

	// use different styles here
	return fmt.Sprintf("%s\n%s", m.txtStyle.Render(m.it.View()), m.quitStyle.Render(m.t.View()) )
	// return m.txtStyle.Render(m.t.View())
}
