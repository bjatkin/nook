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
	matchPath,
	matchFlag,
	matchString,
	matchComment,
	matchIdentifier,
}

type lexer struct {
	source    []byte
	pos       uint
	nextToken token.Token
}

func newLexer(source []byte) lexer {
	lexer := lexer{
		source:    source,
		pos:       0,
		nextToken: token.Token{},
	}

	lexer.take()

	return lexer
}

func (l *lexer) peek() token.Token {
	return l.nextToken
}

func isWhitespace(char byte) bool {
	return char == ' ' || char == '\n' || char == '\t'
}

func (l *lexer) take() token.Token {
	currentToken := l.nextToken

	for int(l.pos) < len(l.source) && isWhitespace(l.source[l.pos]) {
		l.pos++
	}

	if int(l.pos) >= len(l.source) {
		l.nextToken = token.Token{
			Pos:  l.pos,
			Kind: token.EOF,
		}
		return currentToken
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
	l.nextToken = token.Token{
		Pos:   start,
		Value: string(l.source[start:l.pos]),
		Kind:  bestMatch.kind,
	}

	return currentToken
}
