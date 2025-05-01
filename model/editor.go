package model

import (
	"fmt"
	"strings"

	"github.com/bjatkin/nook/script/parser"
	"github.com/bjatkin/nook/script/vm"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type changeDirMsg string
type debugInfoMsg string

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
	width            int
	height           int
	vm               vm.VM
}

func (e editor) Update(msg tea.Msg) (editor, tea.Cmd) {
	e.debug = []string{}
	// TODO: remove me after testing is done
	indent, nesting := getNesting(e.content)
	e.debug = append(e.debug, fmt.Sprintf("indent '%d' nesting '%d'", indent, nesting))

	switch msg := msg.(type) {
	case debugInfoMsg:
		e.debug = append(e.debug, fmt.Sprintf("debug msg: '%s'", string(msg)))
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "left", "right", "f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "f10", "f11", "f12":
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
			if nesting <= 0 {
				p := parser.NewParser([]byte(e.content))
				ast := p.Parse()
				e.debug = append(e.debug, fmt.Sprint("ast: ", ast))

				// TODO: eval actually needs to happen in a command in case it's a long runing operation
				result, err := e.vm.Eval(ast)

				// reset the content
				e.history = append(e.history, entry{
					content: e.content,
					output:  result.String(),
					err:     err,
				})

				e.content = "("
				e.copyHistoryIndex = len(e.history)
				e.debug = append(e.debug, fmt.Sprint("pwd: ", e.vm.WorkingDir()))

				if e.workingDir != e.vm.WorkingDir() {
					e.workingDir = e.vm.WorkingDir()
					msg := changeDirMsg(e.workingDir)
					e.debug = append(e.debug, fmt.Sprintf("msg: %s", e.workingDir))
					return e, func() tea.Msg { return msg }
				} else {
					return e, nil
				}
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

var divider = lipgloss.NewStyle().
	Background(lipgloss.Color("#222222"))

var paren = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#9E9E9E"))

var space = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#777777"))

var cursor = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#2eff85")).
	Blink(true)

var errorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#c42132"))

var outputStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#edbe21"))

func (e editor) View() string {
	view := ""

	for _, entry := range e.history {
		view += renderContent(entry.content)
		view += "\n"
		if entry.err != nil {
			view += "\t" + errorStyle.Render(entry.err.Error())
		} else {
			view += "\t" + outputStyle.Render(entry.output)
		}
		view += "\n"
		view += divider.Render(strings.Repeat(" ", int(e.width))) + "\n"
	}

	// TODO: use the script parser to get an ast and do syntax hilighting
	// TODO: make it possible to move the cursor around
	view += renderContent(e.content) + cursor.Render("█")

	if e.showDebug {
		view += "\n\n" + strings.Join(e.debug, "\n")
	}

	// TODO: use containers for this
	height := strings.Count(view, "\n")
	for height+1 < int(e.height) {
		view += "\n"
		height++
	}

	return view
}

func renderContent(content string) string {
	view := ""
	tabString := "····"
	spaceString := "·"

	for _, char := range content {
		switch {
		case char == '(' || char == ')':
			view += paren.Render(string(char))
		case char == ' ':
			view += space.Render(spaceString)
		case char == '\t':
			view += space.Render(tabString)
		default:
			view += string(char)
		}
	}

	return view
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
