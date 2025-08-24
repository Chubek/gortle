package lexer

import (
	"fmt"
	"patina-lang/internal/token"
	"strings"
	"unicode/utf8"
)

const EOF rune = -1

type StateFn func(*Lexer) StateFn

type Lexer struct {
	name   string
	input  string
	start  int
	pos    int
	width  int
	ch     byte
	tokens chan token.Token
	state  StateFn
}

func New(name, input string, initState StateFn) (l *Lexer) {
	l = &Lexer{
		name:   name,
		input:  input,
		tokens: make(chan token.Token, 2),
		state:  initState,
	}
	return l
}

func (l *Lexer) GetNextToken() token.Token {
	for {
		select {
		case token := <-l.tokens:
			return token
		default:
			l.state = l.state(l)
		}
	}
	panic("Unreachable")
}

func (l *Lexer) emit(t token.TokenType) {
	l.tokens <- token.Token{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *Lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return EOF
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *Lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *Lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *Lexer) errorf(format string, args ...interface{}) StateFn {
	l.tokens <- token.Token{
		token.TOKEN_Error,
		fmt.Sprintf(format, args...),
	}
	return nil
}
