package sakuramml

import (
	"fmt"
	"sort"
)

// Item struct
type Item struct {
	Key    string
	Value  string
	Length int
}

// SutotonConverter struct
type SutotonConverter struct {
	items []Item
}

// SetDefaultItem func
func (c *SutotonConverter) SetDefaultItem() {
	c.AddSutoton("ド", "c")
	c.AddSutoton("レ", "d")
	c.AddSutoton("ミ", "e")
	c.AddSutoton("ファ", "f")
	c.AddSutoton("フ", "f")
	c.AddSutoton("ソ", "g")
	c.AddSutoton("ラ", "a")
	c.AddSutoton("シ", "b")
	c.AddSutoton("ン", "r")
	c.AddSutoton("ッ", "r")
	c.AddSutoton("ー", "^")
	c.AddSutoton("↑", ">")
	c.AddSutoton("↓", "<")
	c.AddSutoton("【", "[")
	c.AddSutoton("】", "]")
	c.AddSutoton("｛", "{")
	c.AddSutoton("｝", "}")
	c.AddSutoton("（", "(")
	c.AddSutoton("）", ")")
	c.AddSutoton("０", "0")
	c.AddSutoton("１", "1")
	c.AddSutoton("２", "2")
	c.AddSutoton("３", "3")
	c.AddSutoton("４", "4")
	c.AddSutoton("５", "5")
	c.AddSutoton("６", "6")
	c.AddSutoton("７", "7")
	c.AddSutoton("８", "8")
	c.AddSutoton("９", "9")
	c.AddSutoton("♭", "-")
	c.AddSutoton("＃", "#")
	c.AddSutoton("♯", "#")
	c.AddSutoton("音量", "v")
	c.AddSutoton("音階", "o")
	c.AddSutoton("音符", "l")
	c.AddSutoton("音色", "@")
	c.AddSutoton("ゲート", "q")
	c.AddSutoton("トラック", "Track=")
	c.AddSutoton("チャンネル", "CH=")
	c.AddSutoton("テンポ", "Tempo")
	c.AddSutoton("読む", "Include")
	c.Sort()
}

// NewSutotonConverter func
func NewSutotonConverter() *SutotonConverter {
	conv := SutotonConverter{}
	conv.items = []Item{}
	return &conv
}

// AddSutoton func
func (conv *SutotonConverter) AddSutoton(key, value string) {
	if key == "" {
		return
	}
	newItem := Item{Key: key, Value: value, Length: len([]rune(key))}
	conv.items = append(conv.items, newItem)
}

// Find func
func (conv *SutotonConverter) Find(key string) int {
	for i := 0; i < len(conv.items); i++ {
		if conv.items[i].Key == key {
			return i
		}
	}
	return -1
}

// Sort func
func (conv *SutotonConverter) Sort() {
	sort.Slice(conv.items, func(i, j int) bool {
		return conv.items[i].Length > conv.items[j].Length
	})
}

// Convert
func (conv *SutotonConverter) Convert(text string) (string, error) {
	src := []rune(text)
	res := ""
	i := 0
	line := 0
	length := len(src)
MainLoop:
	for i < length {
		c := src[i]
		// Not multi bytes
		if c < 0x80 {
			if c == '\n' {
				line++
			}
			if c == '/' {
				// skip comment
				if StrCompareKey(src, i, "//") {
					comment := StrGetToken(src, &i, "\n")
					res += comment + "\n"
					line++
					continue
				}
				if StrCompareKey(src, i, "/*") {
					comment := StrGetRangeComment(src, &i)
					res += comment
					line += CountKey(comment, "\n")
					continue
				}
			}
			if c == '{' {
				// 明示的文字列はストトンで置換しない
				if StrCompareKey(src, i, "{\"") {
					str := StrGetToken(src, &i, "\"}")
					res += str + "\"}"
					line += CountKey(str, "\n")
					continue
				}
			}
			if c == '~' {
				// Sutoton New Sutoton
				if StrCompareKey(src, i, "~{") {
					i += 2
					key := StrGetToken(src, &i, "}")
					StrSkipSpace(src, &i)
					if src[i] != '=' { // 定義ではなかった?!
						return "", fmt.Errorf("[ERROR](%d)ストトン{%s}の定義エラー", line+1, key)
					}
					i++ // skip "="
					StrSkipSpace(src, &i)
					if src[i] != '{' {
						return "", fmt.Errorf("[ERROR](%d)ストトン{%s}の定義エラー", line+1, key)
					}
					i++ // skip "{"
					value := StrGetToken(src, &i, "}")
					line += CountKey(value, "\n")
					// Check Sutoton
					si := conv.Find(key)
					if si < 0 {
						conv.AddSutoton(key, value)
						conv.Sort()
					} else {
						conv.items[si].Value = value // Replace
					}
					continue
				}
			}
			res += string(c)
			i++
			continue
		}
		// Multi Bytes => Convert
		for j := 0; j < len(conv.items); j++ {
			it := conv.items[j]
			if StrCompareKey(src, i, it.Key) {
				res += it.Value
				i += it.Length
				continue MainLoop
			}
		}
		res += string(c)
		i++
	}
	return res, nil
}

// Convert
func SutotonConvert(src string) (string, error) {
	c := NewSutotonConverter()
	c.SetDefaultItem()
	return c.Convert(src)
}
