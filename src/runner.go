package logic

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

/*
Struct that will execute the main program
*/
type Runner struct {
	input              string
	displayHelp        bool
	generateGraph      bool
	generateTruthTable bool
}

func NewRunner(input string, displayHelp bool, generateGraph bool, generateTruthTable bool) *Runner {
	return &Runner{input: input,
		generateGraph:      generateGraph,
		generateTruthTable: generateTruthTable,
		displayHelp:        displayHelp,
	}
}

/*
Run the program
*/
func (runner Runner) Run() {
	if runner.displayHelp {
		runner.help()
	}

	lexer := NewLexer(runner.input)
	tokens, err := lexer.Tokenize()

	if err != nil {
		logrus.Error(err)
		return
	}

	parser := NewParser(tokens)

	result, err := parser.Parse()

	if err != nil {
		fmt.Println(err)
		return
	}

	variables := make(map[string]bool)
	variables["a"] = false
	variables["b"] = false
	fmt.Println(result.Eval(variables))
}

func (runner Runner) help() {

}
