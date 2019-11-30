package utils

import (
	"testing"
)

func TestBytesToHex(t *testing.T) {
	// - 1 -
	a := BytesToHex([]byte{1, 2, 3, 4})
	e := "01020304"
	if a != e {
		t.Errorf("TestBytesToHex: %s != %s", a, e)
		return
	}
	// - 2 -
	a2 := BytesToHex([]byte{0xFF, 2, 3, 4})
	e2 := "ff020304"
	if a2 != e2 {
		t.Errorf("TestBytesToHex: %s != %s", a2, e2)
		return
	}
}
func TestGetTokenRune(t *testing.T) {
	src := []rune("012::345--789")
	index := 0
	//
	act1 := StrGetToken(src, &index, "::")
	exp1 := "012"
	if act1 != exp1 {
		t.Errorf("TestGetTokenRune : %s != %s", act1, exp1)
		return
	}
	//
	act2 := StrGetToken(src, &index, "--")
	exp2 := "345"
	if act2 != exp2 {
		t.Errorf("TestGetTokenRune : %s != %s", act2, exp2)
		return
	}
	//
	act3 := StrGetToken(src, &index, "::")
	exp3 := "789"
	if act3 != exp3 {
		t.Errorf("TestGetTokenRune : %s != %s", act3, exp3)
		return
	}
}

func TestStrSkipSpaceRet(t *testing.T) {
	src := []rune("   012  \n 345  789")
	index := 0
	//
	StrSkipSpace(src, &index)
	act1 := src[index]
	exp1 := '0'
	if act1 != exp1 {
		t.Errorf("TestStrSkipSpaceRet : %c != %c", act1, exp1)
		return
	}
	//
	StrGetToken(src, &index, " ")
	//
	StrSkipSpaceRet(src, &index)
	act2 := src[index]
	exp2 := '3'
	if act2 != exp2 {
		t.Errorf("TestStrSkipSpaceRet : %c != %c", act2, exp2)
		return
	}
}

func TestStrGetRangeComment(t *testing.T) {
	src := []rune("012 /*345*/  /*123/*456*/789*/")
	index := 0
	StrGetToken(src, &index, " ")
	//
	act1 := StrGetRangeComment(src, &index)
	exp1 := "/*345*/"
	if act1 != exp1 {
		t.Errorf("TestStrGetRangeComment : %s != %s", act1, exp1)
		return
	}
	//
	StrSkipSpace(src, &index)
	act2 := StrGetRangeComment(src, &index)
	exp2 := "/*123/*456*/789*/"
	if act2 != exp2 {
		t.Errorf("TestStrGetRangeComment : %s != %s", act2, exp2)
		return
	}
}

func TestStrCountKey(t *testing.T) {
	act1 := CountKey("11*22*33*44*","*")
	exp1 := 4
	if act1 != exp1 {
		t.Errorf("TestStrCountKey : %d != %d", act1, exp1)
		return
	}

	act2 := CountKey("<>1111<>22222<>3333<>4444<>5555","<>")
	exp2 := 5
	if act2 != exp2 {
		t.Errorf("TestStrCountKey : %d != %d", act2, exp2)
		return
	}
}