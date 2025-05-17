package model

import (
	"github.com/bjatkin/nook/ui/colors"
	"github.com/bjatkin/nook/ui/layout"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type footer struct {
	mode  string // should this be an enum?
	width int
}

func (f footer) Update(msg tea.Msg) footer {
	return f
}

func (f footer) View() string {
	if f.width == 0 {
		return ""
	}

	mode := layout.Text{
		Text:  layout.Pad(3, 3, f.mode),
		Style: lipgloss.NewStyle().Background(colors.Blue3).Foreground(colors.Blue1),
	}
	left := layout.NewHContainer(f.width/2, layout.LeftToRight, lipgloss.NewStyle().Background(colors.Blue1).Foreground(colors.Blue1))
	left.Content = append(left.Content, mode)

	version := layout.Text{
		Text:  " v0.0.1 ",
		Style: lipgloss.NewStyle().Background(colors.Green3).Foreground(colors.Green1),
	}
	right := layout.NewHContainer((f.width/2)-1, layout.RightToLeft, lipgloss.NewStyle().Background(colors.Blue1).Foreground(colors.Blue1))
	right.Content = append(right.Content, version)

	// TODO: seems like weird style issues pop up when you use the full terminal width
	// not sure if this is a bubble tea issue or if it's a terminal quirk
	// either way I need some better helper functions to calculate widths and such.
	// I should use the layout package for that. Not sure how I'm gonna get the terminal
	// width from within the layout package but I'll figure it out.
	cont := layout.NewHContainer(f.width-1, layout.LeftToRight, colors.Debug1)
	cont.Content = append(cont.Content, left, right)
	return cont.String()
}
