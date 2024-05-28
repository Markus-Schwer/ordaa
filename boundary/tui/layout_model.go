package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
)

const (
	MENUS = iota
	ORDERS
)

var NAVBAR = []int{MENUS, ORDERS}

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
	activeTab int
	menuModel tea.Model
}

func NewLayoutModel(
	ctx context.Context,
	renderer *lipgloss.Renderer,
	pty ssh.Pty,
	repo entity.Repository,
) *LayoutModel {
	info := NewLayoutInfo(ctx, renderer, pty, repo)
	return &LayoutModel{
		LayoutInfo: info,
		activeTab:  MENUS,
		menuModel:  NewMenuModel(info),
	}
}

func (m *LayoutModel) Init() tea.Cmd {
	m.activeTab = MENUS
	return nil
}

func (m *LayoutModel) Update(msg tea.Msg) (mdl tea.Model, cmd tea.Cmd) {
	switch m.activeTab {
	case MENUS:
		m.menuModel, cmd = m.menuModel.Update(msg)
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			cmd = tea.Quit
		case "tab":
			if m.activeTab == ORDERS {
				m.activeTab = MENUS
				m.menuModel.Init()
			} else {
				m.activeTab = m.activeTab + 1
			}
		}
	}
	mdl = m
	return
}

func (m *LayoutModel) renderActiveBold(tab int) (rendered string) {
	spacedText := m.txtStyle.
		PaddingLeft(2).
		PaddingRight(2).
		Bold(tab == m.activeTab)
	switch tab {
	case ORDERS:
		rendered = spacedText.
			Render("ORDERS")
	case MENUS:
		rendered = spacedText.
			Render("MENUS")
	}
	return
}

func (m *LayoutModel) View() (content string) {
	content = m.txtStyle.Width(m.width).Align(lipgloss.Center).Render("DOTINDER") + "\n\n"
	for _, v := range NAVBAR {
		content += m.renderActiveBold(v)
	}
	content += "\n\n"
	m.headerOffset = lipgloss.Height(content)
	// log.Ctx(m.ctx).Info().Msgf("header: %v", content)
	switch m.activeTab {
	case MENUS:
		content += m.menuModel.View()
	case ORDERS:
		content += "I AM THE ORDERS DUMMY"
	}
	return
}
