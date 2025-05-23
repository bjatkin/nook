package model

import (
	"slices"
	"strings"

	"github.com/bjatkin/nook/ui/colors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type addHistoryEntry struct {
	command string
	output  string
}

type historyEntry struct {
	command string
	output  string
}

type history struct {
	entries []historyEntry

	// TODO: add width and height to the editor and support scrolling if the
	// input get's too large
	width  int
	height int
}

func (h history) Update(msg tea.Msg) (history, tea.Cmd) {
	switch msg := msg.(type) {
	case resizeContent:
		h.width = msg.width
		h.height = msg.height
		return h, nil
	case addHistoryEntry:
		h.entries = append(h.entries, historyEntry{
			command: msg.command,
			output:  msg.output,
		})

		return h, nil
	default:
		return h, nil
	}
}

func (h history) View() string {
	if len(h.entries) == 0 {
		return ""
	}

	// TODO: cache this style on the history struct
	dividerStyle := lipgloss.NewStyle().Background(colors.Blue1)

	view := []string{}
	for _, entry := range slices.Backward(h.entries) {
		command := renderCommand(h.width, entry.command)
		output := renderOutput(h.width, entry.output)
		divider := dividerStyle.Render(strings.Repeat(" ", h.width))
		view = append(view, command+"\n"+output+"\n"+divider)
	}
	return strings.Join(view, "\n")
}

func renderCommand(width int, command string) string {
	// TODO: this probably won't change in the history so we can just cache it
	// on the history struct when it's created
	styles := styles(colors.Blue2)
	view := []string{}
	for _, line := range strings.Split(command, "\n") {
		pad := width - len(line)
		padding := styles["default"].Render(strings.Repeat(" ", pad))
		view = append(view, renderLine(line, styles)+padding)
	}

	return strings.Join(view, "\n")
}

func renderOutput(width int, output string) string {
	view := []string{}

	// TODO: cache these styles
	gutterStyle := lipgloss.NewStyle().Background(colors.Blue1).Foreground(colors.Blue3)
	lineStyle := lipgloss.NewStyle().Background(colors.Blue1).Foreground(colors.White)

	for _, line := range strings.Split(output, "\n") {
		pad := width - len(line) + 4
		padding := lineStyle.Render(strings.Repeat(" ", pad))
		view = append(view, gutterStyle.Render("  | ")+lineStyle.Render(line)+padding)
	}
	return strings.Join(view, "\n")
}
