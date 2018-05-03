package analyzer

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type schemaLoader struct {
	LoadedSchema bytes.Buffer
}

func (s *schemaLoader) ProcessFile(path string, info os.FileInfo, err error) error {
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

func LoadGraphQLSchema(rootPath string) (string, error) {
	loader := schemaLoader{}
	if err := filepath.Walk(rootPath, loader.ProcessFile); err != nil {
		return "", err
	}

	return loader.LoadedSchema.String(), nil
}
