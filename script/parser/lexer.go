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
	matchCommand,
	matchWhitespace,
	matchIdentifier,
}

type lexer struct {
	source               []byte
	pos                  uint
	nextToken            token.Token
	includeIgnoredTokens bool
}

func newLexer(source []byte) lexer {
	lexer := lexer{
		source:               source,
		pos:                  0,
		nextToken:            token.Token{},
		includeIgnoredTokens: false,
	}

	lexer.take()

	return lexer
}

func newVerboseLexer(source []byte) lexer {
	lexer := lexer{
		source:               source,
		pos:                  0,
		nextToken:            token.Token{},
		includeIgnoredTokens: false,
	}

	lexer.take()

	return lexer
}

func (l *lexer) lex() []token.Token {
	tokens := []token.Token{}
	for l.peek().Kind != token.EOF {
		tok := l.take()
		if l.includeIgnoredTokens {
			tokens = append(tokens, tok)
		}
		if tok.Kind == token.Whitespace || tok.Kind == token.Comment {
			continue
		}

		tokens = append(tokens, tok)
	}

	return tokens
}

// TODO: we should move away from lexing like this and instead to
// lex everything in one go, if we want performance we can just do
// in using a go-routine.
func (l *lexer) peek() token.Token {
	return l.nextToken
}

func isWhitespace(char byte) bool {
	return char == ' ' || char == '\n' || char == '\t' || char == ','
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
