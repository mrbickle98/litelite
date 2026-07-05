package lexer

type Lexer struct {
	src     []byte
	pos     int
	readPos int
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
	_ = Token{}
	panic("lexer.NextToken not implemented")
}

func (l *Lexer) readChar() {
}

func (l *Lexer) peekChar() byte {
	return 0
}

func (l *Lexer) skipWhitespaceAndComments() {
}

func (l *Lexer) readIdent() Token {
	return Token{}
}

func (l *Lexer) readNumber() Token {
	return Token{}
}

func (l *Lexer) readString() Token {
	return Token{}
}

func (l *Lexer) readBlob() Token {
	return Token{}
}

func (l *Lexer) readQuotedIdent(quote byte) Token {
	return Token{}
}