package graphkit

import (
	"fmt"
	"strings"

	"text/scanner"

	"github.com/dalloriam/graphkit/nodes"
)

type schemaParser struct {
	schema string
	stream *scanner.Scanner

	currentTok  rune
	currentText string

	nextTok  rune
	nextText string
}

func newSchemaParser(schema string) *schemaParser {
	reader := strings.NewReader(schema)
	scan := scanner.Scanner{}
	scan.Init(reader)
	scan.Whitespace = 1<<'\t' | 1<<'\r' | 1<<' '

	s := &schemaParser{
		schema: schema,
		stream: &scan,
	}

	s.Next()
	return s
}

func (p *schemaParser) Next() {
	p.currentText = p.nextText
	p.nextTok = p.currentTok

	p.nextTok = p.stream.Scan()
	p.nextText = p.stream.TokenText()
}

func (p *schemaParser) accept(body string) (ok bool) {
	if ok = (p.nextText == body); ok {
		p.Next()
	}
	return
}

func (p *schemaParser) expect(body string) error {
	if !p.accept(body) {
		return fmt.Errorf("expected '%s', got '%s'", body, p.nextText)
	}
	return nil
}

func (p *schemaParser) Parse() (*Schema, error) {

	schema := &Schema{
		Types: make(map[string]*nodes.Block),
	}

	for p.nextTok != scanner.EOF && p.currentTok != scanner.EOF {
		if p.currentText == "#" {
			p.parseComment()
		}

		switch p.currentText {

		case "schema":
			blk, err := p.parseRoot()
			if err != nil {
				return nil, err
			}

			for _, f := range blk.Fields {
				if f.Name == "query" {
					schema.RootQuery = f.Type.Name
				} else if f.Name == "mutation" {
					schema.RootMutation = f.Type.Name
				}
			}

		case "input", "type":
			blk, err := p.parseBlock()
			if err != nil {
				return nil, err
			}
			schema.Types[blk.Name] = blk
		}
		p.Next()
	}
	return schema, nil
}

func (p *schemaParser) parseRoot() (*nodes.Block, error) {
	block := &nodes.Block{}
	block.Type = p.currentText

	if err := p.expect("{"); err != nil {
		return nil, err
	}

	var fields []*nodes.Field

	for !p.accept("}") {
		if p.nextText == "\n" {
			p.Next()
			continue
		}

		if p.nextText == "#" {
			p.parseComment()
		}

		f, err := p.parseField()
		if err != nil {
			return nil, err
		}

		fields = append(fields, f)
	}
	block.Fields = fields

	return block, nil
}

func (p *schemaParser) parseBlock() (*nodes.Block, error) {
	block := &nodes.Block{}
	block.Type = p.currentText
	p.Next()

	block.Name = p.currentText

	if err := p.expect("{"); err != nil {
		return nil, err
	}

	var fields []*nodes.Field

	for !p.accept("}") {
		if p.nextText == "\n" {
			p.Next()
			continue
		}

		if p.nextText == "#" {
			p.parseComment()
		}

		f, err := p.parseField()
		if err != nil {
			return nil, err
		}

		fields = append(fields, f)
	}
	block.Fields = fields

	return block, nil
}

func (p *schemaParser) parseField() (*nodes.Field, error) {
	var line []string

	field := &nodes.Field{}

	p.Next()

	field.Name = p.currentText

	var parameters []*nodes.Parameter

	if p.accept("(") {
		for !p.accept(")") {
			p.Next()
			param, err := p.parseParameter()
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, param)
		}
	}

	field.Parameters = parameters

	if err := p.expect(":"); err != nil {
		return nil, err
	}

	t, err := p.parseType()
	if err != nil {
		return nil, err
	}

	field.Type = t

	for p.currentText != "\n" && p.nextText != "}" {
		line = append(line, p.currentText)
		p.Next()
	}

	return field, nil
}

func (p *schemaParser) parseParameter() (*nodes.Parameter, error) {
	param := &nodes.Parameter{}
	param.Name = p.currentText

	if err := p.expect(":"); err != nil {
		return nil, err
	}

	t, err := p.parseType()
	if err != nil {
		return nil, err
	}

	param.Type = t

	if p.accept("=") {
		p.Next()
		param.DefaultValue = p.currentText
	}

	if p.nextText == "," {
		p.Next()
	}

	return param, nil
}

func (p *schemaParser) parseType() (*nodes.Type, error) {
	isList := p.accept("[")
	p.Next()

	t := &nodes.Type{}
	t.Name = p.currentText

	t.Nullable = !p.accept("!")

	if isList {
		t.Repetition = nodes.MULTIPLE
		if err := p.expect("]"); err != nil {
			return nil, err
		}
	}

	return t, nil
}

func (p *schemaParser) parseComment() {
	for p.currentText != "\n" {
		fmt.Printf("Skipping %s\n", p.currentText)
		p.Next()
	}
	p.Next()
}
