package tui

import (
	"context"
	"math"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/rs/zerolog/log"
	"github.com/Markus-Schwer/ordaa/internal/boundary/tui/components"
	"github.com/Markus-Schwer/ordaa/internal/boundary/tui/components/menu"
	"github.com/Markus-Schwer/ordaa/internal/boundary/tui/components/orders"
	"github.com/Markus-Schwer/ordaa/internal/entity"
)

const (
	MENUS = iota
	ORDERS
	HOME

	TABS
	BODY
)

var NAVBAR = []int{HOME, ORDERS, MENUS}

type LayoutInfo struct {
	ctx          context.Context
	term         string
	profile      string
	bg           string
	txtStyle     lipgloss.Style
	quitStyle    lipgloss.Style
	width        int
	height       int
	headerOffset int
	footerOffset int
	repo         entity.Repository
}

func (l *LayoutInfo) BoxStyle() lipgloss.Style {
	return l.txtStyle
}

func (l *LayoutInfo) ContentWidth() int {
	return int(math.Min(float64(l.width), 60))
}

func (l *LayoutInfo) ContentHeight() int {
	return l.height - l.headerOffset - l.footerOffset - 10
}

func (l *LayoutInfo) Repository() entity.Repository {
	return l.repo
}

func NewLayoutInfo(
	ctx context.Context,
	renderer *lipgloss.Renderer,
	pty ssh.Pty,
	repo entity.Repository,
) *LayoutInfo {
	return &LayoutInfo{
		ctx:       ctx,
		term:      pty.Term,
		profile:   renderer.ColorProfile().String().String(),
		txtStyle:  renderer.NewStyle().Foreground(lipgloss.Color("10")),
		quitStyle: renderer.NewStyle().Foreground(lipgloss.Color("8")),
		width:     pty.Window.Width,
		height:    pty.Window.Height,
		repo:      repo,
	}
}

type LayoutModel struct {
	*LayoutInfo
	activeTab   components.ComponentKey
	activeBox   int
	activeModel tea.Model
	// subModels   map[components.ComponentKey]tea.Model
	helpModel tea.Model
}

func NewLayoutModel(
	ctx context.Context,
	renderer *lipgloss.Renderer,
	pty ssh.Pty,
	repo entity.Repository,
) *LayoutModel {
	info := NewLayoutInfo(ctx, renderer, pty, repo)
	// subModels := make(map[components.ComponentKey]tea.Model)
	// subModels[menu.MenuSelectorComponent] = menu.NewMenuItemSelectorModel(ctx, repo, info)
	// subModels[orders.OrderSelectorComponent] = orders.NewOrderSelectorModel(ctx, repo, info)
	return &LayoutModel{
		LayoutInfo: info,
		activeTab:  orders.OrderSelectorComponent,
		activeBox:  BODY,
		// subModels:  subModels,
		helpModel: NewHelpModel(info),
	}
}

func (m *LayoutModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	m.activeTab = orders.OrderSelectorComponent
	m.activeModel = orders.NewOrderSelectorModel(m.ctx, m.repo, m)
	cmds = append(cmds, m.activeModel.Init())
	cmds = append(cmds, m.helpModel.Init())
	return tea.Batch(cmds...)
}

func (m *LayoutModel) Update(msg tea.Msg) (mdl tea.Model, cmd tea.Cmd) {
	m.helpModel.Update(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case components.LayoutEvent:
		log.Ctx(m.ctx).Info().Msgf("LayoutEvent: %v", msg)
		m.activeTab = msg.Component
		if msg.Component == menu.MenuSelectorComponent {
			m.ctx = components.SetMenuUuidToContext(m.ctx, msg.Uuid)
			m.activeModel = menu.NewMenuItemSelectorModel(m.ctx, m.repo, m)
			cmd = m.activeModel.Init()
		}
	default:
		m.activeModel, cmd = m.activeModel.Update(msg)
	}
	mdl = m
	return
}

// func (m *LayoutModel) renderActiveBold(tab int) (rendered string) {
// 	s := m.txtStyle.Padding(0, 2).Foreground(lipgloss.Color("12")).Background(lipgloss.Color("0"))
// 	if tab == m.activeTab {
// 		s = s.Foreground(lipgloss.Color("0")).Background(lipgloss.Color("12"))
// 	}
// 	switch tab {
// 	case HOME:
// 		rendered = s.Render(padToSizeCenter("ordaa", 12))
// 	case ORDERS:
// 		rendered = s.Render(padToSizeCenter("ORDERS", 12))
// 	case MENUS:
// 		rendered = s.Render(padToSizeCenter("MENUS", 12))
// 	}
// 	return
// }

func padToSizeCenter(str string, total int) string {
	pad := int((total - len(str)) / 2)
	return strings.Repeat(" ", pad) + str + strings.Repeat(" ", total-len(str)-pad)
}

func (m *LayoutModel) View() (content string) {
	entries := make([]string, len(NAVBAR))
	// for i, v := range NAVBAR {
	// 	entries[i] = m.renderActiveBold(v)
	// }
	header := "\n" + m.txtStyle.
		Width(m.width).
		Align(lipgloss.Center).
		Render(
			lipgloss.JoinHorizontal(lipgloss.Center, entries...),
		) + "\n\n"
	m.headerOffset = lipgloss.Height(header)
	content += header
	footer := "\n\n" + m.helpModel.View()
	m.footerOffset = lipgloss.Height(footer)
	content += m.txtStyle.
		Width(m.width).
		Align(lipgloss.Center).
		Render(m.activeModel.View())
	content += footer
	return
}

// func (m *LayoutModel) updateSubmodel(modelIdx components.ComponentKey, cmd tea.Cmd) {
// 	msg := cmd()
// 	m.subModels[modelIdx], cmd = m.subModels[modelIdx].Update(msg)
// 	if cmd != nil {
// 		m.updateSubmodel(modelIdx, cmd)
// 	}
// }
