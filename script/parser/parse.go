package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bjatkin/nook/script/ast"
	"github.com/bjatkin/nook/script/token"
)

type Parser struct {
	lexer     Lexer
	tokens    []token.Token
	nextToken uint
	Errors    []error
}

func NewParser(source []byte) *Parser {
	return &Parser{
		lexer: newLexer(source),
	}
}

func (p *Parser) addError(err error) {
	p.Errors = append(p.Errors, err)
}

func (p *Parser) Parse() ast.Expr {
	p.tokens = p.lexer.Lex()

	return p.parse()
}

func (p *Parser) peek() token.Token {
	if p.nextToken > uint(len(p.tokens)) {
		return token.Token{Kind: token.EOF}
	}

	return p.tokens[p.nextToken]
}

func (p *Parser) take() token.Token {
	tok := p.peek()

	p.nextToken++
	return tok
}

func (p *Parser) parse() ast.Expr {
	kind := p.peek().Kind
	switch kind {
	case token.OpenParen:
		_ = p.take() // take the '('
		args := p.parseArgs(token.CloseParen)
		_ = p.take() // take the ')'

		if len(args) == 0 {
			p.addError(fmt.Errorf("() is not a valid s-expression"))
			return nil
		}

		return &ast.SExpr{
			Operator: args[0],
			Operands: args[1:],
		}

	case token.OpenCurly:
		open := p.take() // take the '{'
		args := p.parseArgs(token.CloseCurly)
		_ = p.take() // take the '}'

		return &ast.SExpr{
			Operator: &ast.SCurly{Tok: open},
			Operands: args,
		}

	case token.OpenSquare:
		open := p.take() // take the '['
		args := p.parseArgs(token.CloseSquare)
		_ = p.take() // take the ']'

		return &ast.SExpr{
			Operator: &ast.SSquare{Tok: open},
			Operands: args,
		}
	case token.Atom:
		tok := p.take()
		return &ast.Atom{Tok: tok, Value: tok.Value}

	case token.Command:
		tok := p.take()
		return &ast.SCommand{Tok: tok}

	case token.String:
		tok := p.take()
		value := tok.Value
		value = value[1 : len(value)-1]
		return &ast.String{Tok: tok, Value: value}

	case token.Int:
		tok := p.take()

		var i int64
		var err error
		switch {
		case strings.HasPrefix(tok.Value, "0x"):
			value := strings.TrimPrefix(tok.Value, "0x")
			i, err = strconv.ParseInt(value, 16, 64)
		case strings.HasPrefix(tok.Value, "0b"):
			value := strings.TrimPrefix(tok.Value, "0b")
			i, err = strconv.ParseInt(value, 2, 64)
		case strings.HasPrefix(tok.Value, "0o"):
			value := strings.TrimPrefix(tok.Value, "0o")
			i, err = strconv.ParseInt(value, 8, 64)
		default:
			i, err = strconv.ParseInt(tok.Value, 10, 64)
		}
		if err != nil {
			p.addError(fmt.Errorf("invalid integer '%s' %w", tok.Value, err))
			return nil
		}

		return &ast.Int{Tok: tok, Value: i}

	case token.Float:
		tok := p.take()
		f, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			p.addError(fmt.Errorf("invalid float '%s' %w", tok.Value, err))
			return nil
		}

		return &ast.Float{Tok: tok, Value: f}

	case token.Bool:
		tok := p.take()
		value := false
		switch tok.Value {
		case "true":
			value = true
		case "false":
			value = false
		default:
			p.addError(fmt.Errorf("invalid bool '%s'", tok.Value))
			return nil
		}

		return &ast.Bool{Tok: tok, Value: value}

	case token.Flag:
		tok := p.take()
		return &ast.Flag{Tok: tok, Value: tok.Value}
	case token.Path:
		tok := p.take()
		return &ast.Path{Tok: tok, Value: tok.Value}
	case token.Identifier:
		tok := p.take()
		return &ast.Identifier{Tok: tok, Name: tok.Value}
	case token.Plus, token.Minus, token.Multiply, token.Divide:
		tok := p.take()
		return &ast.Identifier{Tok: tok, Name: tok.Value}
	case token.Let:
		tok := p.take()
		return &ast.SLet{Tok: tok}
	default:
		p.addError(fmt.Errorf("unsupported expression '%#v'", p.take()))
		return nil
	}
}

func (p *Parser) parseArgs(closeToken token.Kind) []ast.Expr {
	args := []ast.Expr{}
	for p.peek().Kind != closeToken {
		if p.peek().Kind == token.EOF {
			p.addError(fmt.Errorf("unclosed expression list"))
			return nil
		}

		arg := p.parse()
		args = append(args, arg)
	}
	return args
}
