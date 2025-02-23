package language

const (
	STR_LIST_START = "["
	STR_LIST_END   = "]"
	STR_STRING     = "'"

	CHR_LIST_START = '['
	CHR_LIST_END   = ']'
	CHR_STRING     = '\''
	CHR_ESCAPE     = '\\'
)

// ValueType is an enumeration of our Lisp value kinds.
type ValueType int

const (
	Unknown ValueType = iota
	Function
	Procedure
	Lazy
	Symbol
	Option
	List
	Bool
	Int
	Float
	String
)

type LazyValue struct {
	Expression  Value
	Value       Value
	Environment *Environment
	Evaluated   bool
}

// OptionValue represents an optional value. Instead of using null, we wrap values in an Option.
type OptionValue struct {
	Value Value
	Some  bool
}

// Value is our generic container for interpreter values.
type Value struct {
	Data        interface{}
	Type        ValueType
	PreventEval bool
}
