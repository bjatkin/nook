package colors

import "github.com/charmbracelet/lipgloss"

var (
	lightBlue = lipgloss.Color("#8caba1")
	darkBlue  = lipgloss.Color("#4b726e")

	lightGreen  = lipgloss.Color("#b3a555")
	darkGreen   = lipgloss.Color("#77743b")
	pastelGreen = lipgloss.Color("#d2c9a5")

	darkPurple = lipgloss.Color("#4b3d44")
	purple     = lipgloss.Color("#574852")

	lightGray = lipgloss.Color("#ab9b8e")

	darkBrown   = lipgloss.Color("#4d4539")
	mediumBrown = lipgloss.Color("#927441")
	lightBrown  = lipgloss.Color("#ba9158")

	darkOrang    = lipgloss.Color("#79444a")
	mediumOrange = lipgloss.Color("#ae5d40")
	lightOrange  = lipgloss.Color("#c77b58")
)

var Default = lipgloss.NewStyle().
	Background(darkPurple).
	Foreground(pastelGreen)

var Primary = lipgloss.NewStyle().
	Background(darkBlue).
	Foreground(pastelGreen).
	Bold(true)

var Background1 = lipgloss.NewStyle().
	Background(purple).
	Foreground(lightGray)

var Secondary = lipgloss.NewStyle().
	Background(darkGreen).
	Foreground(pastelGreen)

var Third = lipgloss.NewStyle().
	Background(darkOrang).
	Foreground(pastelGreen)

var Error = lipgloss.NewStyle().
	Background(darkPurple).
	Foreground(lightOrange)

var Debug1 = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FFFFFF")).
	Background(lipgloss.Color("#c90000"))

var Debug2 = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#000000")).
	Background(lipgloss.Color("#ffcb2e"))
