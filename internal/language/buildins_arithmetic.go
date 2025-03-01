package language

import (
	"errors"
)

var addInts = Value{
	Type: Function,
	Data: func(
		args []Value,
		env *Environment,
	) (
		Value,
		error,
	) {
		var sumInt int64 = 0
		for _, arg := range args {
			evaluatedArg, err := EvaluateUntilConcrete(arg, env)
			if err != nil {
				return Value{}, err
			}
			if evaluatedArg.Type != Int {
				return Value{}, errors.New("arguments to add must be integers")
			}
			sumInt += evaluatedArg.Data.(int64)
		}
		return Value{
			Type: Int,
			Data: sumInt,
		}, nil
	},
}

var subtractInts = Value{
	Type: Function,
	Data: func(
		args []Value,
		env *Environment,
	) (
		Value,
		error,
	) {
		var resultInt int64 = 0
		for i, arg := range args {
			evaluatedArg, err := EvaluateUntilConcrete(arg, env)
			if err != nil {
				return Value{}, err
			}
			if evaluatedArg.Type != Int {
				return Value{}, errors.New("arguments to subtract must be integers")
			}
			number := evaluatedArg.Data.(int64)
			if i == 0 {
				resultInt = number
			} else {
				resultInt -= number
			}
		}
		return Value{
			Type: Int,
			Data: resultInt,
		}, nil
	},
}
