package language

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Repl() {
	env := NewEnv(nil)
	AddBuiltins(env)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Type 'exit' to quit.")
	for {
		fmt.Print(">> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if strings.TrimSpace(line) == "exit" {
			break
		}
		if strings.TrimSpace(line) == "" {
			continue
		}

		expression, err := Parse(line, "<repl>", nil)
		if err != nil {
			fmt.Println("Parse error:", err)
			continue
		}
		result, err := Evaluate(expression, env)
		if err != nil {
			fmt.Println("Eval error:", err)
			continue
		}
		fmt.Println("Result:", toString(result))
	}
}
