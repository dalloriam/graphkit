package main

import (
	"fmt"

	"github.com/dalloriam/graphql-tools/analyzer"
)

func main() {
	rawSchema, err := analyzer.LoadGraphQLSchema("./schema")
	if err != nil {
		panic(err)
	}

	parser := analyzer.NewSchemaParser(rawSchema)

	schema, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Println(schema.Types)
}
