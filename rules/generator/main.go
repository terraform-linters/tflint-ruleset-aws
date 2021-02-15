package main

import (
	"bufio"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	utils "github.com/terraform-linters/tflint-ruleset-aws/rules/generator-utils"
)

type metadata struct {
	RuleName   string
	RuleNameCC string
}

func main() {
	buf := bufio.NewReader(os.Stdin)
	fmt.Print("Rule name? (e.g. aws_instance_invalid_type): ")
	ruleName, err := buf.ReadString('\n')
	if err != nil {
		panic(err)
	}
	ruleName = strings.Trim(ruleName, "\n")

	meta := &metadata{RuleNameCC: utils.ToCamel(ruleName), RuleName: ruleName}

	utils.GenerateFileWithLogs(fmt.Sprintf("rules/%s.go", ruleName), "rules/rule.go.tmpl", meta)
	utils.GenerateFileWithLogs(fmt.Sprintf("rules/%s_test.go", ruleName), "rules/rule_test.go.tmpl", meta)
	utils.GenerateFileWithLogs(fmt.Sprintf("docs/rules/%s.md", ruleName), "docs/rules/rule.md.tmpl", meta)

	src, err := ioutil.ReadFile("rules/provider.go")
	if err != nil {
		panic(err)
	}
	dstf, err := decorator.Parse(src)
	if err != nil {
		panic(err)
	}

	dst.Inspect(dstf, func(n dst.Node) bool {
		switch node := n.(type) {
		case *dst.CompositeLit:
			expr := &dst.CallExpr{
				Fun: &dst.Ident{
					Name: fmt.Sprintf("New%sRule", meta.RuleNameCC),
				},
			}
			expr.Decs.Before = dst.NewLine
			expr.Decs.After = dst.NewLine
			node.Elts = append(node.Elts, expr)
		}
		return true
	})

	fset, astf, err := decorator.RestoreFile(dstf)
	if err != nil {
		panic(err)
	}

	fp, err := os.OpenFile("rules/provider.go", os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	if err := format.Node(fp, fset, astf); err != nil {
		panic(err)
	}
	fmt.Println("Modified: rules/provider.go")

	fmt.Println(`
TODO:
1. Remove all "TODO" comments from generated files.
2. Write implementation of the rule.
3. Add a link to the generated documentation into docs/rules/README.md`)
}
