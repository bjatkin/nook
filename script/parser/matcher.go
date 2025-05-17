package parser

import (
	"github.com/bjatkin/nook/script/token"
)

func isDecimal(char byte) bool {
	return char >= '0' && char <= '9'
}

func isOctal(char byte) bool {
	return char >= '0' && char <= '7'
}

func isBinary(char byte) bool {
	return char == '0' || char == '1'
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

func isAlpha(char byte) bool {
	if char >= 'a' && char <= 'z' {
		return true
	}

	if char >= 'A' && char <= 'Z' {
		return true
	}

	return false
}

type match struct {
	len  uint
	kind token.Kind
}

type matcher func(bytes []byte) *match

func matchSingleChar(bytes []byte) *match {
	if len(bytes) == 0 {
		return nil
	}

	switch bytes[0] {
	case '+':
		return &match{len: 1, kind: token.Plus}
	case '-':
		return &match{len: 1, kind: token.Minus}
	case '/':
		return &match{len: 1, kind: token.Divide}
	case '*':
		return &match{len: 1, kind: token.Multiply}
	case '>':
		return &match{len: 1, kind: token.GreaterThan}
	case '<':
		return &match{len: 1, kind: token.LessThan}
	case '(':
		return &match{len: 1, kind: token.OpenParen}
	case ')':
		return &match{len: 1, kind: token.CloseParen}
	case '{':
		return &match{len: 1, kind: token.OpenCurly}
	case '}':
		return &match{len: 1, kind: token.CloseCurly}
	case '[':
		return &match{len: 1, kind: token.OpenSquare}
	case ']':
		return &match{len: 1, kind: token.CloseSquare}
	case '.':
		return &match{len: 1, kind: token.Path}
	default:
		return nil
	}
}

func matchDoubleChar(bytes []byte) *match {
	if len(bytes) < 2 {
		return nil
	}

	switch string(bytes[:2]) {
	case ">=":
		return &match{len: 2, kind: token.GreaterEqual}
	case "<=":
		return &match{len: 2, kind: token.LessThan}
	case "==":
		return &match{len: 2, kind: token.Equal}
	case "./":
		return &match{len: 2, kind: token.Path}
	case "..":
		return &match{len: 2, kind: token.Path}
	default:
		return nil
	}
}

// matchLongPath matches only paths that start with either '/' or './' or '../'
func matchLongPath(bytes []byte) *match {
	if !matchPathPrefix(bytes) {
		return nil
	}

	// TODO: this needs to be WAAAAAYYY more robust, I'm missing a ton of valid paths here
	for i, char := range bytes {
		if char == '.' {
			continue
		}
		if char == '/' {
			continue
		}
		if char == '_' {
			continue
		}
		if char == '-' {
			continue
		}
		if isAlpha(char) {
			continue
		}
		if isDecimal(char) {
			continue
		}

		// a single slash is not a valid path. It must be followed by at least one other character
		// TODO: really I should be matching path's based of individual sections since an empty
		// directory is probably never valid
		if i == 1 && bytes[0] == '/' {
			return nil
		}

		return &match{len: uint(i), kind: token.Path}
	}

	return &match{len: uint(len(bytes)), kind: token.Path}
}

func matchPathPrefix(bytes []byte) bool {
	if len(bytes) > 0 && bytes[0] == '/' {
		return true
	}

	if len(bytes) >= 2 &&
		bytes[0] == '.' &&
		bytes[1] == '/' {
		return true
	}

	if len(bytes) >= 3 &&
		bytes[0] == '.' &&
		bytes[1] == '.' &&
		bytes[2] == '/' {
		return true
	}

	return false
}

func matchFlag(bytes []byte) *match {
	// need two bytes, first for - and second for valid flag name
	if len(bytes) < 2 {
		return nil
	}

	if bytes[0] != '-' {
		return nil
	}

	for i, char := range bytes[1:] {
		if char == '-' {
			continue
		}
		if isAlpha(char) {
			continue
		}
		if i > 0 && isDecimal(char) {
			continue
		}
		if i == 0 {
			return nil
		}

		return &match{len: uint(i + 1), kind: token.Flag}
	}

	return &match{len: uint(len(bytes)), kind: token.Flag}
}

func matchString(bytes []byte) *match {
	if bytes[0] != '"' {
		return nil
	}

	escape := false
	for i, char := range bytes[1:] {
		if char == '\\' {
			escape = true
			continue
		}
		if char == '"' && !escape {
			return &match{len: uint(i + 2), kind: token.String}
		}
		if char == '\n' {
			break
		}

		escape = false
	}

	return nil
}

func matchIdentifier(bytes []byte) *match {
	if len(bytes) == 0 {
		return nil
	}

	for i, char := range bytes {
		if isAlpha(char) {
			continue
		}
		if i > 0 && isDecimal(char) {
			continue
		}
		if i > 0 {
			kind := identifierKind(string(bytes[:i]))
			return &match{len: uint(i), kind: kind}
		}

		return nil
	}

	kind := identifierKind(string(bytes))
	return &match{len: uint(len(bytes)), kind: kind}
}

func matchComment(bytes []byte) *match {
	if bytes[0] != '#' {
		return nil
	}

	for i, char := range bytes {
		if char == '\n' {
			return &match{len: uint(i), kind: token.Comment}
		}
	}

	return &match{len: uint(len(bytes)), kind: token.Comment}
}

func matchAtom(bytes []byte) *match {
	if bytes[0] != '\'' {
		return nil
	}
	// atoms must be at least 2 bytes, 1 for the ' and one for a valid identifier name
	if len(bytes) < 1 {
		return nil
	}

	for i, char := range bytes[1:] {
		if isAlpha(char) {
			continue
		}
		if i > 0 {
			return &match{len: uint(i + 1), kind: token.Atom}
		}
	}

	return &match{len: uint(len(bytes)), kind: token.Atom}
}

func matchInt(bytes []byte) *match {
	if len(bytes) == 0 {
		return nil
	}

	if len(bytes) > 1 {
		prefix := string(bytes[:2])
		switch prefix {
		case "0x":
			len, ok := matchHex(bytes[2:])
			if ok {
				return &match{len: len + 2, kind: token.Int}
			}
		case "0o":
			len, ok := matchOctal(bytes[2:])
			if ok {
				return &match{len: len + 2, kind: token.Int}
			}
		case "0b":
			len, ok := matchBinary(bytes[2:])
			if ok {
				return &match{len: len + 2, kind: token.Int}
			}
		}
	}

	if len(bytes) > 2 {
		prefix := string(bytes[:3])
		switch prefix {
		case "-0x":
			len, ok := matchHex(bytes[3:])
			if ok {
				return &match{len: len + 3, kind: token.Int}
			}
		case "-0o":
			len, ok := matchOctal(bytes[3:])
			if ok {
				return &match{len: len + 3, kind: token.Int}
			}
		case "-0b":
			len, ok := matchBinary(bytes[3:])
			if ok {
				return &match{len: len + 3, kind: token.Int}
			}
		}
	}

	if bytes[0] == '-' {
		len, ok := matchDecimal(bytes[1:])
		if ok {
			return &match{len: len + 1, kind: token.Int}
		}
	}

	len, ok := matchDecimal(bytes)
	if ok {
		return &match{len: len, kind: token.Int}
	}

	return nil
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

func matchFloat(bytes []byte) *match {
	if !isDecimal(bytes[0]) {
		return nil
	}

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
			return &match{len: uint(i), kind: token.Float}
		}

		return nil
	}

	if mantisa {
		return &match{len: uint(len(bytes)), kind: token.Float}
	}

	return nil
}

func matchCommand(bytes []byte) *match {
	if bytes[0] != '$' {
		return nil
	}

	for i, char := range bytes[1:] {
		if isDecimal(char) {
			continue
		}
		if isAlpha(char) {
			continue
		}
		if char == '_' || char == '-' {
			continue
		}

		return &match{
			len:  uint(i + 1),
			kind: token.Command,
		}
	}

	return &match{
		len:  uint(len(bytes)),
		kind: token.Command,
	}
}

func isWhitespace(char byte) bool {
	return char == ' ' || char == '\n' || char == '\t' || char == ','
}

func matchWhitespace(bytes []byte) *match {
	if !isWhitespace(bytes[0]) {
		return nil
	}

	for i, char := range bytes {
		if isWhitespace(char) {
			continue
		}

		return &match{
			len:  uint(i),
			kind: token.Whitespace,
		}
	}

	return &match{
		len:  uint(len(bytes)),
		kind: token.Whitespace,
	}
}

func matchUnknownToken(bytes []byte) *match {
	for i, char := range bytes {
		if isWhitespace(char) {
			return &match{len: uint(i), kind: token.Invalid}
		}
	}

	return &match{len: uint(len(bytes)), kind: token.Invalid}
}

func identifierKind(value string) token.Kind {
	switch value {
	case "let":
		return token.Let
	case "true":
		return token.Bool
	case "false":
		return token.Bool
	case "int":
		return token.IntType
	case "float":
		return token.FloatType
	case "bool":
		return token.BoolType
	case "str":
		return token.StringType
	case "path":
		return token.PathType
	case "flag":
		return token.FlagType
	case "atom":
		return token.AtomType
	case "command":
		return token.CommandType
	default:
		return token.Identifier
	}
}
