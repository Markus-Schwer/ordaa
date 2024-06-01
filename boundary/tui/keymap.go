package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding
	Right      key.Binding
	Quit       key.Binding
	ToggleHelp key.Binding
}

var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "Move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "Move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("←/h", "Move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("l", "right"),
		key.WithHelp("→/l", "Move right"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c", "esc"),
		key.WithHelp("q/ctrl+c/esc", "Quit"),
	),
	ToggleHelp: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "Toggle help"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.ToggleHelp, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.ToggleHelp, k.Quit},
	}
}

type HelpModel struct {
	*LayoutInfo
	keyMap *KeyMap
	help   help.Model
}

func NewHelpModel(info *LayoutInfo) *HelpModel {
	h := help.New()
	h.Width = info.width
	h.ShowAll = false
	return &HelpModel{
		LayoutInfo: info,
		keyMap:     &DefaultKeyMap,
		help:       help.New(),
	}
}

func (m *HelpModel) Init() tea.Cmd {
	return nil
}

func (m *HelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.ToggleHelp):
			m.help.ShowAll = !m.help.ShowAll
		}
	}
	return m, nil
}
func (m *HelpModel) View() string {
	return m.help.View(m.keyMap)
}
