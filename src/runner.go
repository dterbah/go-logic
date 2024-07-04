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
	input              string
	displayHelp        bool
	generateGraph      bool
	truthTable         bool
	simplifyExpression bool
}

func NewRunner(input string, displayHelp bool, generateGraph bool, generateTruthTable bool, simplifyExpression bool) *Runner {
	return &Runner{input: input,
		generateGraph:      generateGraph,
		truthTable:         generateTruthTable,
		displayHelp:        displayHelp,
		simplifyExpression: simplifyExpression,
	}
}

/*
Run the program
*/
func (runner Runner) Run() {
	var simplifiedExpr Expression
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

	if runner.simplifyExpression {
		simplifiedExpr = result.Simplify()
	}

	if runner.truthTable {
		runner.generateTruthTable(result, *variables, simplifiedExpr)
	}

	if runner.generateGraph {
		fmt.Println("🚀 Dot Graph is being generated ...")
		var graph string
		if simplifiedExpr != nil {
			graph = GenerateDot(simplifiedExpr)
		} else {
			graph = GenerateDot(result)
		}

		err := os.WriteFile("graph.dot", []byte(graph), 0644)
		if err != nil {
			fmt.Println("Error when generating your dot graph")
		} else {
			fmt.Println("✅ Your graph is ready to be displayed !")
		}
	}
}

func (runner Runner) help() {

}

func (runner Runner) generateTruthTable(expr Expression, variables set.Set[string], simplifiedExpr Expression) {
	nbrVariables := variables.Size()
	iterations := int(math.Pow(2, float64(nbrVariables)))

	data := [][]string{}
	headers := append(variables.ToArray(), runner.input)
	if simplifiedExpr != nil {
		headers = append(headers, fmt.Sprintf("Simplified : %s", simplifiedExpr.String()))
	}

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

		if simplifiedExpr != nil {
			tableRow = append(tableRow, boolutil.BoolToString(result))
		}

		data = append(data, tableRow)
	}

	// display the truth table
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
