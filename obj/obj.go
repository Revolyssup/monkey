package obj

import "fmt"

type DataType string

const (
	INTEGER_OBJ = "Integer"
	BOOLEAN_OBJ = "Bool"
	NULL_OBJ    = "Null"
	RETURN_OBJ  = "Return"
	ERROR_OBJ   = "Error"
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

//Implementing  Return
type Return struct {
	Value Object
}

func (ret *Return) DataType() DataType {
	return RETURN_OBJ
}

func (ret *Return) Inspect() string {
	return ret.Value.Inspect()
}

//Implementing Error object is similar to Return as they both stop the execution of program and return something
type Error struct {
	ErrMsg string
}

func (err *Error) DataType() DataType {
	return ERROR_OBJ
}

func (err *Error) Inspect() string {
	return "[MONKE ANGRY:] " + err.ErrMsg
}

//Environment object will passed around recursively in Eval

type Env struct {
	variables map[string]Object
}

func (env *Env) Get(s string) (Object, bool) {
	ob, ok := env.variables[s]
	return ob, ok
}

func (env *Env) Set(s string, ob Object) Object {
	env.variables[s] = ob
	return ob
}

func NewEnvironment() *Env {
	s := make(map[string]Object)
	env := &Env{variables: s}
	return env
}
