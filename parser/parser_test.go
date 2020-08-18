package parser

import (
	"testing"

	"github.com/kujirahand/sakuramml-go/lexer"
)

func TestReadNoteOn(t *testing.T) {
	tok, _ := lexer.Lex("cd")
	nTop, _ := Parse(tok)
	act := nTop.ToStringAllName(" ")
	exp := "Nop NoteOn NoteOn"
	if act != exp {
		t.Errorf("NoteOn : %s != %s", act, exp)
	}
}
