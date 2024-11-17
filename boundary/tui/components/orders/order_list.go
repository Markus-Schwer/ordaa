package orders

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/tui/components"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
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

const OrderSelectorComponent components.ComponentKey = "order_selector_component"

type model struct {
	ctx    context.Context
	layout components.Layout
	list   list.Model
	orders []entity.Order
}

func NewOrderSelectorModel(ctx context.Context, repo entity.Repository, layout components.Layout) tea.Model {
	l := list.New([]list.Item{}, newItemDelegate(ctx, repo), layout.ContentWidth(), layout.ContentHeight())
	l.Title = "Order selection"
	l.Styles.Title = titleStyle
	l.StartSpinner()
	return &model{
		ctx:    ctx,
		layout: layout,
		list:   l,
	}
}

func (m *model) Init() tea.Cmd {
	return func() tea.Msg {
		var listItems = make([]listItem, 0)
		m.layout.Repository().Transaction(func(tx *gorm.DB) error {
			orders, err := m.layout.Repository().GetAllOrders(tx)
			log.Ctx(m.ctx).Info().Msgf("orders: %v", orders)
			if err != nil {
				return err
			}
			for _, v := range orders {
				m, err := m.layout.Repository().GetMenu(tx, v.MenuUuid)
				if err != nil {
					return err
				}
				listItems = append(listItems, listItem{order: &v, menu: m})
			}
			return nil
		})
		return listItems
	}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case []listItem:
		if msg == nil {
			cmds = append(cmds, m.list.StartSpinner())
		} else {
			m.list.StopSpinner()
			cmds = append(cmds, m.list.SetItems(makeListItems(msg)))
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
