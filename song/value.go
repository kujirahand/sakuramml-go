package song

import (
	"fmt"
	"strconv"
)

// SValue is value of sakura
type SValue interface {
	ToInt() int
	ToNum() float64
	ToStr() string
}

// SNumber is number of sakura
type SNumber float64

// ToInt : to int
func (n SNumber) ToInt() int {
	return int(n)
}

// ToNum : to num
func (n SNumber) ToNum() float64 {
	return float64(n)
}

// ToStr : to str
func (n SNumber) ToStr() string {
	return fmt.Sprintf("%f", float64(n))
}

// SStr is number of sakura
type SStr string

// ToInt : to int
func (n SStr) ToInt() int {
	v, _ := strconv.ParseInt(string(n), 10, strconv.IntSize)
	return int(v)
}

// ToNum : to num
func (n SStr) ToNum() float64 {
	v, _ := strconv.ParseFloat(string(n), 64)
	return v
}

// ToStr : to str
func (n SStr) ToStr() string {
	return string(n)
}
