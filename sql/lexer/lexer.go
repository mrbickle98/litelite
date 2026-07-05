package lexer

import (
	"fmt"
	"strings"
)

type Lexer struct {
	src     []byte
	pos     int
	nextPos int
	ch      byte
	line    int
	col     int
}

func NewLexer(src string) *Lexer {
	l := &Lexer{src: []byte(src), line: 1, col: 0}
	l.readChar()
	return l
}

func (l *Lexer) Tokenize() []Token {
	var out []Token
	for {
		t := l.NextToken()
		if t.Type == TokenEOF {
			break
		}
		out = append(out, t)
		if t.Type == TokenError {
			break
		}
	}
	return out
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespaceAndComments()

	line, col := l.line, l.col

	switch l.ch {
	case 0:
		return Token{Type: TokenEOF, Line: line, Col: col}
	case ';':
		l.readChar()
		return Token{Type: TokenSemicolon, Value: ";", Line: line, Col: col}
	case ',':
		l.readChar()
		return Token{Type: TokenComma, Value: ",", Line: line, Col: col}
	case '(':
		l.readChar()
		return Token{Type: TokenLParen, Value: "(", Line: line, Col: col}
	case ')':
		l.readChar()
		return Token{Type: TokenRParen, Value: ")", Line: line, Col: col}
	case '*':
		l.readChar()
		return Token{Type: TokenStar, Value: "*", Line: line, Col: col}
	case '.':
		if isDigit(l.peekChar()) {
			return l.readNumber()
		}
		l.readChar()
		return Token{Type: TokenDot, Value: ".", Line: line, Col: col}
	case '"', '`', '[':
		return l.readQuotedIdent(l.ch)
	case '\'':
		return l.readString()
	case 'x', 'X':
		if l.peekChar() == '\'' {
			return l.readBlob()
		}
		return l.readIdent()
	case '=':
		l.readChar()
		return Token{Type: TokenEq, Value: "=", Line: line, Col: col}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			// one readChar to move ch to "=" and the next readChar to
			// move ch to the next byte such that next thing can be evaluated
			return Token{Type: TokenLE, Value: "<=", Line: line, Col: col}
		}
		if l.peekChar() == '>' {
			l.readChar()
			l.readChar()
			return Token{Type: TokenNE, Value: "<>", Line: line, Col: col}
		}
		l.readChar()
		return Token{Type: TokenLT, Value: "<", Line: line, Col: col}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			return Token{Type: TokenGE, Value: ">=", Line: line, Col: col}
		}
		l.readChar()
		return Token{Type: TokenGT, Value: ">", Line: line, Col: col}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			l.readChar()
			return Token{Type: TokenNE, Value: "!=", Line: line, Col: col}
		}
		tok := Token{Type: TokenError, Value: fmt.Sprintf("unexpected byte %q", l.ch), Line: line, Col: col}
		l.readChar()
		return tok
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			l.readChar()
			return Token{Type: TokenConcat, Value: "||", Line: line, Col: col}
		}
		tok := Token{Type: TokenError, Value: fmt.Sprintf("unexpected byte %q", l.ch), Line: line, Col: col}
		l.readChar()
		return tok
	default:
		if isLetter(l.ch) || l.ch == '_' {
			return l.readIdent()
		}
		if isDigit(l.ch) {
			return l.readNumber()
		}
		tok := Token{Type: TokenError, Value: fmt.Sprintf("unexpected byte %q", l.ch), Line: line, Col: col}
		l.readChar()
		return tok
	}
}

/*
If we are not at the end
read the next byte at "nextPos".
Put that byte into "ch".
Update line and col accordingly
*/
func (l *Lexer) readChar() {
	if l.nextPos >= len(l.src) {
		l.ch = 0
		l.pos = l.nextPos
	} else {
		l.ch = l.src[l.nextPos]
		l.pos = l.nextPos
		if l.ch == '\n' {
			l.line++
			l.col = 0
		} else {
			l.col++
		}
	}
	l.nextPos++
}

func (l *Lexer) peekChar() byte {
	if l.nextPos >= len(l.src) {
		return 0
	}
	return l.src[l.nextPos]
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		switch l.ch {
		case ' ', '\t', '\r', '\n':
			l.readChar()
		case '-':
			if l.peekChar() == '-' {
				for l.ch != 0 && l.ch != '\n' {
					l.readChar()
				}
			} else {
				return
			}
		case '/':
			if l.peekChar() == '*' {
				l.readChar()
				l.readChar()
				for l.ch != 0 {
					if l.ch == '*' && l.peekChar() == '/' {
						l.readChar()
						l.readChar()
						break
					}
					l.readChar()
				}
			} else {
				return
			}
		default:
			return
		}
	}
}

func (l *Lexer) readIdent() Token {
	line, col := l.line, l.col
	start := l.pos
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' || l.ch == '$' {
		l.readChar()
	}
	lexeme := string(l.src[start:l.pos])
	typ := TokenIdent
	if kw, ok := keywords[strings.ToLower(lexeme)]; ok {
		typ = kw
	}
	return Token{Type: typ, Value: lexeme, Line: line, Col: col}
}

func (l *Lexer) readNumber() Token {
	line, col := l.line, l.col
	start := l.pos
	typ := TokenInt

	if l.ch == '0' && (l.peekChar() == 'x' || l.peekChar() == 'X') {
		l.readChar()
		l.readChar()
		for isHexDigit(l.ch) {
			l.readChar()
		}
		return Token{Type: TokenInt, Value: string(l.src[start:l.pos]), Line: line, Col: col}
	}

	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '.' {
		typ = TokenFloat
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	if l.ch == 'e' || l.ch == 'E' {
		typ = TokenFloat
		l.readChar()
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return Token{Type: typ, Value: string(l.src[start:l.pos]), Line: line, Col: col}
}

func (l *Lexer) readString() Token {
	line, col := l.line, l.col
	start := l.pos
	l.readChar()
	for l.ch != 0 {
		if l.ch == '\'' {
			if l.peekChar() == '\'' {
				l.readChar()
				l.readChar()
				continue
			}
			l.readChar()
			return Token{Type: TokenString, Value: string(l.src[start:l.pos]), Line: line, Col: col}
		}
		l.readChar()
	}
	return Token{Type: TokenError, Value: "unterminated string literal", Line: line, Col: col}
}

func (l *Lexer) readBlob() Token {
	line, col := l.line, l.col
	start := l.pos
	l.readChar()
	if l.ch != '\'' {
		return Token{Type: TokenError, Value: "expected ' after blob prefix", Line: line, Col: col}
	}
	l.readChar()
	for isHexDigit(l.ch) {
		l.readChar()
	}
	if l.ch != '\'' {
		return Token{Type: TokenError, Value: "malformed blob literal", Line: line, Col: col}
	}
	l.readChar()
	return Token{Type: TokenBlob, Value: string(l.src[start:l.pos]), Line: line, Col: col}
}

func (l *Lexer) readQuotedIdent(quote byte) Token {
	line, col := l.line, l.col
	start := l.pos
	close := quote
	if quote == '[' {
		close = ']'
	}
	l.readChar()
	for l.ch != 0 {
		if l.ch == close {
			l.readChar()
			return Token{Type: TokenIdent, Value: string(l.src[start:l.pos]), Line: line, Col: col}
		}
		l.readChar()
	}
	return Token{Type: TokenError, Value: "unterminated quoted identifier", Line: line, Col: col}
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isHexDigit(ch byte) bool {
	return isDigit(ch) || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

