package model

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/bjatkin/nook/colors"
	"github.com/bjatkin/nook/layout"
	tea "github.com/charmbracelet/bubbletea"
)

type header struct {
	homeDir    string
	workingDir string
	gitBranch  string
	gitIsDirty bool
	width      int
}

func (h header) Update(msg tea.Msg) (header, tea.Cmd) {
	if changeDir, ok := msg.(changeDirMsg); ok {
		h.workingDir = string(changeDir)
		cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		cmd.Dir = string(changeDir)
		branch, err := cmd.Output()
		if err != nil {
			// just treat this directory as if it's not a git repo
			h.gitBranch = ""
			return h, nil
		}

		h.gitBranch = strings.Trim(string(branch), "\n")
		cmd = exec.Command("git", "status", "--porcelain")
		cmd.Dir = string(changeDir)
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

	cont := layout.NewHContainer(h.width-1, layout.LeftToRight, colors.Background1)

	dir := strings.TrimPrefix(h.workingDir, h.homeDir+"/")
	dir = layout.Pad(3, 3, dir)

	cont.Content = append(
		cont.Content,
		layout.Text{
			Text:  dir,
			Style: colors.Primary,
		},
	)

	if h.gitBranch == "" {
		return cont.String()
	}

	style := colors.Secondary
	if h.gitIsDirty {
		style = colors.Third
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
