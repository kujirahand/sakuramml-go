package token

// TType type
type TType int

const (
	// Comment token
	Comment TType = iota
	// Word token
	Word
	// Macro token
	Macro
	// Number token
	Number
	// Flag token
	Flag
	// String token
	String
)

// Token struct
type Token struct {
	Type  TType
	Label string
	Line  int
}

// Tokens Slice
type Tokens []Token

// TokensToString for Debug
func TokensToString(tokens Tokens, delimiter string) string {
	s := ""
	for i, t := range tokens {
		// s += fmt.Sprintf("%3d: %5s %s\n", i, t.Type, t.Label)
		s += t.Label
		if i != len(tokens)-1 {
			s += delimiter
		}
	}
	return s
}

// Desk struct
type Desk struct {
	tokens Tokens
	index  int
}

// NewDesk func
func NewDesk(tt Tokens) Desk {
	d := Desk{tokens: tt, index: 0}
	return d
}

// Length tokens
func (desk *Desk) Length() int {
	return len(desk.tokens)
}

// HasNext func
func (desk *Desk) HasNext() bool {
	return (desk.index < len(desk.tokens))
}

// Peek func
func (desk *Desk) Peek() *Token {
	if desk.HasNext() {
		return &desk.tokens[desk.index]
	}
	return nil
}

// PeekN func
func (desk *Desk) PeekN(n int) *Token {
	idx := n + desk.index
	// range check
	if idx < 0 {
		return nil
	}
	if idx >= desk.Length() {
		return nil
	}
	return &desk.tokens[idx]
}

// Next func
func (desk *Desk) Next() *Token {
	if desk.HasNext() {
		v := desk.tokens[desk.index]
		desk.index++
		return &v
	}
	return nil
}

// Back func
func (desk *Desk) Back() {
	if desk.index > 0 {
		desk.index--
	}
}

// IsType func
func (desk *Desk) IsType(tt TType) bool {
	t := desk.Peek()
	if t == nil {
		return false
	}
	return t.Type == tt
}

// IsLabel func
func (desk *Desk) IsLabel(s string) bool {
	t := desk.Peek()
	if t == nil {
		return false
	}
	return t.Label == s
}
