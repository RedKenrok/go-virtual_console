package language

import "fmt"

// typeToString returns a string representation of a Type.
func typeToString(
	t ValueType,
) string {
	switch t {
	case Bool:
		return "bool"
	case Float:
		return "float"
	case Function:
		return "function"
	case Int:
		return "int"
	case Lazy:
		return "lazy"
	case List:
		return "list"
	case Option:
		return "option"
	case Procedure:
		return "procedure"
	case String:
		return "string"
	case Symbol:
		return "symbol"
	}
	return "unknown"
}

// valueToString returns a string representation of a Value.
func valueToString(
	value Value,
) string {
	switch value.Type {
	case Bool:
		if value.Data.(bool) {
			return "bool<true>"
		}
		return "bool<false>"

	case Float:
		return fmt.Sprintf("float<%g>", value.Data.(float64))

	case Function:
		return "function<>"

	case Int:
		return fmt.Sprintf("int<%d>", value.Data.(int64))

	case Lazy:
		return "lazy<" + valueToString(value.Data.(LazyData).Expression) + ">" + environmentToString(value.Data.(LazyData).Environment)

	case List:
		list := value.Data.([]Value)
		result := "list<"
		for i, elem := range list {
			if i > 0 {
				result += " "
			}
			result += valueToString(elem)
		}
		return result + ">"

	case Option:
		option := value.Data.(OptionValue)
		if option.Some {
			return "some<" + valueToString(option.Value) + ">"
		}
		return "none<>"

	case Procedure:
		return "procedure<>"

	case String:
		return "string<" + value.Data.(string) + ">"

	case Symbol:
		return "symbol<" + value.Data.(string) + ">"

	default:
		return "unknown<>"
	}
}

// environmentToString returns a string representation of an Environment.
func environmentToString(
	env *Environment,
) string {
	result := "{"
	for key, value := range env.Values {
		result += key + ": " + valueToString(value) + ", "
	}
	result += "}"
	return result
}

// valueEqual compares two Values for equality (with a little leniency for numbers).
func valueEqual(
	a Value,
	b Value,
	env *Environment,
) bool {
	if a.Type == b.Type {
		switch a.Type {
		case Bool:
			return a.Data.(bool) == b.Data.(bool)

		case Int:
			return a.Data.(int64) == b.Data.(int64)

		case Float:
			return a.Data.(float64) == b.Data.(float64)

		case String:
			return a.Data.(string) == b.Data.(string)

		case Symbol:
			aSymbol := a.Data.(string)
			bSymbol := b.Data.(string)
			if aSymbol == bSymbol {
				return true
			}

			aValue, err := env.Get(aSymbol)
			if err != nil {
				return false
			}
			bValue, err := env.Get(bSymbol)
			if err != nil {
				return false
			}
			return valueEqual(
				aValue,
				bValue,
				env,
			)

		case List:
			aList := a.Data.([]Value)
			bList := b.Data.([]Value)
			if len(aList) != len(bList) {
				return false
			}
			for i := range aList {
				if !valueEqual(aList[i], bList[i], env) {
					return false
				}
			}
			return true

		case Option:
			aOption := a.Data.(OptionValue)
			bOption := b.Data.(OptionValue)
			if aOption.Some != bOption.Some {
				return false
			}
			if !aOption.Some {
				return true
			}
			return valueEqual(
				aOption.Value,
				bOption.Value,
				env,
			)

		default:
			return false
		}
	}
	return false
}

func printEvaluation(
	message string,
	value Value,
	env *Environment,
) {
	spaces := ""
	for i := 0; i < evalDepth; i++ {
		spaces += "  "
	}
	println(spaces+message, valueToString(value), environmentToString(env))
}

func printKeyValue(
	message string,
	key string,
	value *Value,
) {
	spaces := ""
	for i := 0; i < evalDepth; i++ {
		spaces += "  "
	}
	if value == nil {
		println(spaces+message, key)
	} else {
		println(spaces+message, key, valueToString(*value))
	}
}
