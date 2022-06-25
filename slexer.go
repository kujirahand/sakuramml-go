package sakuramml

// Lexer 最低限必要な構造体を定義
type SLexer struct {
	src    []rune
	index  int
	lineno int
	fileno int
}

func newSLexer(src string, fileno int) *SLexer {
	l := SLexer{
		src:    []rune(src),
		index:  0,
		lineno: 0,
		fileno: fileno,
	}
	return &l
}

func (p *SLexer) readInt(defValue int) int {
	if p.isEOF() {
		return defValue
	}
	if !isDigit(p.peek()) {
		return defValue
	}
	result := int(p.nextRune() - '0')
	for !p.isEOF() {
		c := p.peek()
		if !isDigit(c) {
			break
		}
		p.next()
		result *= 10
		result += int(c - '0')
	}
	return result
}

func (p *SLexer) isEOF() bool {
	for p.index >= len(p.src) {
		return true
	}
	return false
}

func (p *SLexer) peek() rune {
	if p.isEOF() {
		return rune(0)
	}
	return p.src[p.index]
}

func (p *SLexer) peekNext() rune {
	if (p.index + 1) < len(p.src) {
		return p.src[p.index+1]
	}
	return rune(0)
}

func (p *SLexer) next() {
	p.index++
}

func (p *SLexer) nextRune() rune {
	r := p.src[p.index]
	p.index++
	return r
}

func (p *SLexer) skipSpace() {
	for !p.isEOF() {
		c := p.peek()
		if c == ' ' || c == '\t' || c == '\r' {
			p.next()
			continue
		}
		break
	}
}

func (p *SLexer) getLineInfo() LineInfo {
	return LineInfo{LineNo: p.lineno, FileNo: p.fileno}
}

func (p *SLexer) getChars(n int) string {
	result := ""
	for !p.isEOF() {
		if n > 0 {
			result += string(p.nextRune())
			n--
			continue
		}
		break
	}
	return result
}

func (p *SLexer) testStr(test string) bool {
	r := []rune(test)
	i := p.index
	realStr := p.getChars(len(r))
	p.index = i
	return test == realStr
}

func (p *SLexer) skipTo(test string) int {
	countLF := 0
	for !p.isEOF() {
		if p.peek() == '\n' {
			countLF++
		}
		if p.testStr(test) {
			p.index += len([]rune(test))
			break
		}
		p.next()
	}
	return countLF
}
