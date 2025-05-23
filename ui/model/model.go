package model

import (
	"fmt"
	"strings"

	"github.com/bjatkin/nook/script/checker"
	"github.com/bjatkin/nook/script/vm"
	"github.com/bjatkin/nook/ui/colors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	header, err := newHeader()
	if err != nil {
		return Model{}, fmt.Errorf("failed to build header: %w", err)
	}

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
		pad := buildPad(view, m.tWidth, m.tHeight)
		return view + pad + "\n" + footer
	case showHistory:
		return header + "\n" + editor + "\n" + history
	case m.showFooter:
		view := header + "\n" + editor + "\n"
		pad := buildPad(view, m.tWidth, m.tHeight)
		return view + pad + "\n" + footer
	default:
		return header + "\n" + editor
	}
}

func buildPad(view string, width, height int) string {
	lines := len(strings.Split(view, "\n"))
	if lines > height {
		return ""
	}

	padLine := strings.Repeat(" ", width)
	paddingLines := []string{}
	for i := 0; i < height-lines; i++ {
		paddingLines = append(paddingLines, padLine)
	}

	padding := strings.Join(paddingLines, "\n")
	return lipgloss.NewStyle().Background(colors.Gray1).Render(padding)
}
