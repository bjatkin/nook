package main

import (
	"fmt"
	"os"

	"github.com/bjatkin/nook/model"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model, err := model.NewModel()
	if err != nil {
		fmt.Println("err:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("err:", err)
		os.Exit(1)
	}
}
