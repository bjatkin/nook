package parser

import (
	"github.com/bjatkin/nook/script/token"
)

var matchers = []matcher{
	matchSingleChar,
	matchDoubleChar,
	matchFloat,
	matchInt,
	matchAtom,
	matchLongPath,
	matchFlag,
	matchString,
	matchComment,
	matchCommand,
	// matchWhitespace,
	matchIdentifier,
}

type Lexer struct {
	source               []byte
	pos                  uint
	includeIgnoredTokens bool
}

func newLexer(source []byte) Lexer {
	lexer := Lexer{
		source:               source,
		pos:                  0,
		includeIgnoredTokens: false,
	}

	return lexer
}

func NewVerboseLexer(source []byte) Lexer {
	lexer := Lexer{
		source:               source,
		pos:                  0,
		includeIgnoredTokens: true,
	}

	return lexer
}

func (l *Lexer) Lex() []token.Token {
	tokens := []token.Token{}
	tok := l.next()
	for tok.Kind != token.EOF {
		switch {
		case l.includeIgnoredTokens:
			tokens = append(tokens, tok)
		case tok.Kind != token.Whitespace && tok.Kind != token.Comment:
			tokens = append(tokens, tok)
		}

		tok = l.next()
	}

	return tokens
}

func (l *Lexer) next() token.Token {
	if int(l.pos) >= len(l.source) {
		return token.Token{
			Pos:  l.pos,
			Kind: token.EOF,
		}
	}

	var bestMatch *match
	for _, match := range matchers {
		found := match(l.source[l.pos:])
		if found == nil {
			continue
		}

		if bestMatch == nil {
			bestMatch = found
		}
		if bestMatch.len <= found.len {
			bestMatch = found
		}
	}
	if bestMatch == nil {
		bestMatch = matchUnknownToken(l.source[l.pos:])
	}

	start := l.pos
	l.pos += bestMatch.len
	return token.Token{
		Pos:   start,
		Value: string(l.source[start:l.pos]),
		Kind:  bestMatch.kind,
	}
}
