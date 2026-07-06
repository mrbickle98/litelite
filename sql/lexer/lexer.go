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
	case eof:
		return Token{Type: TokenEOF, Line: line, Col: col}
	case semicolon:
		l.readChar()
		return Token{Type: TokenSemicolon, Value: ";", Line: line, Col: col}
	case comma:
		l.readChar()
		return Token{Type: TokenComma, Value: ",", Line: line, Col: col}
	case lparen:
		l.readChar()
		return Token{Type: TokenLParen, Value: "(", Line: line, Col: col}
	case rparen:
		l.readChar()
		return Token{Type: TokenRParen, Value: ")", Line: line, Col: col}
	case star:
		l.readChar()
		return Token{Type: TokenStar, Value: "*", Line: line, Col: col}
	case dot:
		if isDigit(l.peekChar()) {
			return l.readNumber()
		}
		l.readChar()
		return Token{Type: TokenDot, Value: ".", Line: line, Col: col}
	case doubleQuote, backtick, lBracket:
		return l.readQuotedIdent(l.ch)
	case singleQuote:
		return l.readString()
	case xLower, xUpper:
		if l.peekChar() == singleQuote {
			return l.readBlob()
		}
		return l.readIdent()
	case equals:
		l.readChar()
		return Token{Type: TokenEq, Value: "=", Line: line, Col: col}
	case lt:
		if l.peekChar() == equals {
			l.readChar()
			l.readChar()
			// one readChar to move ch to "=" and the next readChar to
			// move ch to the next byte such that next thing can be evaluated
			return Token{Type: TokenLE, Value: "<=", Line: line, Col: col}
		}
		if l.peekChar() == gt {
			l.readChar()
			l.readChar()
			return Token{Type: TokenNE, Value: "<>", Line: line, Col: col}
		}
		l.readChar()
		return Token{Type: TokenLT, Value: "<", Line: line, Col: col}
	case gt:
		if l.peekChar() == equals {
			l.readChar()
			l.readChar()
			return Token{Type: TokenGE, Value: ">=", Line: line, Col: col}
		}
		l.readChar()
		return Token{Type: TokenGT, Value: ">", Line: line, Col: col}
	case bang:
		if l.peekChar() == equals {
			l.readChar()
			l.readChar()
			return Token{Type: TokenNE, Value: "!=", Line: line, Col: col}
		}
		tok := Token{Type: TokenError, Value: fmt.Sprintf("unexpected byte %q", l.ch), Line: line, Col: col}
		l.readChar()
		return tok
	case pipe:
		if l.peekChar() == pipe {
			l.readChar()
			l.readChar()
			return Token{Type: TokenConcat, Value: "||", Line: line, Col: col}
		}
		tok := Token{Type: TokenError, Value: fmt.Sprintf("unexpected byte %q", l.ch), Line: line, Col: col}
		l.readChar()
		return tok
	default:
		if isLetter(l.ch) || l.ch == underscore {
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

func (l *Lexer) readChar() {
	if l.nextPos >= len(l.src) {
		l.ch = eof
		l.pos = l.nextPos
	} else {
		l.ch = l.src[l.nextPos]
		l.pos = l.nextPos
		if l.ch == newline {
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
		return eof
	}
	return l.src[l.nextPos]
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		switch l.ch {
		case space, tab, cr, newline:
			l.readChar()
		case minus:
			if l.peekChar() == minus {
				for l.ch != eof && l.ch != newline {
					l.readChar()
				}
			} else {
				return
			}
		case slash:
			if l.peekChar() == star {
				l.readChar()
				l.readChar()
				for l.ch != eof {
					if l.ch == star && l.peekChar() == slash {
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
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == underscore || l.ch == dollar {
		l.readChar()
	}
	lexeme := string(l.src[start:l.pos])

	// default assumption is that we have a TokenIdent unless
	// we have an explicit match with any keywords
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

	if l.ch == zero && (l.peekChar() == xLower || l.peekChar() == xUpper) {
		l.readChar()
		l.readChar()
		// consume until we have exhausted all possible "ch" which could belong to
		// a hex
		for isHexDigit(l.ch) {
			l.readChar()
		}
		return Token{Type: TokenInt, Value: string(l.src[start:l.pos]), Line: line, Col: col}
	}

	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == dot {
		typ = TokenFloat
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	if l.ch == eLower || l.ch == eUpper {
		typ = TokenFloat
		l.readChar()
		if l.ch == plus || l.ch == minus {
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
	for l.ch != eof {
		if l.ch == singleQuote {
			if l.peekChar() == singleQuote {
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
	if l.ch != singleQuote {
		return Token{Type: TokenError, Value: "expected ' after blob prefix", Line: line, Col: col}
	}
	l.readChar()
	for isHexDigit(l.ch) {
		l.readChar()
	}
	if l.ch != singleQuote {
		return Token{Type: TokenError, Value: "malformed blob literal", Line: line, Col: col}
	}
	l.readChar()
	return Token{Type: TokenBlob, Value: string(l.src[start:l.pos]), Line: line, Col: col}
}

func (l *Lexer) readQuotedIdent(quote byte) Token {
	line, col := l.line, l.col
	start := l.pos
	close := quote
	if quote == lBracket {
		close = rBracket
	}
	l.readChar()
	for l.ch != eof {
		if l.ch == close {
			l.readChar()
			return Token{Type: TokenIdent, Value: string(l.src[start:l.pos]), Line: line, Col: col}
		}
		l.readChar()
	}
	return Token{Type: TokenError, Value: "unterminated quoted identifier", Line: line, Col: col}
}