package obj

import "fmt"

type DataType string

const (
	INTEGER_OBJ = "Integer"
	BOOLEAN_OBJ = "Bool"
	NULL_OBJ    = "Null"
)

//All variables will be wrapped inside of an object-like struct.

type Object interface {
	DataType() DataType
	Inspect() string
}

//Implementing Integers

type Integer struct {
	Value int64
}

func (integer *Integer) DataType() DataType {
	return INTEGER_OBJ
}
func (integer *Integer) Inspect() string {
	return fmt.Sprintf("%d", integer.Value)
}

//Implementing Booleans
type Boolean struct {
	Value bool
}

func (boolean *Boolean) DataType() DataType {
	return BOOLEAN_OBJ
}

func (boolean *Boolean) Inspect() string {
	return fmt.Sprintf("%t", boolean.Value)
}

//Implementing Null
type Null struct{} //it holds no value
func (null *Null) DataType() DataType {
	return NULL_OBJ
}

func (null *Null) Inspect() string {
	return "null"
}
