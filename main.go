package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bjatkin/nook/model"
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
