package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TokenKind int

const (
	EOF TokenKind = iota

	// 1 rune tokens
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	PLUS
	MINUS
	SLASH
	STAR
	SEMICOLON

	// 1 or 2 rune tokens
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals
	IDENTIFIER
	STRING
	NUMBER

	// Keywords
	TRUE
	FALSE
	AND
	OR
	VAR
	STRUCT
	FUN
	RETURN
	WHILE
	FOR
	IF
	ELSE
	PRINT
	NIL
)

var KEYWORDS = map[string]TokenKind{
	"false":  FALSE,
	"and":    AND,
	"or":     OR,
	"var":    VAR,
	"struct": STRUCT,
	"fun":    FUN,
	"return": RETURN,
	"while":  WHILE,
	"for":    FOR,
	"if":     IF,
	"else":   ELSE,
	"print":  PRINT,
	"nil":    NIL,
}

type LiteralKind int

const (
	LiteralNumber LiteralKind = iota
	LiteralString
	LiteralNone
)

type Literal struct {
	kind   LiteralKind
	number int
	str    string
}

func newLiteralNumber(value int) Literal {
	return Literal{
		kind:   LiteralNumber,
		number: value,
	}
}

func newLiteralString(value string) Literal {
	return Literal{
		kind: LiteralString,
		str:  value,
	}
}

func newLiteralNone() Literal {
	return Literal{
		kind: LiteralNone,
	}
}

type Token struct {
	kind    TokenKind
	lexeme  string
	literal Literal
	line    int
}

func (t *Token) String() string {
	if t.kind == NUMBER {
		return t.lexeme + " Literal: " + string(t.literal.number)
	} else if t.kind == STRING {
		return t.lexeme + " Literal: " + string(t.literal.str)
	}

	return t.lexeme
}

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func newScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  make([]Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) scanTokens() ([]Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			return nil, err
		}
	}

	s.tokens = append(s.tokens, Token{EOF, "", newLiteralNone(), s.line})
	return s.tokens, nil
}

func (s *Scanner) scanToken() error {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
		break
	case ')':
		s.addToken(RIGHT_PAREN)
		break
	case '{':
		s.addToken(LEFT_BRACE)
		break
	case '}':
		s.addToken(RIGHT_BRACE)
		break
	case ',':
		s.addToken(COMMA)
		break
	case '.':
		s.addToken(DOT)
		break
	case '+':
		s.addToken(MINUS)
		break
	case '-':
		s.addToken(PLUS)
		break
	case ';':
		s.addToken(SEMICOLON)
		break
	case '*':
		s.addToken(STAR)
		break
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
		break
	case '!':
		if s.match('=') {
			s.advance()
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
		break
	case '=':
		if s.match('=') {
			s.advance()
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(BANG)
		}
		break
	case '<':
		if s.match('=') {
			s.advance()
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
		break
	case '>':
		if s.match('=') {
			s.advance()
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
		break
	case ' ', '\r', '\t':
		break
	case '\n':
		s.line += 1
		break
	case '"':
		s.scanString()
		break
	default:
		if isDigit(c) {
			err := s.scanNumber()
			if err != nil {
				return err
			}
		} else if isAlpha(c) {
			s.scanIdentifier()
		} else {
			reportError(s.line, "Unexpected charactor.")
		}
		break
	}
	return nil
}

func (s *Scanner) scanString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line += 1
		}
		s.advance()
	}

	if s.isAtEnd() {
		reportError(s.line, "Unterminated string.")
	}

	s.advance() // eat right side double quotation

	value := s.source[s.start+1 : s.current-1]
	literal := newLiteralString(value)
	s.addTokenWithLiteral(STRING, literal)
}

func (s *Scanner) scanNumber() error {
	for isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	value, err := strconv.Atoi(s.source[s.start:s.current])
	if err != nil {
		return err
	}

	literal := newLiteralNumber(value)
	s.addTokenWithLiteral(NUMBER, literal)

	return nil
}

func (s *Scanner) scanIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	kind, ok := KEYWORDS[text]
	if !ok {
		kind = IDENTIFIER
	}
	s.addToken(kind)
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() rune {
	c := s.source[s.current]
	s.current += 1
	return rune(c)
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if rune(s.source[s.current]) != expected {
		return false
	}
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\000'
	}
	return rune(s.source[s.current])
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return rune(s.source[s.current+1])
}

func (s *Scanner) addToken(kind TokenKind) {
	s.addTokenWithLiteral(kind, Literal{kind: LiteralNone})
}

func (s *Scanner) addTokenWithLiteral(kind TokenKind, literal Literal) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{kind, text, literal, s.line})
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

type Runner struct {
	hadError bool
}

func newRunner() Runner {
	return Runner{
		hadError: false,
	}
}

func (r *Runner) runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = r.Run(string(bytes))
	if r.hadError {
		os.Exit(1)
	}
	return err
}

func (r *Runner) runPrompt() error {
	rd := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := rd.ReadString('\n')
		if err != nil {
			return err
		}
		line = strings.ReplaceAll(line, "\n", "")
		err = r.Run(line)
		if err != nil {
			return err
		}
		r.hadError = false
	}
	return nil
}

func (r *Runner) Run(source string) error {
	scanner := newScanner(source)
	tokens, err := scanner.scanTokens()
	if err != nil {
		return err
	}
	fmt.Println("Tokens", len(tokens))

	for _, token := range tokens {
		fmt.Println(token.String())
	}
	return nil
}

func main() {
	args := os.Args[1:]
	runner := newRunner()
	if len(args) > 1 {
		fmt.Println("Usage: qword [script]")
		os.Exit(1)
	} else if len(args) == 1 {
		runner.runFile(args[0])
	} else {
		runner.runPrompt()
	}
}

// TODO: Introduce Error handling interfaces
func reportError(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Printf("[line %d] Error %s: %s\n", line, where, message)
}
