package sakuramml

import "strconv"

// Variable type
type VType = int

const (
	VTypeNone  = 0
	VTypeInt   = 1
	VTypeStr   = 2
	VTypeArray = 3
)

// Value struct
type Value struct {
	Type   VType
	SValue string
	IValue int
	AValue []*Value
}

// NewValue func
func NewValue() *Value {
	v := Value{Type: VTypeNone}
	return &v
}

// NewValueInt func init Value by int
func NewValueInt(value int) *Value {
	v := NewValue()
	v.Type = VTypeInt
	v.IValue = value
	return v
}

// NewValueStr func init Value by string
func NewValueStr(value string) *Value {
	v := NewValue()
	v.Type = VTypeStr
	v.SValue = value
	return v
}

// NewVaueArray func
func NewVaueArray() *Value {
	v := NewValue()
	v.Type = VTypeArray
	v.AValue = []*Value{}
	return v
}

func (v *Value) ToInt() int {
	return v.IValue
}

func (v *Value) ToString() string {
	if v.Type == VTypeInt {
		return strconv.Itoa(v.IValue)
	}
	return v.SValue
}

func (v *Value) AddValue(cv *Value) {
	v.AValue = append(v.AValue, cv)
}

// Variable struct
type Variable struct {
	values map[string]*Value
}

func NewVariable() *Variable {
	vv := Variable{}
	vv.values = map[string]*Value{}
	return &vv
}

func (vv *Variable) GetValue(name string) *Value {
	v, ok := vv.values[name]
	if !ok {
		return nil
	}
	return v
}

func (vv *Variable) GetIValue(name string, def int) int {
	v := vv.GetValue(name)
	if v == nil {
		return def
	}
	if v.Type == VTypeInt {
		return v.IValue
	} else {
		return def
	}
}

func (vv *Variable) GetSValue(name string, def string) string {
	v := vv.GetValue(name)
	if v == nil {
		return def
	}
	if v.Type == VTypeStr {
		return v.SValue
	} else {
		return strconv.Itoa(v.IValue)
	}
}

func (vv *Variable) Exists(name string) bool {
	v := vv.GetValue(name)
	return (v != nil)
}

func (vv *Variable) SetIValue(name string, value int) {
	vv.values[name] = NewValueInt(value)
}

func (vv *Variable) SetSValue(name string, value string) {
	vv.values[name] = NewValueStr(value)
}
