package model

import (
	"fmt"
	"strings"

	"github.com/bjatkin/nook/script/parser"
	"github.com/bjatkin/nook/script/vm"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type entry struct {
	content string
	output  string
	err     error
}

type editor struct {
	history          []entry
	copyHistoryIndex int
	content          string
	indent           int
	width            uint
	vm               vm.VM
}

func (e editor) Update(msg tea.Msg) (editor, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "left", "right":
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
		case "backspace":
			if len(e.content) == 0 {
				return e, nil
			}

			// deleting a closing paren
			if e.content[len(e.content)-1] == ')' {
				e.indent++
			}

			// deleteing an opening paren
			if e.content[len(e.content)-1] == '(' {
				e.indent--
			}

			e.content = e.content[:len(e.content)-1]
			return e, nil
		case "tab":
			e.content += "\t"
			return e, nil
		case "enter":
			if e.indent <= 0 {
				p := parser.NewParser([]byte(e.content))
				ast := p.Parse()
				result, err := e.vm.Eval(ast)

				// reset the content
				e.history = append(e.history, entry{
					content: e.content,
					output:  fmt.Sprint(result),
					err:     err,
				})

				e.content = "("
				e.indent = 1
				e.copyHistoryIndex = len(e.history)
				return e, nil
			}

			e.content += "\n" + strings.Repeat("\t", e.indent)
			return e, nil
		case "(":
			e.content += "("
			e.indent += 1
			return e, nil
		case ")":
			e.content += ")"
			e.indent -= 1
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
	view += renderContent(e.content) + cursor.Render("█")

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
