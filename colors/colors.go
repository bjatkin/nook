package colors

import "github.com/charmbracelet/lipgloss"

var Default = lipgloss.NewStyle().
	Background(lipgloss.Color("#4b3d44")).
	Foreground(lipgloss.Color("#d2c9a5"))

var Emphasis = lipgloss.NewStyle().
	Background(lipgloss.Color("#4b726e")).
	Foreground(lipgloss.Color("#d2c9a5")).
	Bold(true)

var Emphasis2 = lipgloss.NewStyle().
	Background(lipgloss.Color("#574852")).
	Foreground(lipgloss.Color("#ab9b8e"))

var Second = lipgloss.NewStyle().
	Background(lipgloss.Color("#927441")).
	Foreground(lipgloss.Color("#d2c9a5"))

var Error = lipgloss.NewStyle().
	Background(lipgloss.Color("#4b3d44")).
	Foreground(lipgloss.Color("#c77b58"))

var Debug1 = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FFFFFF")).
	Background(lipgloss.Color("#c90000"))

var Debug2 = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#000000")).
	Background(lipgloss.Color("#ffcb2e"))
