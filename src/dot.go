package logic

import "strings"

func GenerateDot(expression Expression) string {
	var builder strings.Builder
	builder.WriteString("digraph G {\n")
	expression.ToDot(&builder, "")
	builder.WriteString("}\n")
	return builder.String()
}
