package orders

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/Markus-Schwer/ordaa/internal/boundary/tui/components"
	"github.com/Markus-Schwer/ordaa/internal/boundary/tui/components/menu"
	"github.com/Markus-Schwer/ordaa/internal/entity"
)

type delegateKeyMap struct {
	sel  key.Binding
	help key.Binding
}

type listItem struct {
	order *entity.Order
	menu  *entity.Menu
}

func (l listItem) Title() string { return fmt.Sprintf("[%s] %s", l.order.State, l.menu.Name) }

//	func (l listItem) Description() string {
//		return fmt.Sprintf("deadline: %s", l.order.OrderDeadline.String())
//	}
func (l listItem) Description() string {
	return "description"

}
func (l listItem) FilterValue() string {
	return strings.Join([]string{l.order.State, l.menu.Name}, " ")
}

func makeListItems(i []listItem) []list.Item {
	var items []list.Item
	for _, item := range i {
		items = append(items, listItem(item))
	}
	return items
}

// TODO: all nils in UpdateFunc should be replaced with proper error handling
func newItemDelegate(ctx context.Context, repo entity.Repository) (d list.DefaultDelegate) {
	d = list.NewDefaultDelegate()
	keys := newDelegateKeyMap()
	d.UpdateFunc = func(msg tea.Msg, m *list.Model) (cmd tea.Cmd) {
		i, ok := m.SelectedItem().(listItem)
		if !ok {
			return nil
		}
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.sel):
				return func() tea.Msg {
					return components.LayoutEvent{
						Component: menu.MenuSelectorComponent,
						Uuid:      i.menu.Uuid,
					}
				}
			}
		}
		return nil
	}
	d.ShortHelpFunc = keys.ShortHelp
	d.FullHelpFunc = keys.FullHelp
	return
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		sel: key.NewBinding(
			key.WithKeys("+", "enter"),
			key.WithHelp("+ | enter", "select order"),
		),
		help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "show help"),
		),
	}
}

func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.help,
	}
}

func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.sel,
		},
	}
}
