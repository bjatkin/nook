package model

import (
	"strings"
	"time"

	"github.com/bjatkin/nook/script/checker"
	"github.com/bjatkin/nook/script/normalizer"
	"github.com/bjatkin/nook/script/parser"
	"github.com/bjatkin/nook/script/vm"
	"github.com/bjatkin/nook/ui/colors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type runResult struct {
	result string
	errors []error
}

type activeCell struct {
	editor      codeEditor
	errorOutput []error
	vm          *vm.VM
	typeChecker *checker.Checker
	running     bool
	width       int
	height      int
}

func (a activeCell) Update(msg tea.Msg) (activeCell, tea.Cmd) {
	switch msg := msg.(type) {
	case resizeContent:
		a.width = msg.width
		a.height = msg.height

		editor, cmd := a.editor.Update(msg)
		a.editor = editor
		return a, cmd

	case runResult:
		a.errorOutput = msg.errors
		a.running = false
		if len(a.errorOutput) > 0 {
			return a, nil
		}

		code := a.editor.Text()
		a.editor = codeEditor{
			mode:    "INSERT",
			cursor:  cursor{row: 0, col: 1},
			content: []string{"("},
			width:   a.editor.width,
			height:  a.editor.height,
		}
		return a, func() tea.Msg {
			return addHistoryEntry{
				command: code,
				output:  msg.result,
			}
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			code := a.editor.Text()
			runCode := isCompleteExpression(code)
			if !runCode {
				editor, cmd := a.editor.Update(msg)
				a.editor = editor
				return a, cmd
			}

			// TODO: we should only run code if we're in normal mode
			if a.running {
				// if we're already running, don't run the code again
				return a, nil
			}

			a.running = true
			return a, func() tea.Msg {
				start := time.Now()
				result, errors := a.runCode([]byte(code))

				// minimum runtime is 1/16th second so the UI has time to update
				minDuration := time.Second / 16
				now := time.Now()
				duration := now.Sub(start)
				if duration < minDuration {
					minNano := minDuration.Nanoseconds()
					currentNano := duration.Nanoseconds()
					time.Sleep(time.Duration(minNano - currentNano))
				}

				return runResult{
					result: result,
					errors: errors,
				}
			}
		default:
			editor, cmd := a.editor.Update(msg)
			a.editor = editor
			return a, cmd
		}
	default:
		editor, cmd := a.editor.Update(msg)
		a.editor = editor
		return a, cmd
	}
}

func (a activeCell) runCode(code []byte) (string, []error) {
	p := parser.NewParser(code)
	ast := p.Parse()
	if len(p.Errors) != 0 {
		return "", p.Errors
	}

	n := normalizer.Normalizer{}
	ast = n.Normalize(ast)
	if len(n.Errors) > 0 {
		return "", n.Errors
	}

	_ = a.typeChecker.Infer(ast)
	if len(a.typeChecker.Errors) > 0 {
		typeCheckErrs := a.typeChecker.Errors
		a.typeChecker.Errors = []error{}
		return "", typeCheckErrs
	}

	result, err := a.vm.Eval(ast)
	if err != nil {
		return "", []error{err}
	}

	return result.String(), nil
}

func (a activeCell) View() string {
	editor := ""
	background := colors.Black
	if a.running {
		background = colors.Green2
	}

	editor = a.editor.View(background)
	if len(a.errorOutput) == 0 {
		return editor
	}

	errStyle := lipgloss.NewStyle().Background(background).Foreground(colors.Yellow3)
	errorLines := []string{}
	for _, err := range a.errorOutput {
		errLines := strings.Split(err.Error(), "\n")
		for _, line := range errLines {
			line = "  │ " + line
			pad := a.width - len(line)
			line = line + strings.Repeat(" ", pad)
			errorLines = append(errorLines, errStyle.Render(line))
		}
	}

	return editor + "\n" + strings.Join(errorLines, "\n")
}

func isCompleteExpression(code string) bool {
	depth := 0
	inComment := false
	exprCount := 0
	for _, char := range code {
		if char == '\n' {
			inComment = false
			continue
		}

		if char == '#' {
			inComment = true
			continue
		}

		if !inComment && char == '(' {
			depth++
		}
		if !inComment && char == ')' {
			depth--
		}
		if !inComment {
			exprCount++
		}
	}

	return depth <= 0 && exprCount > 0
}
