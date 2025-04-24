package parser

import (
	"github.com/bjatkin/nook/script/token"
)

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

	var nextToken token.Token
	kind, matchSingle := matchSingleChar(l.source[l.pos])
	if matchSingle {
		nextToken = token.Token{
			Pos:   l.pos,
			Value: string(l.source[l.pos]),
			Kind:  kind,
		}
	}

	kind, matchDouble := matchDoubleChar(l.source[l.pos:])
	if matchDouble {
		nextToken = token.Token{
			Pos:   l.pos,
			Value: string(l.source[l.pos : l.pos+2]),
			Kind:  kind,
		}
	}

	if matchSingle || matchDouble {
		l.pos += uint(len(nextToken.Value))
		l.nextToken = nextToken
		return currentToken
	}

	count, match := matchFloat(l.source[l.pos:])
	if match {
		start := l.pos
		l.pos += count
		l.nextToken = token.Token{
			Pos:   start,
			Value: string(l.source[start:l.pos]),
			Kind:  token.Float,
		}
		return currentToken
	}

	count, match = matchInt(l.source[l.pos:])
	if match {
		start := l.pos
		l.pos += count
		l.nextToken = token.Token{
			Pos:   start,
			Value: string(l.source[start:l.pos]),
			Kind:  token.Int,
		}
		return currentToken
	}

	count, match = matchAtom(l.source[l.pos:])
	if match {
		start := l.pos
		l.pos += count
		l.nextToken = token.Token{
			Pos:   start,
			Value: string(l.source[start:l.pos]),
			Kind:  token.Atom,
		}
		return currentToken
	}

	count, match = matchString(l.source[l.pos:])
	if match {
		start := l.pos
		l.pos += count
		l.nextToken = token.Token{
			Pos:   start,
			Value: string(l.source[start:l.pos]),
			Kind:  token.String,
		}
		return currentToken
	}

	count, match = matchComment(l.source[l.pos:])
	if match {
		start := l.pos
		l.pos += count
		l.nextToken = token.Token{
			Pos:   start,
			Value: string(l.source[start:l.pos]),
			Kind:  token.Comment,
		}
		return currentToken
	}

	count, match = matchIdentifier(l.source[l.pos:])
	if match {
		start := l.pos
		l.pos += count
		value := string(l.source[start:l.pos])
		l.nextToken = token.Token{
			Pos:   start,
			Value: value,
			Kind:  identifierKind(value),
		}
		return currentToken
	}

	count = unknownToken(l.source[l.pos:])
	start := l.pos
	l.pos += count
	value := string(l.source[start:l.pos])
	l.nextToken = token.Token{
		Pos:   start,
		Value: value,
		Kind:  token.Invalid,
	}
	return currentToken
}

func matchSingleChar(char byte) (token.Kind, bool) {
	switch char {
	case '+':
		return token.Plus, true
	case '-':
		return token.Minus, true
	case '/':
		return token.Divide, true
	case '*':
		return token.Multiply, true
	case '>':
		return token.GreaterThan, true
	case '<':
		return token.LessThan, true
	case '(':
		return token.OpenParen, true
	case ')':
		return token.CloseParen, true
	default:
		return token.Invalid, false
	}
}

func matchDoubleChar(bytes []byte) (token.Kind, bool) {
	if len(bytes) <= 2 {
		return token.Invalid, false
	}

	switch string(bytes[:2]) {
	case ">=":
		return token.GreaterEqual, true
	case "<=":
		return token.LessThan, true
	case "==":
		return token.Equal, true
	}

	return token.Invalid, false
}

func isDecimal(char byte) bool {
	return char >= '0' && char <= '9'
}

func matchDecimal(bytes []byte) (uint, bool) {
	if len(bytes) == 0 {
		return 0, false
	}

	for i, char := range bytes {
		if isDecimal(char) {
			continue
		}
		if char == '_' {
			continue
		}
		if i > 0 {
			return uint(i), true
		}

		return 0, false
	}

	return uint(len(bytes)), true
}

func isHex(char byte) bool {
	if isDecimal(char) {
		return true
	}

	if char >= 'a' && char <= 'f' {
		return true
	}

	if char >= 'A' && char <= 'F' {
		return true
	}

	return false
}

func matchHex(bytes []byte) (uint, bool) {
	if len(bytes) == 0 {
		return 0, false
	}

	for i, char := range bytes {
		if isHex(char) {
			continue
		}
		if char == '_' {
			continue
		}
		if i > 0 {
			return uint(i), true
		}

		return 0, false
	}

	return uint(len(bytes)), true
}

func isOctal(char byte) bool {
	return char >= '0' && char <= '7'
}

func matchOctal(bytes []byte) (uint, bool) {
	if len(bytes) == 0 {
		return 0, false
	}

	for i, char := range bytes {
		if isOctal(char) {
			continue
		}
		if char == '_' {
			continue
		}
		if i > 0 {
			return uint(i), true
		}

		return 0, false
	}

	return uint(len(bytes)), true
}

func isBinary(char byte) bool {
	return char == '0' || char == '1'
}

func matchBinary(bytes []byte) (uint, bool) {
	if len(bytes) == 0 {
		return 0, false
	}

	for i, char := range bytes {
		if isBinary(char) {
			continue
		}
		if char == '_' {
			continue
		}
		if i > 0 {
			return uint(i), true
		}

		return 0, false
	}

	return uint(len(bytes)), true
}

func matchInt(bytes []byte) (uint, bool) {
	prefix := string(bytes[0])
	if len(bytes) > 1 {
		prefix = string(bytes[:2])
	}

	switch {
	case prefix == "0x":
		len, ok := matchHex(bytes[2:])
		if ok {
			return len + 2, true
		}
	case prefix == "0o":
		len, ok := matchOctal(bytes[2:])
		if ok {
			return len + 2, true
		}
	case prefix == "0b":
		len, ok := matchBinary(bytes[2:])
		if ok {
			return len + 2, true
		}
	default:
		len, ok := matchDecimal(bytes)
		if ok {
			return len, true
		}
	}

	return 0, false
}

func matchFloat(bytes []byte) (uint, bool) {
	mantisa := false
	for i, char := range bytes {
		if isDecimal(char) {
			continue
		}
		if !mantisa && char == '.' {
			mantisa = true
			continue
		}
		if i > 0 && mantisa {
			return uint(i), true
		}

		return 0, false
	}

	if mantisa {
		return uint(len(bytes)), true
	}

	return 0, false
}

func isAlpha(char byte) bool {
	if char >= 'a' && char <= 'z' {
		return true
	}

	if char >= 'A' && char <= 'Z' {
		return true
	}

	return false
}

func matchAtom(bytes []byte) (uint, bool) {
	if bytes[0] != ':' {
		return 0, false
	}

	for i, char := range bytes[1:] {
		if isAlpha(char) {
			continue
		}
		if i > 0 {
			return uint(i + 1), true
		}
	}

	if len(bytes) > 1 {
		return uint(len(bytes)), true
	}

	return 0, false
}

func matchString(bytes []byte) (uint, bool) {
	if bytes[0] != '"' {
		return 0, false
	}

	escape := false
	for i, char := range bytes[1:] {
		if char == '\\' {
			escape = true
			continue
		}
		if char == '"' && !escape {
			return uint(i + 2), true
		}
		if char == '\n' {
			break
		}

		escape = false
	}

	return 0, false
}

func matchIdentifier(bytes []byte) (uint, bool) {
	for i, char := range bytes {
		if isAlpha(char) {
			continue
		}
		if i > 0 && isDecimal(char) {
			continue
		}
		if i > 0 {
			return uint(i), true
		}

		return 0, false
	}
	if len(bytes) > 0 {
		return uint(len(bytes)), true
	}

	return 0, false
}

func matchComment(bytes []byte) (uint, bool) {
	if bytes[0] != '#' {
		return 0, false
	}

	for i, char := range bytes {
		if char == '\n' {
			return uint(i), true
		}
	}

	return uint(len(bytes)), true
}

func identifierKind(value string) token.Kind {
	switch value {
	case "let":
		return token.Let
	case "x":
		return token.Exec
	case "true":
		return token.Bool
	case "false":
		return token.Bool
	default:
		return token.Identifier
	}
}

func unknownToken(bytes []byte) uint {
	for i, char := range bytes {
		if isWhitespace(char) {
			return uint(i)
		}
	}

	return uint(len(bytes))
}
