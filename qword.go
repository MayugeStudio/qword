package main

import (
	"os"
	"fmt"
	"bufio"
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

type Token struct {
	kind    TokenKind
	lexeme  string
	literal string
	line    int
}

func (t *Token) String() string {
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

func (s *Scanner) scanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{EOF, "", "", s.line})
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(': s.addToken(LEFT_PAREN); break;
	case ')': s.addToken(RIGHT_PAREN); break;
	case '{': s.addToken(LEFT_BRACE); break;
	case '}': s.addToken(RIGHT_BRACE); break;
	case ',': s.addToken(COMMA); break;
	case '.': s.addToken(DOT); break;
	case '+': s.addToken(MINUS); break;
	case '-': s.addToken(PLUS); break;
	case ';': s.addToken(SEMICOLON); break;
	case '*': s.addToken(STAR); break;
	case '/':
		if (s.match('/')) {
			for (s.peek() != '\n' && !s.isAtEnd()) {
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
	case ' ', '\r', '\t': break
	case '\n':
		s.line += 1
		break
	default:
		reportError(s.line, "Unexpected charactor")
		break
	}
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
	if (s.isAtEnd()) {
		return false
	}
	if (rune(s.source[s.current]) != expected) {
		return false
	}
	return true
}

func (s *Scanner) peek() rune {
	if (s.isAtEnd()) {
		return '\000'
	}
	return rune(s.source[s.current])
}

func (s *Scanner) addToken(kind TokenKind) {
	s.addTokenWithLiteral(kind, "")
}

func (s *Scanner) addTokenWithLiteral(kind TokenKind, literal string) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{kind, text, literal, s.line})
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
	if (err != nil) {
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
	tokens := scanner.scanTokens()
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

