package language

import (
	"errors"
)

// Evaluate evaluates an expression within a given environment.
func Evaluate(
	expression Value,
	env *Environment,
) (Value, error) {
	switch expression.Type {
	case Symbol:
		return env.Get(expression.Data.(string))

	case Bool, Int, Float, String, Option:
		return expression, nil

	case Lazy:
		thunk := expression.Data.(LazyValue)
		if !thunk.Evaluated {
			value, err := Evaluate(
				thunk.Expression,
				thunk.Environment,
			)
			if err != nil {
				return Value{}, err
			}
			thunk.Value = value
			thunk.Evaluated = true
			// Update the thunk so future accesses are fast.
			expression.Data = thunk
		}
		return expression.Data.(LazyValue).Value, nil

	case List:
		list := expression.Data.([]Value)
		if len(list) == 0 {
			return expression, nil
		}
		value, err := Evaluate(list[0], env)
		if err != nil {
			return Value{}, err
		}

		if value.Type != Function && value.Type != Procedure {
			return Value{}, errors.New("first element in list is not a function or procedure")
		}

		if value.Type == Function {
			// Lazily evaluate arguments.
			if function, ok := value.Data.(func([]Value, *Environment) (Value, error)); ok {
				return function(list[1:], env)
			}
			return Value{}, errors.New("function is not callable")
		}
		// Type must be a procedure.

		// Eagerly evaluate arguments.
		var evaluatedArgs []Value
		for _, arg := range list[1:] {
			evaluatedArg, err := Evaluate(arg, env)
			if err != nil {
				return Value{}, err
			}
			evaluatedArgs = append(evaluatedArgs, evaluatedArg)
		}
		if procedure, ok := value.Data.(func([]Value, *Environment) (Value, error)); ok {
			return procedure(evaluatedArgs, env)
		}
		return Value{}, errors.New("procedure is not callable")

	default:
		return expression, nil
	}
}
