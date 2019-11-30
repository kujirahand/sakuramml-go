package lexer

import (
	"sakuramml/token"
	"testing"
)

func lexExec(t *testing.T, code, expected string) {
	tt, _ := Lex(code)
	a := token.TokensToString(tt, " ")
	if a != expected {
		t.Errorf("TestLex : %s != %s", a, expected)
	}
}

func TestLex(t *testing.T) {
	lexExec(t, "abc", "a b c")
	lexExec(t, "TR(3)a", "TR ( 3 ) a")
	lexExec(t, "o5", "o 5")
	lexExec(t, "TR=3 abc", "TR = 3 a b c")
}

func TestLex2(t *testing.T) {
	lexExec(t, "o5cde", "o 5 c d e")
	lexExec(t, "TR=3 [c]", "TR = 3 [ c ]")
	lexExec(t, "/* cde */", "")
	lexExec(t, "///hello\ncde", "/*hello*/ c d e")
}
