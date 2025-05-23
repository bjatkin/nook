package layout

import (
	"strings"

	"github.com/bjatkin/nook/script/parser"
	"github.com/bjatkin/nook/script/token"
	"github.com/bjatkin/nook/ui/colors"
	"github.com/charmbracelet/lipgloss"
)

type Code struct {
	lines  []string
	width  int
	height int
}

func NewCode(code string) Code {
	lines := strings.Split(code, "\n")
	height := len(lines)
	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}

	return Code{
		lines:  lines,
		width:  width,
		height: height,
	}
}

func (t Code) Shape() (int, int) {
	return t.width, t.height
}

func (t Code) Render(width, height int) []string {
	styles := styles(colors.Blue1)
	lines := []string{}
	for i := 0; i < height; i++ {
		if i >= len(t.lines) {
			lines = append(lines, "")
			continue
		}

		line := setWidth(t.lines[i], t.width)
		line = renderLine(line, styles)
		lines = append(lines, line)
	}

	return lines
}

func styles(background lipgloss.TerminalColor) map[string]lipgloss.Style {
	return map[string]lipgloss.Style{
		"muted":   lipgloss.NewStyle().Foreground(colors.Gray2).Background(background),
		"keyword": lipgloss.NewStyle().Foreground(colors.Blue3).Background(background).Bold(true),
		"symbol":  lipgloss.NewStyle().Foreground(colors.Red3).Background(background).Bold(true),
		"number":  lipgloss.NewStyle().Foreground(colors.Purple4).Background(background),
		"string":  lipgloss.NewStyle().Foreground(colors.Yellow4).Background(background),
		"comment": lipgloss.NewStyle().Foreground(colors.Green3).Background(background),
		"default": lipgloss.NewStyle().Foreground(colors.White).Background(background),
	}
}

func renderLine(line string, styles map[string]lipgloss.Style) string {
	lexer := parser.NewVerboseLexer([]byte(line))
	view := ""
	for _, tok := range lexer.Lex() {
		if tok.Value == "\n" {
			continue
		}

		view += styleToken(tok, styles)
	}

	return view
}

func styleToken(tok token.Token, styles map[string]lipgloss.Style) string {
	if len(tok.Value) == 0 {
		return ""
	}

	switch tok.Kind {
	case token.OpenParen, token.CloseParen:
		return styles["muted"].Render(tok.Value)
	case token.Whitespace:
		value := strings.ReplaceAll(tok.Value, " ", "·")
		value = strings.ReplaceAll(value, "\t", "*---")
		return styles["muted"].Render(value)
	case token.Let, token.Bool, token.GreaterThan, token.GreaterEqual,
		token.LessThan, token.LessEqual, token.Equal, token.Command:
		return styles["keyword"].Render(tok.Value)
	case token.Plus, token.Minus, token.Divide, token.Multiply:
		return styles["symbol"].Render(tok.Value)
	case token.Int, token.Float:
		return styles["number"].Render(tok.Value)
	case token.String:
		return styles["string"].Render(tok.Value)
	case token.Comment:
		return styles["comment"].Render(tok.Value)
	default:
		return styles["default"].Render(tok.Value)
	}
}
