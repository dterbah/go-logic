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
	generateTruthTable := flag.Bool("t", true, "Generate truth table")
	simplifyExpression := flag.Bool("s", false, "Simplify the expression")
	printHelp := flag.Bool("h", false, "Display list of operator")
	flag.Parse()

	if *logicExpression == "" {
		fmt.Println("The -e option is required.")
		flag.Usage()
		os.Exit(1)
	}

	runner := logic.NewRunner(*logicExpression, *printHelp, *generateGraph, *generateTruthTable, *simplifyExpression)
	runner.Run()
}
