package model

import (
	"strings"

	"github.com/bjatkin/nook/ui/colors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type footer struct {
	mode  string // should this be an enum? (ya, but it needs to be editor wide I think)
	width int
}

func (f footer) Update(msg tea.Msg) (footer, tea.Cmd) {
	switch msg := msg.(type) {
	case resizeContent:
		f.width = msg.width
	case changeMode:
		f.mode = string(msg)
	}

	return f, nil
}

/*
func (f footer) Shape() (int, int) {
	return f.width, 1
}

func (f footer) Render(width, height int) []string {
	// TODO: maybe we should gaurentee that this will never happen
	// since we can just skip this call if width or height is 0
	if width == 0 || height == 0 {
		return nil
	}

	modeStyle := lipgloss.NewStyle().Background(colors.Blue3).Foreground(colors.Blue1)
	mode := layout.NewText(
		layout.Pad(3, 3, f.mode),
		modeStyle,
	)

	padd := layout.NewText(strings.Repeat(" ", 100), lipgloss.NewStyle().Background(colors.Blue1))
	leftDiv := &layout.Div_{
		Direction: layout.LeftToRight_,
		Contents:  []layout.Content_{mode, padd},
		Width:     width / 2,
		Height:    1,
	}

	padd = layout.NewText(strings.Repeat(" ", 15), lipgloss.NewStyle().Background(colors.Green1))
	versionStyle := lipgloss.NewStyle().Background(colors.Green3).Foreground(colors.Green1)
	version := " v0.0.1 " // TODO: inject this with go releaser maybe?
	rightDiv := &layout.Div_{
		Direction: layout.RightToLeft_,
		Contents: []layout.Content_{
			layout.NewText(version, versionStyle),
			padd,
		},
		Width:  width / 2,
		Height: 1,
	}

	statusDiv := &layout.Div_{
		Direction: layout.LeftToRight_,
		Contents:  []layout.Content_{leftDiv, rightDiv},
		Width:     width,
		Height:    1,
	}

	content := layout.Div_{
		Direction: layout.TopToBottom_,
		Contents:  []layout.Content_{statusDiv},
		Width:     width,
		Height:    height,
	}

	return content.Render(width, height)
}
*/

func (f footer) View() string {
	if f.width == 0 {
		return ""
	}

	mode := "   " + f.mode + "   "
	modeStyle := lipgloss.NewStyle().Background(colors.Blue3).Foreground(colors.Blue1)

	version := " v0.0.1 "
	versionStyle := lipgloss.NewStyle().Background(colors.Green3).Foreground(colors.Green1)

	pad := strings.Repeat(" ", f.width-len(mode)-len(version))
	padStyle := lipgloss.NewStyle().Background(colors.Blue1)

	return modeStyle.Render(mode) + versionStyle.Render(version) + padStyle.Render(pad)
}
