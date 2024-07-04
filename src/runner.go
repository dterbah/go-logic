package logic

import (
	"errors"
	"fmt"
	"math"
	"os"

	boolutil "github.com/dterbah/go-logic/src/utils"
	"github.com/dterbah/gods/set"
	comparator "github.com/dterbah/gods/utils"
	"github.com/goccy/go-graphviz"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
)

const DOT_GRAPH_IMAGE_PATH = "graph.png"

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

		err := exportDotGraph(graph)

		if err != nil {
			fmt.Println(err)
			fmt.Println("❌ Error during the generation of graph")
		} else {
			fmt.Println("✅ Graph created !")
		}
	}
}

func (runner Runner) help() {

}

func exportDotGraph(dotGraph string) error {
	graph, err := graphviz.ParseBytes([]byte(dotGraph))

	if err != nil {
		return errors.New("error during the export of graph")
	}

	g := graphviz.New()

	if err := g.RenderFilename(graph, graphviz.PNG, DOT_GRAPH_IMAGE_PATH); err != nil {
		return err
	}

	return nil
}

func createTruthTableData(expr Expression, variables set.Set[string], simplifiedExpr Expression) [][]string {
	nbrVariables := variables.Size()
	iterations := int(math.Pow(2, float64(nbrVariables)))

	data := [][]string{}

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

	return data
}

func (runner Runner) generateTruthTable(expr Expression, variables set.Set[string], simplifiedExpr Expression) {
	headers := append(variables.ToArray(), runner.input)
	if simplifiedExpr != nil {
		headers = append(headers, fmt.Sprintf("Simplified : %s", simplifiedExpr.String()))
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)

	data := createTruthTableData(expr, variables, simplifiedExpr)

	// display the truth table
	for _, v := range data {
		table.Append(v)
	}

	table.Render()
}
