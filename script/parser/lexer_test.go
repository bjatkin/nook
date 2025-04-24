package parser

import (
	"testing"
)

func Test_matchHex(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		wantLen uint
		wantOk  bool
	}{
		{
			name:    "single hex number",
			args:    args{bytes: []byte("A")},
			wantLen: 1,
			wantOk:  true,
		},
		{
			name:    "long hex number",
			args:    args{bytes: []byte("01_CF")},
			wantLen: 5,
			wantOk:  true,
		},
		{
			name:    "not a hex number",
			args:    args{bytes: []byte("Hello")},
			wantLen: 0,
			wantOk:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := matchHex(tt.args.bytes)
			if got != tt.wantLen {
				t.Errorf("matchHex() got = %v, wantLen %v", got, tt.wantLen)
			}
			if got1 != tt.wantOk {
				t.Errorf("matchHex() got1 = %v, wantOk %v", got1, tt.wantOk)
			}
		})
	}
}

func Test_matchOctal(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		wantLen uint
		wantOk  bool
	}{
		{
			name:    "octal number",
			args:    args{bytes: []byte("23_02")},
			wantLen: 5,
			wantOk:  true,
		},
		{
			name:    "not octal number",
			args:    args{bytes: []byte("Hello")},
			wantLen: 0,
			wantOk:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := matchOctal(tt.args.bytes)
			if got != tt.wantLen {
				t.Errorf("matchOctal() got = %v, want %v", got, tt.wantLen)
			}
			if got1 != tt.wantOk {
				t.Errorf("matchOctal() got1 = %v, want %v", got1, tt.wantOk)
			}
		})
	}
}

func Test_matchBinary(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		wantLen uint
		wantOk  bool
	}{
		{
			name:    "binary number",
			args:    args{bytes: []byte("01_01")},
			wantLen: 5,
			wantOk:  true,
		},
		{
			name:    "not octal number",
			args:    args{bytes: []byte("Hello")},
			wantLen: 0,
			wantOk:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := matchBinary(tt.args.bytes)
			if got != tt.wantLen {
				t.Errorf("matchBinary() got = %v, want %v", got, tt.wantLen)
			}
			if got1 != tt.wantOk {
				t.Errorf("matchBinary() got1 = %v, want %v", got1, tt.wantOk)
			}
		})
	}
}

func Test_matchNumber(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		wantLen uint
		wantOk  bool
	}{
		{
			name:    "decimal number",
			args:    args{bytes: []byte("012_456")},
			wantLen: 7,
			wantOk:  true,
		},
		{
			name:    "hex number",
			args:    args{bytes: []byte("0xAB_C0")},
			wantLen: 7,
			wantOk:  true,
		},
		{
			name:    "octal number",
			args:    args{bytes: []byte("0o17_46")},
			wantLen: 7,
			wantOk:  true,
		},
		{
			name:    "binary number",
			args:    args{bytes: []byte("0b0110_1111")},
			wantLen: 11,
			wantOk:  true,
		},
		{
			name:    "keyword",
			args:    args{bytes: []byte("else")},
			wantLen: 0,
			wantOk:  false,
		},
		{
			name:    "invalid number",
			args:    args{bytes: []byte("10)")},
			wantLen: 2,
			wantOk:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := matchInt(tt.args.bytes)
			if got != tt.wantLen {
				t.Errorf("matchNumber() got = %v, want %v", got, tt.wantLen)
			}
			if got1 != tt.wantOk {
				t.Errorf("matchNumber() got1 = %v, want %v", got1, tt.wantOk)
			}
		})
	}
}

func Test_matchFloat(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		wantLen uint
		wantOk  bool
	}{
		{
			name:    "float number",
			args:    args{bytes: []byte("43.54")},
			wantLen: 5,
			wantOk:  true,
		},
		{
			name:    "int number",
			args:    args{bytes: []byte("1234")},
			wantLen: 0,
			wantOk:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := matchFloat(tt.args.bytes)
			if got != tt.wantLen {
				t.Errorf("matchFloat() got = %v, want %v", got, tt.wantLen)
			}
			if got1 != tt.wantOk {
				t.Errorf("matchFloat() got1 = %v, want %v", got1, tt.wantOk)
			}
		})
	}
}

func Test_matchAtom(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		wantLen uint
		wantOk  bool
	}{
		{
			name:    "valid atom",
			args:    args{bytes: []byte(":validAtom")},
			wantLen: 10,
			wantOk:  true,
		},
		{
			name:    "short atom",
			args:    args{bytes: []byte(":ok")},
			wantLen: 3,
			wantOk:  true,
		},
		{
			name:    "invalid atom",
			args:    args{bytes: []byte("test")},
			wantLen: 0,
			wantOk:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLen, gotOk := matchAtom(tt.args.bytes)
			if gotLen != tt.wantLen {
				t.Errorf("matchAtom() got = %v, want %v", gotLen, tt.wantLen)
			}
			if gotOk != tt.wantOk {
				t.Errorf("matchAtom() got1 = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_matchString(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		wantLen uint
		wantOk  bool
	}{
		{
			name:    "valid string",
			args:    args{bytes: []byte("\"hello world\"")},
			wantLen: 13,
			wantOk:  true,
		},
		{
			name:    "invalid string",
			args:    args{bytes: []byte("test")},
			wantLen: 0,
			wantOk:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := matchString(tt.args.bytes)
			if got != tt.wantLen {
				t.Errorf("matchString() got = %v, want %v", got, tt.wantLen)
			}
			if got1 != tt.wantOk {
				t.Errorf("matchString() got1 = %v, want %v", got1, tt.wantOk)
			}
		})
	}
}

func Test_lexer_matchIdentifier(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		wantLen uint
		wantOk  bool
	}{
		{
			name:    "invalid identifier",
			args:    args{bytes: []byte("10")},
			wantLen: 0,
			wantOk:  false,
		},
		{
			name:    "valid identifier",
			args:    args{bytes: []byte("test10")},
			wantLen: 6,
			wantOk:  true,
		},
		{
			name:    "trailing space",
			args:    args{bytes: []byte("a ")},
			wantLen: 1,
			wantOk:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := matchIdentifier(tt.args.bytes)
			if got != tt.wantLen {
				t.Errorf("lexer.matchIdentifier() got = %v, want %v", got, tt.wantLen)
			}
			if got1 != tt.wantOk {
				t.Errorf("lexer.matchIdentifier() got1 = %v, want %v", got1, tt.wantOk)
			}
		})
	}
}
