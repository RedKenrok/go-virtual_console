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

	case Int:
		return fmt.Sprintf("%d", value.Data.(int64))

	case Float:
		return fmt.Sprintf("%g", value.Data.(float64))

	case Symbol:
		return value.Data.(string)

	case String:
		return STR_STRING + value.Data.(string) + STR_STRING

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

	case Procedure:
		return "<procedure>"

	case Function:
		return "<function>"

	case Option:
		option := value.Data.(OptionValue)
		if option.Some {
			return "some " + toString(option.Value)
		}
		return "none"

	default:
		return "unknown"
	}
}

// valueEqual compares two Values for equality (with a little leniency for numbers).
func valueEqual(
	a Value,
	b Value,
) bool {
	if a.Type == b.Type {
		switch a.Type {
		case Bool:
			return a.Data.(bool) == b.Data.(bool)

		case Int:
			return a.Data.(int64) == b.Data.(int64)

		case Float:
			return a.Data.(float64) == b.Data.(float64)

		case Symbol, String:
			return a.Data.(string) == b.Data.(string)

		case List:
			aList := a.Data.([]Value)
			bList := b.Data.([]Value)
			if len(aList) != len(bList) {
				return false
			}
			for i := range aList {
				if !valueEqual(aList[i], bList[i]) {
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
			return valueEqual(aOption.Value, bOption.Value)

		default:
			return false
		}
	}
	return false
}
