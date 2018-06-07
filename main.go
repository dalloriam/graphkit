package main

import (
	"fmt"

	"github.com/dalloriam/graphql-tools/analyzer"
)

func main() {

	schema, err := analyzer.LoadSchema("./schema")
	if err != nil {
		panic(err)
	}

	fmt.Println(schema.RootMutation)
}
