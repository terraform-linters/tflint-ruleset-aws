package utils

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"text/template"
)

var heading = regexp.MustCompile("(^[A-Za-z])|_([A-Za-z])")

// ToCamel converts a string to CamelCase
func ToCamel(str string) string {
	exceptions := map[string]string{
		"ami":         "AMI",
		"db":          "DB",
		"alb":         "ALB",
		"elb":         "ELB",
		"elasticache": "ElastiCache",
		"iam":         "IAM",
		"ip":          "IP",
		"sql":         "SQL",
		"vm":          "VM",
		"os":          "OS",
		"id":          "ID",
		"tls":         "TLS",
		"oauth":       "OAuth",
		"ttl":         "TTL",
		"api":         "API",
		"uri":         "URI",
		"url":         "URL",
		"http":        "HTTP",
		"ui":          "UI",
		"dns":         "DNS",
		"ssh":         "SSH",
		"acl":         "ACL",
		"xss":         "XSS",
		"docdb":       "DocDB",
		"dynamodb":    "DynamoDB",
		"memorydb":    "MemoryDB",
	}
	parts := strings.Split(str, "_")
	replaced := make([]string, len(parts))
	for i, s := range parts {
		conv, ok := exceptions[s]
		if ok {
			replaced[i] = conv
		} else {
			replaced[i] = s
		}
	}
	str = strings.Join(replaced, "_")

	return heading.ReplaceAllStringFunc(str, func(s string) string {
		return strings.ToUpper(strings.Replace(s, "_", "", -1))
	})
}

// GenerateFile generates a new file from the passed template and metadata
func GenerateFile(fileName string, tmplName string, meta interface{}) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	tmpl := template.Must(template.ParseFiles(tmplName))
	err = tmpl.Execute(file, meta)
	if err != nil {
		panic(err)
	}
}

// GenerateFileWithLogs generates a new file from the passed template and metadata
// The difference from GenerateFile function is to output logs
func GenerateFileWithLogs(fileName string, tmplName string, meta interface{}) {
	GenerateFile(fileName, tmplName, meta)
	fmt.Printf("Created: %s\n", fileName)
}

// CleanDir removes generated .go files from dir that are no longer in the
// generated list. Only files containing "DO NOT EDIT" are considered generated.
func CleanDir(dir string, generated []string) {
	entries, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		panic(err)
	}
	for _, path := range entries {
		if slices.Contains(generated, filepath.Base(path)) {
			continue
		}
		if !isGenerated(path) {
			continue
		}
		fmt.Printf("Removing stale file: %s\n", path)
		if err := os.Remove(path); err != nil {
			panic(err)
		}
	}
}

func isGenerated(path string) bool {
	f, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments|parser.PackageClauseOnly)
	if err != nil {
		panic(err)
	}
	return ast.IsGenerated(f)
}
