package logic

import (
	"fmt"
	"math"
	"os"

	boolutil "github.com/dterbah/go-logic/src/utils"
	"github.com/dterbah/gods/set"
	comparator "github.com/dterbah/gods/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
)

/*
Struct that will execute the main program
*/
type Runner struct {
	input         string
	displayHelp   bool
	generateGraph bool
	truthTable    bool
}

func NewRunner(input string, displayHelp bool, generateGraph bool, generateTruthTable bool) *Runner {
	return &Runner{input: input,
		generateGraph: generateGraph,
		truthTable:    generateTruthTable,
		displayHelp:   displayHelp,
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
	variablesTokens := tokens.Filter(func(element Token) bool {
		return element.Is(VAR)
	})

	variables := set.New(comparator.StringComparator)
	variablesTokens.ForEach(func(element Token, index int) {
		variables.Add(element.Value)
	})

	variablesTokens.ForEach(func(element Token, index int) {
		variables.Add(element.Value)
	})

	result, err := parser.Parse()

	if err != nil {
		fmt.Println(err)
		return
	}

	variables2 := make(map[string]bool)
	variables2["a"] = false
	variables2["b"] = false
	fmt.Println(result.Eval(variables2))

	if runner.truthTable {
		runner.generateTruthTable(result, *variables)
	}
}

func (runner Runner) help() {

}

func (runner Runner) generateTruthTable(expr Expression, variables set.Set[string]) {
	nbrVariables := variables.Size()
	iterations := int(math.Pow(2, float64(nbrVariables)))

	data := [][]string{}
	headers := append(variables.ToArray(), runner.input)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)

	for i := 0; i < iterations; i++ {
		variablesMap := make(map[string]bool)
		tableRow := []string{}
		// read all first bits of i
		variables.ForEach(func(element string, index int) {
			mask := int(math.Pow(2, float64(index)))
			value := (i & mask) != 0
			variablesMap[element] = (i & mask) != 0

			tableRow = append(tableRow, boolutil.BoolToString(value))
		})

		result := expr.Eval(variablesMap)
		tableRow = append(tableRow, boolutil.BoolToString(result))

		data = append(data, tableRow)
	}

	// display the truth table
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
