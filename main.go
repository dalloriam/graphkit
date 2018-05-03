package main

import (
	"github.com/dalloriam/graphql-tools/analyzer"
)

func main() {
	schema, err := analyzer.LoadGraphQLSchema("./schema")
	if err != nil {
		panic(err)
	}

	parser := analyzer.NewSchemaParser(schema)
	if err := parser.Parse(); err != nil {
		panic(err)
	}
}
