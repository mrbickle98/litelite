package lexer

import (
	"testing"
)

func TestNewLexer_NonNil(t *testing.T) {
	l := NewLexer("")
	if l == nil {
		t.Fatal("NewLexer returned nil")
	}
}

func TestTokenize_Select(t *testing.T) {
	l := NewLexer("SELECT id, name FROM users WHERE age >= 21 ;")
	toks := l.Tokenize()
	want := []Token{
		{TokenSelect, "SELECT", 1, 1},
		{TokenIdent, "id", 1, 8},
		{TokenComma, ",", 1, 10},
		{TokenIdent, "name", 1, 12},
		{TokenFrom, "FROM", 1, 17},
		{TokenIdent, "users", 1, 22},
		{TokenWhere, "WHERE", 1, 28},
		{TokenIdent, "age", 1, 34},
		{TokenGE, ">=", 1, 38},
		{TokenInt, "21", 1, 41},
		{TokenSemicolon, ";", 1, 44},
	}
	compareTokens(t, toks, want)
}

func TestTokenize_StringWithEscape(t *testing.T) {
	l := NewLexer("SELECT 'it''s ok'")
	toks := l.Tokenize()
	want := []Token{
		{TokenSelect, "SELECT", 1, 1},
		{TokenString, "'it''s ok'", 1, 8},
	}
	compareTokens(t, toks, want)
}

func TestTokenize_UnterminatedString(t *testing.T) {
	l := NewLexer("SELECT 'oops")
	toks := l.Tokenize()
	if len(toks) == 0 || toks[len(toks)-1].Type != TokenError {
		t.Fatalf("expected error token, got %v", toks)
	}
}

func TestTokenize_Blob(t *testing.T) {
	l := NewLexer("INSERT INTO t VALUES (x'deadbeef', X'00')")
	toks := l.Tokenize()
	if len(toks) < 7 {
		t.Fatalf("expected at least 7 tokens, got %d", len(toks))
	}
	if toks[5].Type != TokenBlob || toks[5].Value != "x'deadbeef'" {
		t.Errorf("first blob wrong: %v", toks[5])
	}
	if toks[6].Type != TokenComma {
		t.Errorf("comma missing: %v", toks[6])
	}
	if toks[7].Type != TokenBlob || toks[7].Value != "X'00'" {
		t.Errorf("second blob wrong: %v", toks[7])
	}
}

func TestTokenize_BlobMalformed(t *testing.T) {
	l := NewLexer("x'zz")
	toks := l.Tokenize()
	if len(toks) == 0 || toks[0].Type != TokenError {
		t.Fatalf("expected error token, got %v", toks)
	}
}

func TestTokenize_QuotedIdent(t *testing.T) {
	l := NewLexer(`SELECT "from", ` + "`where`" + `, [select]`)
	toks := l.Tokenize()
	if len(toks) != 6 {
		t.Fatalf("expected 6 tokens, got %d (%v)", len(toks), toks)
	}
	for _, i := range []int{1, 3, 5} {
		if toks[i].Type != TokenIdent {
			t.Errorf("quoted ident %d not TokenIdent: %v", i, toks[i])
		}
	}
	compareTokens(t, []Token{toks[1]}, []Token{{TokenIdent, `"from"`, 1, 8}})
	compareTokens(t, []Token{toks[3]}, []Token{{TokenIdent, "`where`", 1, 16}})
	compareTokens(t, []Token{toks[5]}, []Token{{TokenIdent, "[select]", 1, 25}})
}

func TestTokenize_Numbers(t *testing.T) {
	cases := map[string]TokenType{
		"42":     TokenInt,
		"0x1F":   TokenInt,
		"3.14":   TokenFloat,
		".5":     TokenFloat,
		"1e10":   TokenFloat,
		"1.5e-3": TokenFloat,
	}
	for src, wantType := range cases {
		l := NewLexer(src)
		toks := l.Tokenize()
		if len(toks) != 1 {
			t.Errorf("%q: expected 1 token, got %d (%v)", src, len(toks), toks)
			continue
		}
		if toks[0].Type != wantType {
			t.Errorf("%q: type = %v, want %v", src, toks[0].Type, wantType)
		}
		if toks[0].Value != src {
			t.Errorf("%q: value = %q, want %q", src, toks[0].Value, src)
		}
	}
}

func TestTokenize_Comments(t *testing.T) {
	l := NewLexer("SELECT -- a comment\nid /* block */ FROM t")
	toks := l.Tokenize()
	want := []Token{
		{TokenSelect, "SELECT", 1, 1},
		{TokenIdent, "id", 2, 1},
		{TokenFrom, "FROM", 2, 16},
		{TokenIdent, "t", 2, 21},
	}
	compareTokens(t, toks, want)
}

func TestTokenize_TrackingLineCol(t *testing.T) {
	l := NewLexer("SELECT id\nFROM\nusers")
	toks := l.Tokenize()
	want := []Token{
		{TokenSelect, "SELECT", 1, 1},
		{TokenIdent, "id", 1, 8},
		{TokenFrom, "FROM", 2, 1},
		{TokenIdent, "users", 3, 1},
	}
	compareTokens(t, toks, want)
}

func TestTokenize_UnexpectedByte(t *testing.T) {
	l := NewLexer("SELECT @")
	toks := l.Tokenize()
	if len(toks) != 2 || toks[1].Type != TokenError {
		t.Fatalf("expected error token, got %v", toks)
	}
}

func compareTokens(t *testing.T, got, want []Token) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("token count: got %d, want %d\ngot=%v\nwant=%v", len(got), len(want), got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("token %d: got %v, want %v", i, got[i], want[i])
		}
	}
}