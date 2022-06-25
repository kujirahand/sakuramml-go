package sakuramml

import (
	"testing"
)

func TestSutotonSort(t *testing.T) {
	conv := NewSutotonConverter()
	conv.AddSutoton("音符", "l")
	conv.AddSutoton("音階", "o")
	conv.AddSutoton("トラック", "Track=")
	conv.AddSutoton("トラ", "TR=")
	conv.Sort()

	act, _ := conv.Convert("音符4 cde")
	exp := "l4 cde"
	if act != exp {
		t.Errorf("TestSutotonSort %s != %s", act, exp)
	}
	act2, _ := conv.Convert("トラック1 cde")
	exp2 := "Track=1 cde"
	if act2 != exp2 {
		t.Errorf("TestSutotonSort %s != %s", act2, exp2)
	}
}

func TestSutoton(t *testing.T) {
	conv := NewSutotonConverter()

	act1, _ := conv.Convert("~{音符}={l}~{トラック}={TR=}トラック2音符4 cde")
	exp1 := "TR=2l4 cde"
	if act1 != exp1 {
		t.Errorf("TestSutoton %s != %s", act1, exp1)
	}

	act2, _ := conv.Convert("~{トラック}={Track=}トラック2音符4 cde")
	exp2 := "Track=2l4 cde"
	if act2 != exp2 {
		t.Errorf("TestSutoton (上書き) %s != %s", act2, exp2)
	}
}
