package model

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bjatkin/nook/ui/colors"
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

func newHeader() (header, error) {
	dir, err := os.Getwd()
	if err != nil {
		return header{}, fmt.Errorf("failed to get working directory: %w", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return header{}, fmt.Errorf("failed to get users home directory: %w", err)
	}

	branch, isDirty := getGitBranch(dir)
	return header{
		homeDir:    home,
		workingDir: dir,
		gitBranch:  branch,
		gitIsDirty: isDirty,
	}, nil
}

/*
func (h header) Shape() (int, int) {
	return h.width, 1
}

func (h header) Render(width, height int) []string {
	if width == 0 || height == 0 {
		return nil
	}

	dir := strings.TrimPrefix(h.workingDir, h.homeDir+"/")
	dir = layout.Pad(3, 3, dir)
	dirStyle := lipgloss.NewStyle().Background(colors.Blue3).Foreground(colors.Blue1).Bold(true)

	gitStyle := lipgloss.NewStyle().Background(colors.Blue3).Foreground(colors.Blue1).Bold(true)
	gitBranch := h.gitBranch
	if gitBranch != "" {
		gitBranch = layout.Pad(1, 1, h.gitBranch)
		gitStyle = lipgloss.NewStyle().Background(colors.Green3).Foreground(colors.Green1)
	}

	if h.gitIsDirty {
		gitStyle = lipgloss.NewStyle().Background(colors.Yellow3).Foreground(colors.Yellow1)
	}

	backgroundStyle := lipgloss.NewStyle().Background(colors.Blue1)
	statusDiv := &layout.Div_{
		Direction: layout.LeftToRight_,
		Contents: []layout.Content_{
			layout.NewText(dir, dirStyle),
			layout.NewText(gitBranch, gitStyle),
			layout.NewText(strings.Repeat(" ", h.width), backgroundStyle),
		},
		Width:  width,
		Height: 1,
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

func (h header) Update(msg tea.Msg) (header, tea.Cmd) {
	switch msg := msg.(type) {
	case resizeContent:
		h.width = msg.width
	case runResult:
		dir, err := os.Getwd()
		if err != nil {
			return h, nil
		}

		if h.workingDir == dir {
			return h, nil
		}
		h.workingDir = dir

		branch, isDirty := getGitBranch(h.workingDir)
		h.gitBranch = branch
		h.gitIsDirty = isDirty
		return h, nil
	}

	return h, nil
}

func getGitBranch(dir string) (string, bool) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branch, err := cmd.Output()
	if err != nil {
		return "", true
	}

	gitBranch := strings.Trim(string(branch), "\n")
	cmd = exec.Command("git", "status", "--porcelain")
	isClean, err := cmd.Output()
	if err != nil {
		return gitBranch, false
	}

	gitIsDirty := len(isClean) > 0
	return gitBranch, gitIsDirty
}

func (h header) View() string {
	if h.width == 0 {
		return ""
	}

	dir := strings.TrimPrefix(h.workingDir, h.homeDir+"/")
	dir = "   " + dir + "   "
	dirStyle := lipgloss.NewStyle().Background(colors.Blue3).Foreground(colors.Blue1).Bold(true)

	pad := strings.Repeat(" ", h.width-len(dir))
	padStyle := lipgloss.NewStyle().Background(colors.Blue1)
	if h.gitBranch == "" {
		return dirStyle.Render(dir) + padStyle.Render(pad)
	}

	gitStyle := lipgloss.NewStyle().Background(colors.Green3).Foreground(colors.Green1)
	if h.gitIsDirty {
		gitStyle = lipgloss.NewStyle().Background(colors.Yellow3).Foreground(colors.Yellow1)
	}

	gitBranch := " " + h.gitBranch + " "
	pad = strings.Repeat(" ", h.width-len(dir)-len(gitBranch))
	return dirStyle.Render(dir) + gitStyle.Render(gitBranch) + padStyle.Render(pad)
}
