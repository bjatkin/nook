package model

import (
	"fmt"
	"os"

	"github.com/bjatkin/nook/script/vm"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	tWidth, tHeight uint
	header          header
	editor          editor
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

	return Model{
		header: header{
			workingDir: pwd,
			homeDir:    home,
		},
		editor: editor{
			content: "(",
			indent:  1,
			vm:      vm.NewVM(),
		},
	}, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.tHeight = uint(size.Height)
		m.tWidth = uint(size.Width)
		m.header.width = m.tWidth
		m.editor.width = m.tWidth
		return m, nil
	}

	editor, cmd := m.editor.Update(msg)
	m.editor = editor

	return m, cmd
}

func (m Model) View() string {
	header := m.header.View()
	editor := m.editor.View()
	return header + "\n" + editor
}
