package utils

import (
	"fmt"
)

// BytesToHex func
func BytesToHex(b []byte) string {
	s := ""
	for _, v := range b {
		s += fmt.Sprintf("%02x", v)
	}
	return s
}

// MidiRange func
func MidiRange(v int) int {
	if v < 0 {
		v = 0
	} else if v > 127 {
		v = 127
	}
	return v
}

// InRange func
func InRange(min, v, max int) int {
	if min > v {
		v = min
	} else if v > max {
		v = max
	}
	return v
}

// StrCompareKey func
func StrCompareKey(src []rune, index int, key string) bool {
	keyRune := []rune(key)
	return StrCompareRuneKey(src, index, keyRune)
}

// StrCompareRuneKey func
func StrCompareRuneKey(src []rune, index int, keyRune []rune) bool {
	maxLen := len(src)
	for i := 0; i < len(keyRune); i++ {
		if i >= maxLen {
			return false
		}
		if keyRune[i] != src[index + i] {
			return false
		}
	}
	return true
}

// StrGetToken func
func StrGetToken(src []rune, index *int, delimiter string) string {
	maxLen := len(src)
	delm := []rune(delimiter)
	delmLen := len(delm)
	res := ""
	for *index < maxLen {
		if StrCompareRuneKey(src, *index, delm) {
			*index += delmLen
			break
		}
		res += string(src[*index])
		*index++
	}
	return res
}

// StrSkipSpace func
func StrSkipSpace(src []rune, index *int) {
	length := len(src)
	for *index < length {
		c := src[*index]
		if c == ' ' || c == '\t' {
			*index++
			continue
		}
		break
	}
}

// StrSkipSpaceRet func
func StrSkipSpaceRet(src []rune, index *int) {
	length := len(src)
	for *index < length {
		c := src[*index]
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			*index++
			continue
		}
		break
	}
}

// StrGetRangeComment func
func StrGetRangeComment(src []rune, index *int) string {
	if !StrCompareKey(src, *index, "/*") {
		return ""
	}
	level := 0
	res := ""
	length := len(src)
	for *index < length {
		if StrCompareKey(src, *index, "/*") {
			level++
			*index += 2
			res += "/*"
			continue
		}
		if StrCompareKey(src, *index, "*/") {
			level--
			*index += 2
			res += "*/"
			if level == 0 {
				break
			}
		}
		res += string(src[*index])
		*index++
	}
	return res
}

func CountKey(source, key string) int {
	count := 0
	src := []rune(source)
	index := 0
	length := len(src)
	keyRune := []rune(key)
	for index < length {
		if StrCompareRuneKey(src, index, keyRune) {
			count++
			index += len(keyRune)
			continue
		}
		index++
	}
	return count
}

