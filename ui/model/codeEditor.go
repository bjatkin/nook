package model

import (
	"fmt"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/bjatkin/nook/script/parser"
	"github.com/bjatkin/nook/script/token"
	"github.com/bjatkin/nook/ui/colors"
)

type changeMode string

type cursor struct {
	row int
	col int
}

// codeEditor is an editor component for writing nook script
type codeEditor struct {
	// TODO: modes need to be shared more broadly across the editor and should probably be
	// "enum" variants rather then strings like this
	mode    string
	cursor  cursor
	content []string

	// TODO: add width and height to the editor and support scrolling if the
	// input get's too large
	width  int
	height int
}

// Text returns the content of the editor as a string
func (c *codeEditor) Text() string {
	return strings.Join(c.content, "\n")
}

// Update is the update function for the codeEditor
func (c codeEditor) Update(msg tea.Msg) (codeEditor, tea.Cmd) {
	if msg, ok := msg.(resizeContent); ok {
		c.width = msg.width
		c.height = msg.height
		return c, nil
	}

	switch c.mode {
	case "NORMAL":
		return c.normalUpdate(msg)
	case "INSERT":
		return c.insertUpdate(msg)
	default:
		panic(fmt.Sprintf("invalid mode: '%s'", c.mode))
	}
}

func (c codeEditor) insertUpdate(msg tea.Msg) (codeEditor, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		// TODO: this probably needs to be more robust
		if len(key) == 1 || key == "space" || key == "tab" {
			key = strings.ReplaceAll(key, "space", " ")
			key = strings.ReplaceAll(key, "tab", "\t")

			row := c.cursor.row
			line := c.content[row]
			c.content[row] = insertChar(c.cursor.col, line, key)
			c.cursor.col++
		}

		switch key {
		case "up":
			c.moveCursor(up)
		case "down":
			c.moveCursor(down)
		case "left":
			c.moveCursor(left)
		case "right":
			c.moveCursor(right)
		case "ctrl+c":
			return c, tea.Quit
		case "backspace":
			// we can't delete anything at the beginning of the content
			if c.cursor.row == 0 && c.cursor.col == 0 {
				return c, nil
			}

			row := c.cursor.row
			line := c.content[row]

			// the line is already empty so we need to remove the whole thing
			if len(line) == 0 {
				c.content = removeLine(c.cursor.row, c.content)
				c.cursor.row = max(0, c.cursor.row-1)
				c.cursor.col = len(c.content[c.cursor.row])
			} else {
				c.content[row] = removeChar(c.cursor.col-1, line)
				c.cursor.col = max(0, c.cursor.col-1)
			}

			return c, nil
		case "tab":
			c.content[c.cursor.row] += key
			return c, nil
		case "enter":
			if c.content[0] == "(exit" ||
				c.content[0] == "exit)" ||
				c.content[0] == "(exit)" ||
				c.content[0] == "exit" {
				return c, tea.Quit
			}

			line := c.content[c.cursor.row]
			lineStart := line[:c.cursor.col]
			lineEnd := line[c.cursor.col:]

			start := c.content[:c.cursor.row]
			end := c.content[c.cursor.row+1:]

			c.cursor.row++
			c.cursor.col = 0
			c.content = slices.Concat(start, []string{lineStart}, []string{lineEnd}, end)
			return c, nil
		case "esc":
			c.mode = "NORMAL"
			return c, func() tea.Msg { return changeMode("NORMAL") }
		}
	}

	return c, nil
}

func (c codeEditor) normalUpdate(msg tea.Msg) (codeEditor, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "up", "k":
			c.moveCursor(up)
		case "down", "j":
			c.moveCursor(down)
		case "left", "h":
			c.moveCursor(left)
		case "right", "l":
			c.moveCursor(right)
		case "ctrl+c":
			return c, tea.Quit
		case "i":
			c.mode = "INSERT"
			return c, func() tea.Msg { return changeMode("INSERT") }
		case "a":
			c.mode = "INSERT"
			lineLen := len(c.content[c.cursor.row])
			c.cursor.col = min(lineLen, c.cursor.col+1)
			return c, func() tea.Msg { return changeMode("INSERT") }
		case "A":
			c.mode = "INSERT"
			lineLen := len(c.content[c.cursor.row])
			c.cursor.col = lineLen
			return c, func() tea.Msg { return changeMode("INSERT") }
		case "o":
			c.cursor.row++
			c.cursor.col = 0

			start := c.content[:c.cursor.row]
			end := c.content[c.cursor.row:]

			c.mode = "INSERT"
			c.content = slices.Concat(start, []string{""}, end)
			return c, func() tea.Msg { return changeMode("INSERT") }
		case "O":
			start := c.content[:c.cursor.row]
			end := c.content[c.cursor.row:]

			c.cursor.row++
			c.cursor.col = 0

			c.mode = "INSERT"
			c.content = slices.Concat(start, []string{""}, end)
			return c, func() tea.Msg { return changeMode("INSERT") }
		case "w":
			// skip until the next NookScript token
			line := c.content[c.cursor.row]
			lexer := parser.NewVerboseLexer([]byte(line))
			tokens := lexer.Lex()

			// find the current token
			currentToken := -1
			end := 0
			for i, tok := range tokens {
				end += len(tok.Value)
				if end > c.cursor.col {
					currentToken = i
					break
				}
			}
			if currentToken == -1 {
				currentToken = len(tokens)
			}

			// set the currsor at the begining of the next token
			end = 0
			for i, tok := range tokens {
				start := end
				end += len(tok.Value)
				if tok.Kind == token.Whitespace {
					continue
				}

				if i > currentToken {
					c.cursor.col = start
					break
				}
			}
		case "b":
			// skip backwards to the previous NookScript token
			line := c.content[c.cursor.row]
			lexer := parser.NewVerboseLexer([]byte(line))
			tokens := lexer.Lex()

			// find the current token
			currentToken := -1
			end := 0
			for i, tok := range tokens {
				end += len(tok.Value)
				if end > c.cursor.col {
					currentToken = i
					break
				}
			}
			if currentToken == -1 {
				currentToken = len(tokens)
			}

			// set the currsor at the begining of the previous token
			end = 0
			for i, tok := range tokens {
				start := end
				end += len(tok.Value)
				if i > 0 && tok.Kind == token.Whitespace {
					continue
				}

				if i >= currentToken-2 {
					c.cursor.col = start
					break
				}
			}
		}
	}
	return c, nil
}

type direction int

const (
	up = direction(iota)
	down
	left
	right
)

func (c *codeEditor) moveCursor(direction direction) {
	switch direction {
	case up:
		if c.cursor.row-1 < 0 {
			return
		}

		nextLine := c.content[c.cursor.row]
		c.cursor.row--
		line := c.content[c.cursor.row]

		// convert the column to it's visual position on the previous line
		// and then convert that visual position back to the content position
		// for the next line. This maintains the visual position of the cursor
		// as we move across lines
		visualCol := toVisualColumn(c.cursor.col, nextLine)
		contentCol := toContentColumn(visualCol, line)

		c.cursor.col = contentCol
	case down:
		if c.cursor.row+1 >= len(c.content) {
			return
		}

		prevLine := c.content[c.cursor.row]
		c.cursor.row++
		line := c.content[c.cursor.row]

		// convert the column to it's visual position on the previous line
		// and then convert that visual position back to the content position
		// for the next line. This maintains the visual position of the cursor
		// as we move across lines
		visualColumn := toVisualColumn(c.cursor.col, prevLine)
		contentCol := toContentColumn(visualColumn, line)

		c.cursor.col = contentCol
	case left:
		c.cursor.col = max(0, c.cursor.col-1)
	case right:
		lineLen := len(c.content[c.cursor.row])
		c.cursor.col = min(lineLen, c.cursor.col+1)
	}
}

func toVisualColumn(contentCol int, line string) int {
	prefix := line[:contentCol]
	tabs := strings.Count(prefix, "\t")
	shortPrefix := len(prefix) - tabs
	return shortPrefix + tabs*4
}

func toContentColumn(visualCol int, line string) int {
	contentCol := 0
	trackCol := 0
	for _, c := range line {
		contentCol++
		if c == '\t' {
			trackCol += 4
		} else {
			trackCol++
		}

		if trackCol > visualCol {
			break
		}
	}

	return max(0, contentCol-1)
}

func removeLine(row int, lines []string) []string {
	if len(lines) == 0 {
		return lines
	}
	if row == len(lines)-1 {
		return lines[:row]
	}

	return append(lines[:row], lines[row+1:]...)
}

func removeChar(col int, line string) string {
	if col < 0 {
		return line
	}
	if len(line) < 1 {
		return ""
	}
	if col == len(line)-1 {
		return line[:len(line)-1]
	}

	start := line[:col]
	end := line[col+1:]
	return start + end
}

func insertChar(col int, line, char string) string {
	if col == 0 {
		return char + line
	}
	if col == len(line)-1 {
		return line + char
	}

	start := line[:col]
	end := line[col:]
	return start + char + end
}

func (c codeEditor) View(background lipgloss.Color) string {
	if c.width == 0 || c.height == 0 {
		return ""
	}

	styles := styles(background)
	view := []string{}

	for row, line := range c.content {
		pad := c.width - len(line)
		padding := styles["default"].Render(strings.Repeat(" ", pad))
		if row == c.cursor.row {
			view = append(view, renderCursorLine(c.cursor.col, line, styles)+padding)
		} else {
			view = append(view, renderLine(line, styles)+padding)
		}
	}

	return strings.Join(view, "\n")
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

func renderCursorLine(cursorCol int, line string, styles map[string]lipgloss.Style) string {
	lexer := parser.NewVerboseLexer([]byte(line))
	cursorDrawn := false
	col := 0
	view := ""

	for _, tok := range lexer.Lex() {
		if tok.Value == "\n" {
			continue
		}

		if cursorDrawn {
			view += styleToken(tok, styles)
			continue
		}

		start := col
		col += len(tok.Value)
		if cursorCol < col {
			view += styleCursorToken(cursorCol-start, tok, styles)
			cursorDrawn = true
		} else {
			view += styleToken(tok, styles)
		}
	}

	if !cursorDrawn {
		view += styleCursorToken(0, token.Token{Value: "_", Kind: token.Whitespace}, styles)
	}

	return view
}

func styles(background lipgloss.TerminalColor) map[string]lipgloss.Style {
	cursor := colors.Yellow3
	return map[string]lipgloss.Style{
		"muted":         lipgloss.NewStyle().Foreground(colors.Gray3).Background(background),
		"keyword":       lipgloss.NewStyle().Foreground(colors.Blue3).Background(background).Bold(true),
		"symbol":        lipgloss.NewStyle().Foreground(colors.Red3).Background(background).Bold(true),
		"number":        lipgloss.NewStyle().Foreground(colors.Purple4).Background(background),
		"string":        lipgloss.NewStyle().Foreground(colors.Yellow4).Background(background),
		"comment":       lipgloss.NewStyle().Foreground(colors.Green3).Background(background),
		"default":       lipgloss.NewStyle().Foreground(colors.White).Background(background),
		"cursorMuted":   lipgloss.NewStyle().Foreground(colors.Gray2).Background(cursor),
		"cursorKeyword": lipgloss.NewStyle().Foreground(colors.Blue3).Background(cursor).Bold(true),
		"cursorSymbol":  lipgloss.NewStyle().Foreground(colors.Red3).Background(cursor).Bold(true),
		"cursorNumber":  lipgloss.NewStyle().Foreground(colors.Purple1).Background(cursor),
		"cursorString":  lipgloss.NewStyle().Foreground(colors.Yellow4).Background(cursor),
		"cursorComment": lipgloss.NewStyle().Foreground(colors.Green3).Background(cursor),
		"cursorDefault": lipgloss.NewStyle().Foreground(colors.White).Background(cursor),
	}
}

func styleCursorToken(cursorLoc int, tok token.Token, styles map[string]lipgloss.Style) string {
	startValue := tok.Value[:cursorLoc]
	start := styleToken(token.Token{Value: startValue, Kind: tok.Kind}, styles)

	cursorValue := string(tok.Value[cursorLoc])
	cursor := styleCursor(token.Token{Value: cursorValue, Kind: tok.Kind}, styles)
	if cursorLoc == len(tok.Value)-1 {
		return start + cursor
	}

	endValue := tok.Value[cursorLoc+1:]
	end := styleToken(token.Token{Value: endValue, Kind: tok.Kind}, styles)
	return start + cursor + end
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
		value = strings.ReplaceAll(value, "\t", "├───")
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

func styleCursor(tok token.Token, styles map[string]lipgloss.Style) string {
	if len(tok.Value) == 0 {
		return ""
	}

	switch tok.Kind {
	case token.OpenParen, token.CloseParen:
		return styles["cursorMuted"].Render(tok.Value)
	case token.Whitespace:
		value := strings.ReplaceAll(tok.Value, " ", "·")
		value = strings.ReplaceAll(value, "\t", "├───")
		return styles["cursorMuted"].Render(value)
	case token.Let, token.Bool, token.GreaterThan, token.GreaterEqual,
		token.LessThan, token.LessEqual, token.Equal, token.Command:
		return styles["cursorKeyword"].Render(tok.Value)
	case token.Plus, token.Minus, token.Divide, token.Multiply:
		return styles["cursorSymbol"].Render(tok.Value)
	case token.Int, token.Float:
		return styles["cursorNumber"].Render(tok.Value)
	case token.String:
		return styles["cursorString"].Render(tok.Value)
	case token.Comment:
		return styles["cursorComment"].Render(tok.Value)
	default:
		return styles["cursorDefault"].Render(tok.Value)
	}
}
