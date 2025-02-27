package language

import (
	"testing"
)

// evaluateOrFail is a helper that parses and evaluates a source string in the given environment. It fails the test if there is an error.
func evaluateOrFail(
	input string,
	env *Environment,
	t *testing.T,
) Value {
	expression, err := Parse(input, "<test>", nil)
	if err != nil {
		t.Fatalf("Parse error in input %q: %v", input, err)
	}
	result, err := EvaluateUntilConcrete(expression, env)
	if err != nil {
		t.Fatalf("Eval error in input %q: %v", input, err)
	}
	return result
}

func TestFibonacci(
	t *testing.T,
) {
	env := NewEnv(nil)
	AddBuiltins(env)

	fibDef := `
		[define calc-fib
			[function [n]
				[match n
					[0 0]
					[1 1]
					[_
						[int-add
							[calc-fib
								[int-subtract n 1]
							]
							[calc-fib
								[int-subtract n 2]
							]
						]
					]
				]
			]
		]
	`
	_ = evaluateOrFail(fibDef, env, t)

	tests := []struct {
		input    string
		expected string
	}{
		{"[calc-fib 0]", "int<0>"},
		{"[calc-fib 1]", "int<1>"},
		{"[calc-fib 2]", "int<1>"},
		{"[calc-fib 3]", "int<2>"},
		{"[calc-fib 4]", "int<3>"},
		{"[calc-fib 5]", "int<5>"},
		{"[calc-fib 6]", "int<8>"},
		{"[calc-fib 7]", "int<13>"},
	}

	for _, test := range tests {
		result := evaluateOrFail(test.input, env, t)
		resultString := valueToString(result)
		if resultString != test.expected {
			t.Errorf("For %s: expected %s, got %s", test.input, test.expected, resultString)
		}
	}
}
