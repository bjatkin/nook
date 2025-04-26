package model

import (
	"fmt"
	"os"

	"github.com/bjatkin/nook/script/vm"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	tWidth, tHeight int
	header          header
	editor          editor
	footer          footer
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
	header = header.Update(changeDirMsg(pwd))

	return Model{
		header: header,
		editor: editor{
			content: "(",
			vm:      vm.NewVM(pwd),
		},
		footer: footer{
			mode: "NORMAL",
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

		m.header.width = m.tWidth
		m.editor.width = m.tWidth
		m.editor.height = m.tHeight - 2
		m.footer.width = m.tWidth

		return m, nil
	}

	header := m.header.Update(msg)
	m.header = header
	editor, cmd := m.editor.Update(msg)
	m.editor = editor
	footer := m.footer.Update(msg)
	m.footer = footer

	return m, cmd
}

func (m Model) View() string {
	header := m.header.View()
	editor := m.editor.View()
	return header + "\n" + editor + "\n" + m.footer.View()
}
