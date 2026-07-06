package lexer

const (
	eof         byte = 0
	semicolon   byte = ';'
	comma       byte = ','
	lparen      byte = '('
	rparen      byte = ')'
	star        byte = '*'
	dot         byte = '.'
	doubleQuote byte = '"'
	backtick    byte = '`'
	lBracket    byte = '['
	rBracket    byte = ']'
	singleQuote byte = '\''
	xLower      byte = 'x'
	xUpper      byte = 'X'
	equals      byte = '='
	lt          byte = '<'
	gt          byte = '>'
	bang        byte = '!'
	pipe        byte = '|'
	plus        byte = '+'
	minus       byte = '-'
	slash       byte = '/'
	asterisk    byte = '*'
	underscore  byte = '_'
	dollar      byte = '$'
	space       byte = ' '
	tab         byte = '\t'
	cr          byte = '\r'
	newline     byte = '\n'
	zero        byte = '0'
	eLower      byte = 'e'
	eUpper      byte = 'E'
)

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isHexDigit(ch byte) bool {
	return isDigit(ch) || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}