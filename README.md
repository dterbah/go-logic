## Go Logic !

![CI](https://github.com/dterbah/go-logic/actions/workflows/go-test.yml/badge.svg)
[![codecov](https://codecov.io/gh/dterbah/go-logic/branch/main/graph/badge.svg)](https://codecov.io/gh/dterbah/go-logic)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=dterbah_go-logic&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=dterbah_go-logic)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=dterbah_go-logic&metric=bugs)](https://sonarcloud.io/summary/new_code?id=dterbah_go-logic)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=dterbah_go-logic&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=dterbah_go-logic)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=dterbah_go-logic&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=dterbah_go-logic)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=dterbah_go-logic&metric=sqale_index)](https://sonarcloud.io/summary/new_code?id=dterbah_go-logic)

Be ready to solve your boolean expressions with Go Logic ! This simple CLI is used to take a boolean expression and solve it with all possible combination. You'll find a short guide to help you get the most out of this application.

<img src="./assets/logo.webp" width="150" />

### Installation

To install locally this tool, you can run the following command :

```bash
go install github.com/dterbah/go-logic
```

Once this command is finished, you can directly use it like this :

```bash
go-logic ...
```

### Syntax of the boolean expression

Here is an overview of the syntax for the different boolean operator

| Operator name | Description       | Syntax in Go Logic | Usages             |
| ------------- | ----------------- | ------------------ | ------------------ |
| NOT           | Negation operator | !                  | !a                 |
| OR            | Or operator       | \|, v              | avb, a\|!b         |
| AND           | And operator      | ^,&,.              | a^b, (!a.b).c, c&a |
| XOR           | Xor operator      | +                  | a+b                |
| IMPLIES       | Implie operator   | ->                 | a->b, b->!(a^v)    |

### CLI usage and options

You can use this command using the command `go-logic`. There is multiple options you
can use to have different outputs

| Command | Description                                         | Usage                | Default value | Required |
| ------- | --------------------------------------------------- | -------------------- | ------------- | -------- |
| -e      | Define the expression you want to analyze           | go-logic -e="a+b"    | None          | ✅       |
| -t      | Create and output the truth table of the expression | go-logic -e="a" -t   | True          | ❌       |
| -g      | Create a DOT graph of your expression               | go-logic -e="a^1" -g | False         | ❌       |
| -s      | Simplify the current expression                     | go-logic -e="a+b" -s | False         | ❌       |
