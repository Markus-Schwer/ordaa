package menu

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog/log"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
)

type delegateKeyMap struct {
	increment key.Binding
	decrement key.Binding
	help      key.Binding
}

type listItem entity.MenuItem

func (l listItem) Title() string       { return fmt.Sprintf("[%s] %s", l.ShortName, l.Name) }
func (l listItem) Description() string { return fmt.Sprintf("%4.2f€", float64(l.Price)/100.0) }
func (l listItem) FilterValue() string { return strings.Join([]string{l.Name, l.ShortName}, " ") }

func makeListItems(menu *entity.Menu) []list.Item {
	var items []list.Item
	for _, item := range menu.Items {
		items = append(items, listItem(item))
	}
	return items
}

// TODO: all nils in UpdateFunc should be replaced with proper error handling
func newItemDelegate(_ context.Context, _ entity.Repository) (d list.DefaultDelegate) {
	d = list.NewDefaultDelegate()
	keys := newDelegateKeyMap()
	d.UpdateFunc = func(msg tea.Msg, m *list.Model) (cmd tea.Cmd) {
		i, ok := m.SelectedItem().(listItem)
		if !ok {
			return nil
		}
		// orderUuid, err := components.GetFromContext(ctx, components.ORDER_KEY)
		// if err != nil {
		// 	return nil
		// }
		// userUuid, err := components.GetFromContext(ctx, components.USER_KEY)
		// if err != nil {
		// 	return nil
		// }
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.increment):
				// return func() tea.Msg {
				// 	// // TODO: add order item to db and update status message
				// 	// repo.Transaction(func(tx *gorm.DB) (err error) {
				// 	// 	_, err = repo.CreateOrderItem(tx, orderUuid, &entity.OrderItem{
				// 	// 		Paid: false,
				// 	// 		OrderUuid: orderUuid,
				// 	// 		Price: i.Price,
				// 	// 		User: userUuid,
				// 	// 		MenuItemUuid: i.Uuid,
				// 	// 	})
				// 	// 	return
				// 	// })
				// 	log.Info().Msgf("increment")
				// }
				return m.NewStatusMessage(statusMessageStyle(fmt.Sprintf("Added '%s' to order", i.Name)))
			case key.Matches(msg, keys.decrement):
				return func() tea.Msg {
					log.Info().Msgf("decrement")
					// TODO: remove order item from db and update status message
					return m.NewStatusMessage(statusMessageStyle(fmt.Sprintf("TODO: remove '%s' from order", i.Name)))
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
		increment: key.NewBinding(
			key.WithKeys("+", "enter", "ctrl+a"),
			key.WithHelp("+ | enter | crrl+a", "increment selection"),
		),
		decrement: key.NewBinding(
			key.WithKeys("-", "ctrl+x"),
			key.WithHelp("- | ctrl+x", "decrement selection"),
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
			d.increment,
			d.decrement,
		},
	}
}
