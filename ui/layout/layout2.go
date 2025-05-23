package layout

import (
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type _Content interface {
	Shape() (int, int)
	Render(int, int) []string
}

type _Text struct {
	lines  []string
	width  int
	height int
	style  lipgloss.Style
}

func NewText(text string, style lipgloss.Style) _Text {
	lines := strings.Split(text, "\n")
	height := len(lines)
	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}

	return _Text{
		lines:  lines,
		width:  width,
		height: height,
		style:  style,
	}
}

func (t _Text) Shape() (int, int) {
	return t.width, t.height
}

func (t _Text) Render(width, height int) []string {
	lines := []string{}
	for i := 0; i < height; i++ {
		if i >= len(t.lines) {
			lines = append(lines, strings.Repeat(" ", t.width))
			continue
		}

		line := setWidth(t.lines[i], t.width)
		lines = append(lines, t.style.Render(line))
	}

	return lines
}

func setWidth(line string, width int) string {
	switch {
	case len(line) > width:
		return line[:width]
	case len(line) < width:
		pad := width - len(line)
		return line + strings.Repeat(" ", pad)
	default:
		return line
	}
}

type _FlowDirection int

const (
	_TopToBottom = _FlowDirection(iota)
	_BottomToTop
	_LeftToRight
	_RightToLeft
)

type _Div struct {
	direction _FlowDirection
	contents  []_Content
	width     int
	height    int
}

func (c *_Div) Render() []string {
	switch c.direction {
	case _TopToBottom:
		lines := []string{}
		for _, content := range c.contents {
			_, height := content.Shape()
			block := content.Render(c.width, height)
			for _, line := range block {
				lines = append(lines, line)
				if len(lines) >= c.height {
					return lines
				}
			}
		}

		if c.height < len(lines) {
			return lines[:c.height]
		}
		return lines
	case _BottomToTop:
		lines := []string{}
		for _, content := range slices.Backward(c.contents) {
			_, height := content.Shape()
			block := content.Render(c.width, height)
			for _, line := range block {
				lines = append(lines, line)
				if len(lines) >= c.height {
					return lines
				}
			}
		}

		return lines
	case _LeftToRight:
		rows := []string{}
		for i := 0; i < c.height; i++ {
			rows = append(rows, "")
		}

		currentWidth := 0
		for _, content := range c.contents {
			width, _ := content.Shape()
			if currentWidth+width > c.width {
				width = c.width - currentWidth
			}

			currentWidth += width
			block := content.Render(width, c.height)
			for i := range rows {
				rows[i] += block[i]
			}

			if currentWidth >= c.width {
				return rows
			}
		}

		return rows
	case _RightToLeft:
		rows := []string{}
		for i := 0; i < c.height; i++ {
			rows = append(rows, "")
		}

		currentWidth := 0
		for _, content := range c.contents {
			width, _ := content.Shape()
			if currentWidth+width > c.width {
				width = c.width - currentWidth
			}

			currentWidth += width
			block := content.Render(width, c.height)
			for i := range rows {
				rows[i] = block[i] + rows[i]
			}

			if currentWidth >= c.width {
				return rows
			}
		}

		return rows
	default:
		panic("invalid flow direction")
	}
}
