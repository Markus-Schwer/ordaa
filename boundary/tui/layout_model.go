package tui

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
)

const (
	MENUS = iota
	ORDERS

	TABS
	BODY
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
	activeBox int
	subModels map[int]tea.Model
	helpModel tea.Model
}

func NewLayoutModel(
	ctx context.Context,
	renderer *lipgloss.Renderer,
	pty ssh.Pty,
	repo entity.Repository,
) *LayoutModel {
	info := NewLayoutInfo(ctx, renderer, pty, repo)
	subModels := make(map[int]tea.Model)
	subModels[MENUS] = NewMenuModel(info)
	return &LayoutModel{
		LayoutInfo: info,
		activeTab:  MENUS,
		activeBox:  TABS,
		subModels:  subModels,
		helpModel:  NewHelpModel(info),
	}
}

func (m *LayoutModel) Init() tea.Cmd {
	m.activeTab = MENUS
	for subMdlIdx, subModel := range m.subModels {
		cmd := subModel.Init()
		m.updateSubmodel(subMdlIdx, cmd)
	}
	cmd := m.helpModel.Init()
	if cmd != nil {
		msg := cmd()
		m.helpModel, cmd = m.helpModel.Update(msg)
	}
	return nil
}

func (m *LayoutModel) Update(msg tea.Msg) (mdl tea.Model, cmd tea.Cmd) {
	m.helpModel.Update(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+t":
			m.activeBox = TABS
		case msg.String() == "ctrl+b":
			m.activeBox = BODY
		case key.Matches(msg, DefaultKeyMap.Quit):
			cmd = tea.Quit
		case msg.String() == "m":
			if m.activeBox == TABS {
				m.activeTab = MENUS
			}
		case msg.String() == "o":
			if m.activeBox == TABS {
				m.activeTab = ORDERS
			}
		default:
			m.subModels[m.activeBox], cmd = m.subModels[m.activeTab].Update(msg)
		}
	default:
		m.subModels[m.activeBox], cmd = m.subModels[m.activeTab].Update(msg)
	}
	mdl = m
	return
}

func (m *LayoutModel) renderActiveBold(tab int) (rendered string) {
	s := m.txtStyle.Padding(0, 2).Foreground(lipgloss.Color("12")).Background(lipgloss.Color("0"))
	if tab == m.activeTab {
		s = s.Foreground(lipgloss.Color("0")).Background(lipgloss.Color("12"))
	}
	switch tab {
	case ORDERS:
		rendered = s.Render("[O]RDERS")
	case MENUS:
		rendered = s.Render("[M]ENUS")
	}
	return
}

func (m *LayoutModel) View() (content string) {
	header := m.txtStyle.Width(m.width).Align(lipgloss.Center).Render("DOTINDER") + "\n\n"
	entries := make([]string, len(NAVBAR))
	for i, v := range NAVBAR {
		entries[i] = m.renderActiveBold(v)
	}
	borderStyle := WithBorderAndCorner(m.txtStyle, "t", m.activeBox == TABS)
	header += borderStyle.Render(lipgloss.JoinHorizontal(lipgloss.Bottom, entries...))
	header += "\n\n"
	m.headerOffset = lipgloss.Height(header)
	content += header
	content += WithBorderAndCorner(m.txtStyle, "b", m.activeBox == BODY).Render(m.subModels[m.activeTab].View())
	footer := "\n\n" + m.helpModel.View()
	m.footerOffset = lipgloss.Height(footer)
	content += footer
	return
}

func (m *LayoutModel) updateSubmodel(modelIdx int, cmd tea.Cmd) {
	msg := cmd()
	m.subModels[modelIdx], cmd = m.subModels[modelIdx].Update(msg)
	if cmd != nil {
		m.updateSubmodel(modelIdx, cmd)
	}
}
