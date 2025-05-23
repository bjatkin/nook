package model

import (
	"fmt"
	"os"
	"strings"

	"github.com/bjatkin/nook/script/checker"
	"github.com/bjatkin/nook/script/vm"
	tea "github.com/charmbracelet/bubbletea"
)

type resizeContent struct {
	width  int
	height int
}

type Model struct {
	tWidth, tHeight int
	header          header
	editor          editor
	showFooter      bool
	footer          footer

	history    history
	activeCell activeCell
}

func NewModel() (Model, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return Model{}, fmt.Errorf("failed to get working directory: %w", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return Model{}, fmt.Errorf("failed to get users home directory: %w", err)
	}

	// TODO: this is not the right way to do this, clean this up
	header := header{
		workingDir: pwd,
		homeDir:    home,
	}
	header, _ = header.Update(changeDirMsg(pwd))

	return Model{
		header: header,
		editor: editor{
			content:     "(",
			vm:          *vm.NewVM(),
			typeChecker: checker.NewChecker(),
		},
		footer: footer{
			mode: "INSERT",
		},

		history: history{},
		activeCell: activeCell{
			editor: codeEditor{
				cursor:  cursor{row: 0, col: 1},
				mode:    "INSERT",
				content: []string{"("},
			},
			vm:          vm.NewVM(),
			typeChecker: checker.NewChecker(),
		},
	}, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.tHeight = size.Height
		m.tWidth = size.Width

		m.showFooter = false
		if m.tHeight > 10 {
			m.showFooter = true
		}

		msg = resizeContent{
			width:  size.Width,
			height: size.Height,
		}
	}

	test, testCmd := m.activeCell.Update(msg)
	m.activeCell = test

	history, historyCmd := m.history.Update(msg)
	m.history = history

	header, headerCmd := m.header.Update(msg)
	m.header = header

	footer, footerCmd := m.footer.Update(msg)
	m.footer = footer
	return m, tea.Batch(testCmd, historyCmd, headerCmd, footerCmd)
}

func (m Model) View() string {
	header := m.header.View()
	history := m.history.View()
	showHistory := history != ""

	editor := m.activeCell.View()
	footer := m.footer.View()

	switch {
	case showHistory && m.showFooter:
		view := header + "\n" + editor + "\n" + history + "\n"
		pad := buildPad(view, m.tHeight)
		return view + pad + footer
	case showHistory:
		return header + "\n" + editor + "\n" + history
	case m.showFooter:
		view := header + "\n" + editor + "\n"
		pad := buildPad(view, m.tHeight)
		return view + pad + footer
	default:
		return header + "\n" + editor
	}
}

func buildPad(view string, height int) string {
	lines := len(strings.Split(view, "\n"))
	if lines > height {
		return ""
	}

	return strings.Repeat("\n", height-lines)
}
