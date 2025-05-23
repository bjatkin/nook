package model

import (
	"reflect"
	"testing"
)

func Test_removeChar(t *testing.T) {
	type args struct {
		col  int
		line string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty line",
			args: args{
				col:  0,
				line: "",
			},
			want: "",
		},
		{
			name: "end of line",
			args: args{
				col:  5,
				line: "012345",
			},
			want: "01234",
		},
		{
			name: "beginning of line",
			args: args{
				col:  0,
				line: "012345",
			},
			want: "12345",
		},
		{
			name: "middle of line",
			args: args{
				col:  3,
				line: "012345",
			},
			want: "01245",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeChar(tt.args.col, tt.args.line); got != tt.want {
				t.Errorf("removeChar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_insertChar(t *testing.T) {
	type args struct {
		col  int
		line string
		char string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty line",
			args: args{
				col:  0,
				line: "",
				char: "0",
			},
			want: "0",
		},
		{
			name: "end of line",
			args: args{
				col:  3,
				line: "0123",
				char: "4",
			},
			want: "01234",
		},
		{
			name: "beginning of line",
			args: args{
				col:  0,
				line: "1234",
				char: "0",
			},
			want: "01234",
		},
		{
			name: "middle of line",
			args: args{
				col:  3,
				line: "01245",
				char: "3",
			},
			want: "012345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := insertChar(tt.args.col, tt.args.line, tt.args.char); got != tt.want {
				t.Errorf("insertChar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeLine(t *testing.T) {
	type args struct {
		row   int
		lines []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "remove first line",
			args: args{
				row:   0,
				lines: []string{"first", "second", "third"},
			},
			want: []string{"second", "third"},
		},
		{
			name: "remove last line",
			args: args{
				row:   2,
				lines: []string{"first", "second", "third"},
			},
			want: []string{"first", "second"},
		},
		{
			name: "remove middle line",
			args: args{
				row:   1,
				lines: []string{"first", "second", "third"},
			},
			want: []string{"first", "third"},
		},
		{
			name: "empty lines",
			args: args{
				row:   0,
				lines: []string{},
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeLine(tt.args.row, tt.args.lines); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toVisualColumn(t *testing.T) {
	type args struct {
		contentCol int
		line       string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "visual == content",
			args: args{
				contentCol: 3,
				line:       "012345",
			},
			want: 3,
		},
		{
			name: "single tab",
			args: args{
				contentCol: 3,
				line:       "\t12345",
			},
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toVisualColumn(tt.args.contentCol, tt.args.line); got != tt.want {
				t.Errorf("toVisualColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}
