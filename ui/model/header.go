package model

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bjatkin/nook/ui/colors"
	"github.com/bjatkin/nook/ui/layout"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type header struct {
	homeDir    string
	workingDir string
	gitBranch  string
	gitIsDirty bool
	width      int
}

func (h header) Update(msg tea.Msg) (header, tea.Cmd) {
	if msg, ok := msg.(resizeContent); ok {
		h.width = msg.width
		return h, nil
	}

	if _, ok := msg.(runResult); ok {
		dir, err := os.Getwd()
		if err != nil {
			return h, nil
		}

		if dir == h.workingDir {
			return h, nil
		}

		h.workingDir = dir
		cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		branch, err := cmd.Output()
		if err != nil {
			// just treat this directory as if it's not a git repo
			h.gitBranch = ""
			return h, nil
		}

		h.gitBranch = strings.Trim(string(branch), "\n")
		cmd = exec.Command("git", "status", "--porcelain")
		isClean, err := cmd.Output()
		if err != nil {
			return h, nil
		}

		h.gitIsDirty = len(isClean) > 0
		return h, func() tea.Msg {
			return debugInfoMsg(fmt.Sprintf("git status '%s', '%v'", h.gitBranch, string(isClean)))
		}
	}

	return h, nil
}

func (h header) View() string {
	if h.width == 0 {
		return ""
	}

	// cont := layout.NewHContainer(h.width-1, layout.LeftToRight, colors.Background1)
	cont := layout.NewHContainer(h.width, layout.LeftToRight, lipgloss.NewStyle().Background(colors.Blue1).Foreground(colors.Blue1))

	dir := strings.TrimPrefix(h.workingDir, h.homeDir+"/")
	dir = layout.Pad(3, 3, dir)

	cont.Content = append(
		cont.Content,
		layout.Text{
			Text:  dir,
			Style: lipgloss.NewStyle().Background(colors.Blue3).Foreground(colors.Blue1).Bold(true),
		},
	)

	if h.gitBranch == "" {
		return cont.String()
	}

	style := lipgloss.NewStyle().Background(colors.Green3).Foreground(colors.Green1)
	if h.gitIsDirty {
		style = lipgloss.NewStyle().Background(colors.Yellow3).Foreground(colors.Yellow1)
	}

	cont.Content = append(
		cont.Content,
		layout.Text{
			Text:  layout.Pad(1, 1, h.gitBranch),
			Style: style,
		},
	)

	return cont.String()
}
