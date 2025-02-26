package language

import "fmt"

// toString returns a string representation of a Value.
func toString(
	value Value,
) string {
	switch value.Type {
	case Bool:
		if value.Data.(bool) {
			return "true"
		}
		return "false"

	case Float:
		return fmt.Sprintf("%g", value.Data.(float64))

	case Function:
		return "<function>"

	case Int:
		return fmt.Sprintf("%d", value.Data.(int64))

	case Lazy:
		return "lazy " + toString(value.Data.(LazyData).Expression)

	case List:
		list := value.Data.([]Value)
		result := STR_LIST_START
		for i, elem := range list {
			if i > 0 {
				result += " "
			}
			result += toString(elem)
		}
		result += STR_LIST_END
		return result

	case Option:
		option := value.Data.(OptionValue)
		if option.Some {
			return "some " + toString(option.Value)
		}
		return "none"

	case Procedure:
		return "<procedure>"

	case String:
		return STR_STRING + value.Data.(string) + STR_STRING

	case Symbol:
		return value.Data.(string)

	default:
		return "unknown"
	}
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
