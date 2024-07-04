package main

import (
	"flag"
	"fmt"
	"os"

	logic "github.com/dterbah/go-logic/src"
	"github.com/sirupsen/logrus"
)

// todo: Implement dot graph
func setupLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		ForceColors:            true,
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	setupLogger()
	logicExpression := flag.String("e", "", "Logic expression to evaluate")
	generateGraph := flag.Bool("g", false, "Generate the grap representation of the expression")
	generateTruthTable := flag.Bool("t", false, "Generate truth table")

	flag.Parse()

	if *logicExpression == "" {
		fmt.Println("The -e option is required.")
		flag.Usage()
		os.Exit(1)
	}

	lexer := logic.NewLexer(*logicExpression)
	tokens, err := lexer.Tokenize()

	if err != nil {
		logrus.Error(err)
		return
	}

	parser := logic.NewParser(tokens)

	result, err := parser.Parse()

	if err != nil {
		fmt.Println(err)
		return
	}

	variables := make(map[string]bool)
	variables["a"] = true
	variables["b"] = false
	fmt.Println(result.Eval(variables))

	fmt.Println(*generateGraph, *generateTruthTable)
}
