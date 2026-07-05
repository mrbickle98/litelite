package lexer

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenError

	TokenIdent // Token identifier
	TokenInt
	TokenFloat
	TokenString
	TokenBlob

	TokenSemicolon
	TokenComma
	TokenLParen
	TokenRParen
	TokenStar
	TokenDot
	TokenEq
	TokenLT
	TokenGT
	TokenLE
	TokenGE
	TokenNE
	TokenConcat
)

const (
	TokenKeywordBase TokenType = iota + 100
	TokenSelect
	TokenFrom
	TokenWhere
	TokenInsert
	TokenInto
	TokenValues
	TokenCreate
	TokenTable
)

type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

var keywords map[string]TokenType

func init() {
	keywords = map[string]TokenType{
		"select": TokenSelect,
		"from":   TokenFrom,
		"where":  TokenWhere,
		"insert": TokenInsert,
		"into":   TokenInto,
		"values": TokenValues,
		"create": TokenCreate,
		"table":  TokenTable,
	}
}
