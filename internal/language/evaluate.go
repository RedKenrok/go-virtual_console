package language

import (
	"errors"
)

var evalDepth = 0

// Evaluate evaluates an expression within a given environment.
func Evaluate(
	expression Value,
	env *Environment,
) (Value, error) {
	evalDepth++
	defer func() {
		evalDepth--
	}()

	if evalDepth > 1e9 {
		return Value{}, errors.New("maximum evaluation depth exceeded")
	}

	switch expression.Type {
	case Bool, Float, Int, Option, String:
		return expression, nil

	case Lazy:
		thunk := expression.Data.(LazyData)
		value, err := Evaluate(
			thunk.Expression,
			thunk.Environment,
		)
		if err != nil {
			return Value{}, err
		}
		expression.Type = value.Type
		expression.Data = value.Data
		return expression, nil

	case List:
		list := expression.Data.([]Value)
		if len(list) == 0 {
			return Value{
				Type: Option,
				Data: OptionValue{
					Some: false,
				},
			}, nil
		}
		value, err := Evaluate(
			list[0],
			env,
		)
		if err != nil {
			return Value{}, err
		}

		if value.Type != Function && value.Type != Procedure {
			return Value{}, errors.New("first element in list is not a function or procedure")
		}

		if function, ok := value.Data.(func([]Value, *Environment) (Value, error)); ok {
			return function(list[1:], env)
		}
		return Value{}, errors.New("function or procedure is not callable")

	case Symbol:
		return env.Get(expression.Data.(string))
	}

	return Value{}, errors.New("unknown expression type")
}

func EvaluateUntilConcrete(
	expression Value,
	env *Environment,
) (Value, error) {
	value, err := Evaluate(expression, env)

	if err != nil {
		return Value{}, err
	}
	for value.Type == Lazy || value.Type == List {
		value, err = Evaluate(value, env)
		if err != nil {
			return Value{}, err
		}
	}
	return value, nil
}
