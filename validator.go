package graphkit

import (
	"fmt"

	"github.com/dalloriam/graphkit/nodes"
)

type queryValidator struct {
	schema   *Schema
	visitMap *graph
}

func newQueryValidator(schema *Schema) *queryValidator {
	return &queryValidator{schema, newGraph()}
}

func (v *queryValidator) walk(currentTree *request, currentBlock *nodes.Block) error {
	for _, blockItm := range currentBlock.Fields {
		if blockItm.Name == currentTree.Name {
			if len(currentTree.Children) > 0 {
				fmt.Printf("visiting gql type '%s' from field '%s'\n", blockItm.Type.Name, blockItm.Name)

				if v.visitMap.HasEdge(currentBlock.Type, blockItm.Type.Name, blockItm.Name) {
					return fmt.Errorf("cycle detected: '%s.%s' -> '%s'", currentBlock.Type, blockItm.Name, blockItm.Type.Name)
				}
				v.visitMap.AddEdge(currentBlock.Type, blockItm.Type.Name, blockItm.Name)

				newBlock, err := v.schema.ResolveType(blockItm.Type.Name)
				if err != nil {
					return err
				}

				for _, child := range currentTree.Children {
					if err := v.walk(child, newBlock); err != nil {
						return err
					}
				}
			}
			return nil
		}
	}
	return fmt.Errorf("field '%s' not found", currentTree.Name)
}

func (v *queryValidator) Traverse(req *request) error {
	var err error

	var currentBlock *nodes.Block

	if req.Name == "mutation" {
		currentBlock, err = v.schema.ResolveType(v.schema.RootMutation)
	} else if req.Name == "query" {
		currentBlock, err = v.schema.ResolveType(v.schema.RootQuery)
	}
	if err != nil {
		return err
	}

	for _, child := range req.Children {
		if err := v.walk(child, currentBlock); err != nil {
			return err
		}
	}
	return nil
}

// ValidateQuery validates a GraphQL query against the provided schema.
// It checks for unknown fields as well as for possible malicious queries.
func ValidateQuery(query string, schema *Schema) error {
	validator := newQueryValidator(schema)

	req, err := newRequest(query)
	if err != nil {
		return err
	}

	return validator.Traverse(req)
}
