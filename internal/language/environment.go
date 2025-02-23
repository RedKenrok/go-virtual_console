package language

import "errors"

// Environment represents an environment mapping symbols to Values.
type Environment struct {
	Values map[string]Value
	Outer  *Environment
}

// NewEnv creates a new environment with an optional outer (parent) environment.
func NewEnv(
	outer *Environment,
) *Environment {
	return &Environment{
		Values: make(map[string]Value),
		Outer:  outer,
	}
}

// Get retrieves a variable's value from the environment.
func (
	e *Environment,
) Get(
	key string,
) (
	Value,
	error,
) {
	if value, ok := e.Values[key]; ok {
		return value, nil
	}
	if e.Outer != nil {
		return e.Outer.Get(key)
	}
	return Value{}, errors.New("undefined symbol: " + key)
}

// Set assigns a value to a symbol in the environment.
func (
	e *Environment,
) Set(
	key string,
	value Value,
) {
	e.Values[key] = value
}
