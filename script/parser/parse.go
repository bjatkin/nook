package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bjatkin/nook/script/ast"
	"github.com/bjatkin/nook/script/token"
)

type Parser struct {
	lexer  lexer
	Errors []error
}

func NewParser(source []byte) *Parser {
	return &Parser{
		lexer: newLexer(source),
	}
}

func (p *Parser) AddError(err error) {
	p.Errors = append(p.Errors, err)
}

func (p *Parser) Parse() ast.Expr {
	kind := p.lexer.peek().Kind
	switch kind {
	case token.OpenParen:
		return p.ParseSExpr()

	case token.Atom:
		tok := p.lexer.take()
		return ast.Atom{Tok: tok, Value: tok.Value[1:]}

	case token.String:
		tok := p.lexer.take()
		value := tok.Value
		value = value[1 : len(value)-1]
		return ast.String{Tok: tok, Value: value}

	case token.Int:
		tok := p.lexer.take()

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
			p.AddError(fmt.Errorf("invalid integer '%s' %w", tok.Value, err))
			return nil
		}

		return ast.Int{Tok: tok, Value: i}

	case token.Float:
		tok := p.lexer.take()
		f, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			p.AddError(fmt.Errorf("invalid float '%s' %w", tok.Value, err))
			return nil
		}

		return ast.Float{Tok: tok, Value: f}

	case token.Bool:
		tok := p.lexer.take()
		value := false
		switch tok.Value {
		case "true":
			value = true
		case "false":
			value = false
		default:
			p.AddError(fmt.Errorf("invalid bool '%s'", tok.Value))
			return nil
		}

		return ast.Bool{Tok: tok, Value: value}

	case token.Identifier:
		return ast.Identifier{Value: p.lexer.take()}

	default:
		p.AddError(fmt.Errorf("unsupported expression '%#v'", p.lexer.take()))
		return nil
	}
}

func (p *Parser) ParseSExpr() ast.Expr {
	_ = p.lexer.take()

	operator := p.lexer.take()

	args := []ast.Expr{}
	for p.lexer.peek().Kind != token.CloseParen {
		if p.lexer.peek().Kind == token.EOF {
			p.AddError(fmt.Errorf("unclosed s-expr"))
			return nil
		}

		arg := p.Parse()
		args = append(args, arg)
	}
	_ = p.lexer.take()

	return ast.SExpr{
		Operator: operator,
		Operands: args,
	}
}
