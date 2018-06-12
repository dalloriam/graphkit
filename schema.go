package graphkit

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dalloriam/graphkit/nodes"
)

// Schema represents a parsed GraphQL schema.
type Schema struct {
	Types map[string]*nodes.Block `json:"types,omitempty"`

	RootQuery    string `json:"root_query,omitempty"`
	RootMutation string `json:"root_mutation,omitempty"`
}

// ResolveType resolves a type by name
func (s *Schema) ResolveType(typeName string) (*nodes.Block, error) {
	t, ok := s.Types[typeName]
	if !ok {
		return nil, fmt.Errorf("type %s does not exist", typeName)
	}

	return t, nil
}

// LoadSchema loads a GraphQL schema from a directory.
func LoadSchema(rootPath string) (*Schema, error) {
	raw, err := loadFromDisk(rootPath)
	if err != nil {
		return nil, err
	}

	parser := newSchemaParser(raw)
	parsed, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

type schemaLoader struct {
	LoadedSchema bytes.Buffer
}

func (s *schemaLoader) processFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	fname := info.Name()
	if strings.HasSuffix(fname, ".graphql") && !info.IsDir() {
		raw, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		if err := s.LoadedSchema.WriteByte('\n'); err != nil {
			return err
		}
		if _, err := s.LoadedSchema.Write(raw); err != nil {
			return err
		}
	}

	return nil
}

func loadFromDisk(rootPath string) (string, error) {
	loader := schemaLoader{}
	if err := filepath.Walk(rootPath, loader.processFile); err != nil {
		return "", err
	}

	return loader.LoadedSchema.String(), nil
}
