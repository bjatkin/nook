package layout

import (
	"reflect"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func Test_Div_Render(t *testing.T) {
	type fields struct {
		direction _FlowDirection
		contents  []_Content
		width     int
		height    int
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "top to bottom",
			fields: fields{
				direction: _TopToBottom,
				contents: []_Content{
					NewText("hello\nworld\n!", lipgloss.NewStyle()),
					NewText("good\nbye\nworld\nagain!", lipgloss.NewStyle()),
					NewText("test\ntest", lipgloss.NewStyle()),
					NewText("final\ntest\nline", lipgloss.NewStyle()),
				},
				width:  20,
				height: 10,
			},
			want: []string{
				"hello",
				"world",
				"!    ",
				"good  ",
				"bye   ",
				"world ",
				"again!",
				"test",
				"test",
				"final",
			},
		},
		{
			name: "left to right",
			fields: fields{
				direction: _LeftToRight,
				contents: []_Content{
					NewText("test-ing\nthis\ntext!", lipgloss.NewStyle()),
					NewText(" | \n | \n | \n | \n ! ", lipgloss.NewStyle()),
					NewText("success!\nfailure", lipgloss.NewStyle()),
					NewText(" | \n | \n | \n | \n ! ", lipgloss.NewStyle()),
				},
				width:  20,
				height: 4,
			},
			want: []string{
				"test-ing | success! | ",
				"this     | failure  | ",
				"text!    |          | ",
				"         |          | ",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &_Div{
				direction: tt.fields.direction,
				contents:  tt.fields.contents,
				width:     tt.fields.width,
				height:    tt.fields.height,
			}
			if got := c.Render(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("_Div.Render() = \n%v, want \n%v", strings.Join(got, "\n"), strings.Join(tt.want, "\n"))
			}
		})
	}
}
