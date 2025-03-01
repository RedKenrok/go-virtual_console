package language

import (
	"errors"
)

func AddBuiltins(
	env *Environment,
) {
	env.Set("define", Value{
		Type: Function,
		Data: func(
			args []Value,
			env *Environment,
		) (
			Value,
			error,
		) {
			// [define symbol expression]
			if len(args) != 2 {
				return Value{}, errors.New("define requires 2 arguments")
			}
			symbolValue := args[0]
			if symbolValue.Type != Symbol {
				return Value{}, errors.New("first argument to define must be a symbol")
			}
			result, err := Evaluate(args[1], env)
			if err != nil {
				return Value{}, err
			}
			env.Set(symbolValue.Data.(string), result)
			return result, nil
		},
	})

	env.Set("function", Value{
		Type: Function,
		Data: func(
			args []Value,
			env *Environment,
		) (Value, error) {
			// [function [parameters] body]
			if len(args) != 2 {
				return Value{}, errors.New("function requires 2 arguments")
			}

			parametersValue := args[0]
			if parametersValue.Type != List {
				return Value{}, errors.New("function parameters must be a list")
			}

			parameterList := parametersValue.Data.([]Value)
			var parameters []string
			for _, parameter := range parameterList {
				if parameter.Type != Symbol {
					return Value{}, errors.New("function parameters must be symbols")
				}
				parameters = append(parameters, parameter.Data.(string))
			}

			body := args[1]

			function := func(
				callArgs []Value,
				callEnv *Environment,
			) (Value, error) {
				if len(callArgs) != len(parameters) {
					return Value{}, errors.New("incorrect number of arguments")
				}
				innerEnv := NewEnv(env)
				for index, parameter := range parameters {
					if callArgs[index].Type == List {
						innerEnv.Set(parameter, Value{
							Type: Lazy,
							Data: LazyData{
								Expression:  callArgs[index],
								Environment: callEnv,
							},
						})
					} else {
						innerEnv.Set(parameter, callArgs[index])
					}
				}
				return Value{
					Type: Lazy,
					Data: LazyData{
						Expression:  body,
						Environment: innerEnv,
					},
				}, nil
			}
			return Value{
				Type: Function,
				Data: function,
			}, nil
		},
	})

	env.Set("procedure", Value{
		Type: Function,
		Data: func(
			args []Value,
			env *Environment,
		) (Value, error) {
			// [procedure [parameters] body]
			if len(args) != 2 {
				return Value{}, errors.New("procedure requires 2 arguments")
			}

			parametersValue := args[0]
			if parametersValue.Type != List {
				return Value{}, errors.New("procedure parameters must be a list")
			}

			parameterList := parametersValue.Data.([]Value)
			var parameters []string
			for _, parameter := range parameterList {
				if parameter.Type != Symbol {
					return Value{}, errors.New("procedure parameters must be symbols")
				}
				parameters = append(parameters, parameter.Data.(string))
			}

			body := args[1]

			procedure := func(
				callArgs []Value,
				callEnv *Environment,
			) (Value, error) {
				if len(callArgs) != len(parameters) {
					return Value{}, errors.New("incorrect number of arguments")
				}
				innerEnv := NewEnv(env)
				for index, parameter := range parameters {
					callArg, err := EvaluateUntilConcrete(
						callArgs[index],
						callEnv,
					)
					if err != nil {
						return Value{}, err
					}
					innerEnv.Set(parameter, callArg)
				}
				return Evaluate(body, innerEnv)
			}

			return Value{
				Type: Procedure,
				Data: procedure,
			}, nil
		},
	})

	env.Set("if", Value{
		Type: Function,
		Data: func(
			args []Value,
			env *Environment,
		) (
			Value,
			error,
		) {
			// [if condition then else?]
			if len(args) != 2 && len(args) != 3 {
				return Value{}, errors.New("if requires 2 or 3 arguments")
			}
			condition, err := EvaluateUntilConcrete(args[0], env)
			if err != nil {
				return Value{}, err
			}
			thenExpression := args[1]
			var elseExpression Value
			if len(args) == 3 {
				elseExpression = args[2]
			} else {
				elseExpression = Value{
					Type: Option,
					Data: OptionValue{
						Some: false,
					},
				}
			}
			// Only a bool or Option with Some == false is considered false.
			if (condition.Type == Bool && !condition.Data.(bool)) || (condition.Type == Option && !condition.Data.(OptionValue).Some) {
				return Evaluate(elseExpression, env)
			}
			return Evaluate(thenExpression, env)
		},
	})

	env.Set("match", Value{
		Type: Function,
		Data: func(
			args []Value,
			env *Environment,
		) (
			Value,
			error,
		) {
			// [match expression [pattern result] [pattern result] ...]
			if len(args) < 2 {
				return Value{}, errors.New("match requires an expression and at least one clause")
			}
			matchValue, err := EvaluateUntilConcrete(args[0], env)
			if err != nil {
				return Value{}, err
			}

			var wildcard Value
			for i := 1; i < len(args); i++ {
				clause := args[i]
				if clause.Type != List {
					return Value{}, errors.New("each clause in match must be a list")
				}
				clauseList := clause.Data.([]Value)
				if len(clauseList) != 2 {
					return Value{}, errors.New("each clause in match must have exactly 2 elements")
				}
				pattern := clauseList[0]
				resultExpression := clauseList[1]

				if pattern.Type == Symbol && pattern.Data.(string) == "_" {
					// Grab wildcard pattern as fallback and store for later.
					wildcard = resultExpression
					continue
				}

				patternValue, err := Evaluate(pattern, env)
				if err != nil {
					return Value{}, err
				}
				if valueEqual(matchValue, patternValue, env) {
					return Evaluate(resultExpression, env)
				}
			}

			if wildcard.Type != Unknown {
				return Evaluate(wildcard, env)
			}

			return Value{
				Type: Option,
				Data: OptionValue{
					Some: false,
				},
			}, nil
		},
	})

	// Option constructors.
	env.Set("some", Value{
		Type: Function,
		Data: func(
			args []Value,
			_ *Environment,
		) (
			Value,
			error,
		) {
			if len(args) != 1 {
				return Value{}, errors.New("Some requires exactly one argument")
			}
			return Value{
				Type: Option,
				Data: OptionValue{
					Some:  true,
					Value: args[0],
				},
			}, nil
		},
	})
	env.Set("none", Value{
		Type: Option,
		Data: OptionValue{
			Some: false,
		},
	})

	// Arithmetic operators.
	env.Set("int-add", addInts)
	env.Set("int-subtract", subtractInts)
}
