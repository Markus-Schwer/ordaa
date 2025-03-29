package menu

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gofrs/uuid"
	"gitlab.com/sfz.aalen/hackwerk/ordaa/boundary/tui/components"
	"gitlab.com/sfz.aalen/hackwerk/ordaa/entity"
	"gorm.io/gorm"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

const MenuSelectorComponent components.ComponentKey = "menu_selector_component"

type model struct {
	ctx          context.Context
	layout       components.Layout
	list         list.Model
	initMenuUuid *uuid.UUID
	menu         *entity.Menu
}

func NewMenuItemSelectorModel(ctx context.Context, repo entity.Repository, layout components.Layout) tea.Model {
	l := list.New([]list.Item{}, newItemDelegate(ctx, repo), layout.ContentWidth(), layout.ContentHeight())
	l.Title = "Dish selection"
	l.Styles.Title = titleStyle
	l.SetShowStatusBar(true)
	l.StartSpinner()
	return &model{
		ctx:    ctx,
		layout: layout,
		list:   l,
	}
}

func (m *model) Init() tea.Cmd {
	return func() tea.Msg {
		var err error
		m.initMenuUuid, err = components.GetMenuUuidFromContext(m.ctx)
		if err != nil {
			panic(err)
		}
		var menu *entity.Menu
		m.layout.Repository().Transaction(func(tx *gorm.DB) error {
			var err error
			menu, err = m.layout.Repository().GetMenu(tx, m.initMenuUuid)
			return err
		})
		return menu
	}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case *entity.Menu:
		m.menu = msg
		if m.menu == nil {
			cmds = append(cmds, m.list.StartSpinner())
		} else {
			m.list.StopSpinner()
			cmds = append(cmds, m.list.SetItems(makeListItems(m.menu)))
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	m.list.SetWidth(m.layout.ContentWidth())
	return m.layout.BoxStyle().Render(m.list.View())
}
