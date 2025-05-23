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

// func (h *header) updateWorkingDir() {

// 	panic("here")
// 	if dir == h.workingDir {
// 		// nothing to do
// 		return
// 	}

// 	h.workingDir = dir
// 	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
// 	branch, err := cmd.Output()
// 	if err != nil {
// 		panic("failed to get branch")
// 		// just treat this directory as if it's not a git repo
// 		h.gitBranch = ""
// 		return
// 	}

// 	h.gitBranch = strings.Trim(string(branch), "\n")
// 	cmd = exec.Command("git", "status", "--porcelain")
// 	isClean, err := cmd.Output()
// 	if err != nil {
// 		panic("failed to check is clean")
// 		return
// 	}

// 	h.gitIsDirty = len(isClean) > 0
// }

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
