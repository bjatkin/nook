package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bjatkin/nook/script/ast"
	"github.com/bjatkin/nook/script/token"
)

type Parser struct {
	lexer     lexer
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
	p.tokens = p.lexer.lex()

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
		return p.ParseSExpr()

	case token.Atom:
		tok := p.take()
		return ast.Atom{Tok: tok, Value: tok.Value}

	case token.String:
		tok := p.take()
		value := tok.Value
		value = value[1 : len(value)-1]
		return ast.String{Tok: tok, Value: value}

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

		return ast.Int{Tok: tok, Value: i}

	case token.Float:
		tok := p.take()
		f, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			p.addError(fmt.Errorf("invalid float '%s' %w", tok.Value, err))
			return nil
		}

		return ast.Float{Tok: tok, Value: f}

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

		return ast.Bool{Tok: tok, Value: value}

	case token.Flag:
		tok := p.take()
		return ast.Flag{Tok: tok, Value: tok.Value}
	case token.Path:
		tok := p.take()
		return ast.Path{Tok: tok, Value: tok.Value}
	case token.Identifier:
		return ast.Identifier{Value: p.take()}

	default:
		p.addError(fmt.Errorf("unsupported expression '%#v'", p.take()))
		return nil
	}
}

func (p *Parser) ParseSExpr() ast.Expr {
	_ = p.take()

	operator := p.take()

	args := []ast.Expr{}
	for p.peek().Kind != token.CloseParen {
		if p.peek().Kind == token.EOF {
			p.addError(fmt.Errorf("unclosed s-expr"))
			return nil
		}

		arg := p.parse()
		args = append(args, arg)
	}
	_ = p.take()

	return ast.SExpr{
		Operator: operator,
		Operands: args,
	}
}
