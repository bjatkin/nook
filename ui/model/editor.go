package model

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bjatkin/nook/script/checker"
	"github.com/bjatkin/nook/script/normalizer"
	"github.com/bjatkin/nook/script/parser"
	"github.com/bjatkin/nook/script/token"
	"github.com/bjatkin/nook/script/vm"
	"github.com/bjatkin/nook/ui/colors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type changeDirMsg string
type debugInfoMsg string
type errorPingMsg struct{}

type entry struct {
	content string
	output  string
	err     error
}

type editor struct {
	debug            []string
	showDebug        bool
	history          []entry
	workingDir       string
	copyHistoryIndex int
	content          string
	contentHilight   bool
	output           string
	width            int
	height           int
	vm               vm.VM
	typeChecker      *checker.Checker
}

func (e editor) Update(msg tea.Msg) (editor, tea.Cmd) {
	e.debug = []string{}
	// TODO: remove me after testing is done
	indent, nesting := getNesting(e.content)
	e.debug = append(e.debug, fmt.Sprintf("indent '%d' nesting '%d'", indent, nesting))

	switch msg := msg.(type) {
	case resizeContent:
		e.width = msg.width
		e.height = msg.height
		return e, nil
	case debugInfoMsg:
		e.debug = append(e.debug, fmt.Sprintf("debug msg: '%s'", string(msg)))
	case errorPingMsg:
		e.contentHilight = false
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "esc", "left", "right", "f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "f10", "f11", "f12":
			// capture these events so we don't print them to the editor
		case "up":
			if e.copyHistoryIndex > 0 {
				e.copyHistoryIndex -= 1
			}

			i := e.copyHistoryIndex
			e.content = e.history[i].content
		case "down":
			if e.copyHistoryIndex < len(e.history) {
				e.copyHistoryIndex += 1
			}

			i := e.copyHistoryIndex
			if i == len(e.history) {
				e.content = "("
			} else {
				e.content = e.history[i].content
			}
		case "ctrl+c":
			return e, tea.Quit
		case "ctrl+d":
			e.showDebug = !e.showDebug
		case "backspace":
			if len(e.content) == 0 {
				return e, nil
			}

			e.content = e.content[:len(e.content)-1]
			return e, nil
		case "tab":
			e.content += "\t"
			return e, nil
		case "enter":
			if e.content == "(exit)" || e.content == "exit" {
				return e, tea.Quit
			}

			indent, nesting := getNesting(e.content)
			if nesting <= 0 && !isJustComment(e.content) {
				p := parser.NewParser([]byte(e.content))
				ast := p.Parse()
				if len(p.Errors) != 0 {
					errs := []string{}
					for _, e := range p.Errors {
						errs = append(errs, e.Error())
					}
					e.output = strings.Join(errs, "\n")
					e.contentHilight = true
					return e, func() tea.Msg { time.Sleep(time.Second / 8); return errorPingMsg{} }
				}

				e.debug = append(e.debug, fmt.Sprintf("ast: %#v", ast))

				n := normalizer.Normalizer{}
				ast = n.Normalize(ast)
				if len(n.Errors) > 0 {
					errs := []string{}
					for _, err := range n.Errors {
						errs = append(errs, err.Error())
					}
					e.output = strings.Join(errs, "\n")
					e.contentHilight = true
					return e, func() tea.Msg { time.Sleep(time.Second / 8); return errorPingMsg{} }
				}

				e.debug = append(e.debug, fmt.Sprintf("norm: %#v", ast))

				_ = e.typeChecker.Infer(ast)
				if len(e.typeChecker.Errors) > 0 {
					errs := []string{}
					for _, err := range e.typeChecker.Errors {
						errs = append(errs, err.Error())
					}
					e.output = strings.Join(errs, "\n")

					// clear the errors out between runs
					e.typeChecker.Errors = []error{}
					e.contentHilight = true
					return e, func() tea.Msg { time.Sleep(time.Second / 8); return errorPingMsg{} }
				}

				e.debug = append(e.debug, fmt.Sprintf("checked: %#v", ast))

				// TODO: eval actually needs to happen in a command in case it's a long runing operation
				e.output = ""
				result, err := e.vm.Eval(ast)

				// reset the content
				e.history = append(e.history, entry{
					content: e.content,
					output:  result.String(),
					err:     err,
				})

				e.content = "("
				e.copyHistoryIndex = len(e.history)
				dir, _ := os.Getwd()
				e.debug = append(e.debug, fmt.Sprint("pwd: ", dir))

				// update the working dir in case the VM updates the working dir
				pwd, err := os.Getwd()
				if err != nil {
					pwd = "???"
				}

				if e.workingDir != pwd {
					e.workingDir = pwd
					msg := changeDirMsg(e.workingDir)
					e.debug = append(e.debug, fmt.Sprintf("msg: %s", e.workingDir))
					return e, func() tea.Msg { return msg }
				}

				return e, nil
			}

			e.content += "\n"
			if indent > 0 {
				e.content += strings.Repeat("\t", indent)
			}
			return e, nil
		default:
			e.content += key
			return e, nil
		}
	}

	return e, nil
}

var errorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#c42132"))

func (e editor) View() string {
	if e.height == 0 {
		return ""
	}

	view := []string{}

	outputStyle := lipgloss.NewStyle().Foreground(colors.Blue3).Background(colors.Gray1)
	for _, entry := range e.history {
		view = append(view, renderContent(false, false, entry.content)...)
		if entry.err != nil {
			lines := strings.Split(entry.err.Error(), "\n")
			for _, line := range lines {
				view = append(view, lipgloss.NewStyle().Foreground(colors.Gray2).Render("  │ ")+errorStyle.Render(line))
			}
		} else {
			lines := strings.Split(entry.output, "\n")
			for _, line := range lines {
				view = append(view, lipgloss.NewStyle().Foreground(colors.Gray2).Render("  │ ")+outputStyle.Render(line))
			}
		}
		divider := lipgloss.NewStyle().Background(colors.Blue1)
		view = append(view, divider.Render(strings.Repeat(" ", int(e.width))))
	}

	view = append(view, renderContent(e.contentHilight, true, e.content)...)
	if e.output != "" {
		outputErrStyle := lipgloss.NewStyle().Foreground(colors.Yellow3).Background(colors.Gray1)
		lines := strings.Split(e.output, "\n")
		for _, line := range lines {
			view = append(view, "  │ "+outputErrStyle.Render(line))
		}
	}

	if e.showDebug {
		view = append(view, "")
		view = append(view, "")
		view = append(view, e.debug...)
	}

	for len(view) < int(e.height) {
		view = append(view, "")
	}

	start := len(view) - int(e.height)

	return strings.Join(view[start:], "\n")
}

func renderContent(hilight bool, includeCursor bool, content string) []string {
	view := ""
	tabString := "····"
	spaceString := "·"
	backgroundColor := colors.Gray1
	if hilight {
		backgroundColor = colors.Red3
	}

	l := parser.NewVerboseLexer([]byte(content))
	for _, t := range l.Lex() {
		switch t.Kind {
		case token.OpenParen, token.CloseParen:
			view += lipgloss.NewStyle().Foreground(colors.Gray2).Background(backgroundColor).Render(t.Value)
		case token.Whitespace:
			value := strings.ReplaceAll(t.Value, " ", spaceString)
			value = strings.ReplaceAll(value, "\t", tabString)
			view += lipgloss.NewStyle().Foreground(colors.Gray2).Background(backgroundColor).Render(value)
		case token.Let, token.Bool, token.GreaterThan, token.GreaterEqual,
			token.LessThan, token.LessEqual, token.Equal, token.Command:
			view += lipgloss.NewStyle().Foreground(colors.Blue3).Background(backgroundColor).Bold(true).Render(t.Value)
		case token.Plus, token.Minus, token.Divide, token.Multiply:
			view += lipgloss.NewStyle().Foreground(colors.Red3).Background(backgroundColor).Bold(true).Render(t.Value)
		case token.Int, token.Float:
			view += lipgloss.NewStyle().Foreground(colors.Purple4).Background(backgroundColor).Render(t.Value)
		case token.String, token.Atom:
			view += lipgloss.NewStyle().Foreground(colors.Yellow4).Background(backgroundColor).Render(t.Value)
		case token.Comment:
			view += lipgloss.NewStyle().Foreground(colors.Green3).Background(backgroundColor).Render(t.Value)
		default:
			view += lipgloss.NewStyle().Foreground(colors.White).Background(backgroundColor).Render(t.Value)
		}
	}

	if includeCursor {
		cursor := lipgloss.NewStyle().Foreground(colors.Yellow4)
		view += cursor.Render("█")
	}

	return strings.Split(view, "\n")
}

func isJustComment(code string) bool {
	lines := strings.Split(code, "\n")
	lineCount := 0
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if !strings.HasPrefix(line, "#") {
			lineCount++
		}
	}

	return lineCount == 0
}

func getNesting(code string) (int, int) {
	indent := 0
	lines := strings.Split(code, "\n")
	nesting := 0

	for _, line := range lines {
		open := strings.Count(line, "(")
		close := strings.Count(line, ")")
		nest := open - close
		nesting += nest
		if nest > 0 {
			indent++
		}
		if nest < 0 {
			indent--
		}
	}

	return indent, nesting
}
