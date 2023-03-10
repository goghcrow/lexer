package lexer

import (
	"fmt"
	"strings"
	"testing"
)

const (
	Number TokenKind = iota
	Ident
	NumId
	Space
	Comma
)

func stroftk(k TokenKind) string {
	return map[TokenKind]string{
		Number: "<num>",
		Ident:  "<id>",
		NumId:  "<numid>",
		Space:  "<space>",
		Comma:  ",",
	}[k]
}

func TestLexer(t *testing.T) {
	for _, tt := range []struct {
		input  string
		expect string
		lexer  *Lexer
	}{
		{
			"123",
			"<num>/123",
			BuildLexer(func(lex *Lexicon) {
				lex.Regex(Number, "\\d+")
				lex.Str(Comma, ",")
			}),
		},
		{
			"123,456",
			"<num>/123ð,/,ð<num>/456",
			BuildLexer(func(lex *Lexicon) {
				lex.Regex(Number, "\\d+")
				lex.Str(Comma, ",")
			}),
		},
		{
			"123,456,789",
			"<num>/123ð<num>/456ð<num>/789",
			BuildLexer(func(lex *Lexicon) {
				lex.Regex(Number, "\\d+")
				lex.Str(Comma, ",").Skip()
			}),
		},
		{
			"123, abc, 456, def, ",
			"<num>/123ð<id>/abcð<num>/456ð<id>/def",
			BuildLexer(func(lex *Lexicon) {
				lex.Regex(Number, "\\d+")
				lex.Regex(Ident, "[a-zA-Z]\\w*")
				lex.Regex(Space, "\\s+").Skip()
				lex.Str(Comma, ",").Skip()
			}),
		},
		{
			"123, abc, 456, def, ",
			"<numid>/123ð<numid>/abcð<numid>/456ð<numid>/def",
			BuildLexer(func(lex *Lexicon) {
				lex.Regex(NumId, "\\d+|(?:[a-zA-Z]\\w*)")
				lex.Regex(Space, "\\s+").Skip()
				lex.Str(Comma, ",").Skip()
			}),
		},
	} {
		t.Run(tt.input, func(t *testing.T) {
			toks := tt.lexer.MustLex(tt.input)
			actual := fmtToks(toks)
			if actual != tt.expect {
				t.Errorf("expect %s actual %s", tt.expect, actual)
			}
		})
	}
}

func fmtToks(toks []*Token) string {
	xs := make([]string, len(toks))
	for i, t := range toks {
		xs[i] = fmt.Sprintf("%s/%s", stroftk(t.TokenKind), t.Lexeme)
	}
	return strings.Join(xs, "ð")
}
