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

func TestTokenize_NotImplemented(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic from unimplemented NextToken")
		}
	}()
	l := NewLexer("SELECT 1")
	_ = l.Tokenize()
}