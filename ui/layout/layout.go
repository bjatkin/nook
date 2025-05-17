package layout

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Content interface {
	Lines(maxWidth int) []string
	Shape() (int, int)
}

type Direction int

const (
	TopToBottom = Direction(iota)
	BottomToTop
	LeftToRight
	RightToLeft
)

// TODO: the text field should probably have padding directly on the struct
// that's gonna be way eaiser than trying to correctly pad multi line strings
// with the current padding functions
type Text struct {
	Text  string
	Style lipgloss.Style
}

func (t Text) Lines(maxWidth int) []string {
	if maxWidth <= 0 {
		return nil
	}

	lines := strings.Split(t.Text, "\n")
	width, _ := t.Shape()
	width = min(width, maxWidth)

	for i := range lines {
		l := ToWidth(int(width), lines[i])
		lines[i] = t.Style.Render(l)
	}

	return lines
}

func (t Text) Shape() (int, int) {
	width := 0
	lines := strings.Split(t.Text, "\n")
	for _, line := range lines {
		lineWidth := len(line)
		if width < lineWidth {
			width = lineWidth
		}
	}

	return width, len(lines)
}

type Container struct {
	Width     int
	Height    int
	Direction Direction
	Content   []Content
	Style     lipgloss.Style
}

func NewHContainer(width int, direction Direction, style lipgloss.Style) Container {
	return Container{
		Width:     width,
		Height:    1,
		Direction: direction,
		Style:     style,
	}
}

func (c Container) Shape() (int, int) {
	return c.Width, c.Height
}

func (c Container) String() string {
	lines := c.Lines(int(c.Width))
	return strings.Join(lines, "\n")
}

func (c Container) Lines(maxWidth int) []string {
	if maxWidth <= 0 {
		return nil
	}
	if c.Width == 0 || c.Height == 0 {
		return nil
	}

	buf := []string{}
	for i := 0; i < int(c.Height); i++ {
		buf = append(buf, "")
	}

	switch c.Direction {
	case LeftToRight:
		return c.leftToRightLines(buf)
	case RightToLeft:
		return c.rightToLeftLines(buf)
	case TopToBottom:
		return c.topToBottomLines(buf)
	case BottomToTop:
		return c.bottomToTopLines(buf)
	default:
		panic("unsupported layout direction")
	}
}

func (c Container) leftToRightLines(buf []string) []string {
	width := 0
	for _, content := range c.Content {
		lines := content.Lines(int(c.Width - width))
		if lines == nil {
			break
		}

		w, h := content.Shape()
		empty := strings.Repeat(" ", int(w))

		for i := 0; i < c.Height; i++ {
			if h <= i {
				buf[i] += empty
				continue
			}

			buf[i] += lines[i]
		}

		width += w
	}

	if width < c.Width {
		pad := strings.Repeat(" ", int(c.Width-width))
		pad = c.Style.Render(pad)
		for i := 0; i < int(c.Height); i++ {
			buf[i] += pad
		}
	}

	return buf
}

func (c Container) rightToLeftLines(buf []string) []string {
	width := 0
	for _, content := range c.Content {
		lines := content.Lines(int(c.Width - width))
		if lines == nil {
			break
		}

		w, h := content.Shape()
		empty := strings.Repeat(" ", int(w))

		for i := 0; i < c.Height; i++ {
			if h <= i {
				buf[i] = empty + buf[i]
				continue
			}

			buf[i] = lines[i] + buf[i]
		}

		width += w
	}

	if width < c.Width {
		pad := strings.Repeat(" ", int(c.Width-width))
		pad = c.Style.Render(pad)
		for i := 0; i < int(c.Height); i++ {
			buf[i] = pad + buf[i]
		}
	}

	return buf
}

func (c Container) topToBottomLines(buf []string) []string {
	height := 0
	for _, content := range c.Content {
		lines := content.Lines(int(c.Width))
		if lines == nil {
			break
		}

		w, h := content.Shape()
		pad := strings.Repeat(" ", c.Width-w)
		pad = c.Style.Render(pad)
		for i := 0; i < c.Height-height; i++ {
			if h <= i {
				break
			}

			buf[height] = lines[i] + pad
			height++
		}

		if height >= c.Height {
			break
		}
	}

	if height < c.Height {
		empty := strings.Repeat(" ", int(c.Width))
		empty = c.Style.Render(empty)
		for i := height; i < c.Height; i++ {
			buf[i] = empty
		}
	}

	return buf
}

func (c Container) bottomToTopLines(buf []string) []string {
	height := int(c.Height)
	for _, content := range c.Content {
		lines := content.Lines(int(c.Width))
		if lines == nil {
			break
		}

		w, h := content.Shape()
		pad := strings.Repeat(" ", c.Width-w)
		pad = c.Style.Render(pad)
		offset := height - h
		height -= h
		for i := 0; i < h; i++ {
			if offset+i < 0 {
				continue
			}

			buf[offset+i] = lines[i] + pad
		}

		if height <= 0 {
			break
		}
	}

	if height > 0 {
		empty := strings.Repeat(" ", c.Width)
		empty = c.Style.Render(empty)
		for i := height - 1; i >= 0; i-- {
			buf[i] = empty
		}
	}

	return buf
}

func PadLeft(pad int, text string) string {
	return strings.Repeat(" ", pad) + text
}

func PadRight(pad int, text string) string {
	return text + strings.Repeat(" ", pad)
}

func Pad(left, right int, text string) string {
	return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
}

func ToWidth(width int, text string) string {
	textWidth := len(text)
	if textWidth > width {
		return text[:width]
	}

	extend := width - textWidth
	return text + strings.Repeat(" ", extend)
}
