package analyzer

import (
	"fmt"
	"strings"

	"text/scanner"

	"github.com/dalloriam/graphql-tools/analyzer/nodes"
)

type SchemaParser struct {
	schema string
	stream *scanner.Scanner

	currentTok  rune
	currentText string

	nextTok  rune
	nextText string
}

func NewSchemaParser(schema string) *SchemaParser {
	reader := strings.NewReader(schema)
	scan := scanner.Scanner{}
	scan.Init(reader)
	scan.Whitespace = 1<<'\t' | 1<<'\r' | 1<<' '

	s := &SchemaParser{
		schema: schema,
		stream: &scan,
	}

	s.Next()
	return s
}

func (p *SchemaParser) Next() {
	p.currentText = p.nextText
	p.nextTok = p.currentTok

	p.nextTok = p.stream.Scan()
	p.nextText = p.stream.TokenText()

	//fmt.Println("currentTok: ", p.currentText)
	//fmt.Println("nextTok: ", p.nextText)
}

func (p *SchemaParser) accept(body string) (ok bool) {
	if ok = (p.nextText == body); ok {
		p.Next()
	}
	return
}

func (p *SchemaParser) expect(body string) error {
	if !p.accept(body) {
		return fmt.Errorf("expected '%s', got '%s'", body, p.nextText)
	}
	return nil
}

func (p *SchemaParser) Parse() error {
	for p.nextTok != scanner.EOF && p.currentTok != scanner.EOF {
		switch p.currentText {
		case "#":
			p.parseComment()
		case "input", "type":
			_, err := p.parseBlock()
			if err != nil {
				return err
			}
		}
		p.Next()
	}
	return nil
}

func (p *SchemaParser) parseBlock() (*nodes.Block, error) {
	block := &nodes.Block{}
	block.Type = p.currentText
	p.Next()

	block.Name = p.currentText

	fmt.Println()
	fmt.Println(block.Type, ": ", block.Name)

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

func (p *SchemaParser) parseField() (*nodes.Field, error) {
	var line []string

	field := &nodes.Field{}

	p.Next()

	field.Name = p.currentText

	if p.accept("(") {
		for !p.accept(")") {
			p.Next()
			// TODO: Parse parameters
			//fmt.Println("PARAM: ", p.currentText)
		}
	}

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

	fmt.Printf(
		"    - %s -> <type='%s' repeated='%v' nullable='%v'>\n",
		field.Name,
		field.Type.Name,
		field.Type.Repetition == nodes.MULTIPLE,
		field.Type.Nullable,
	)

	return field, nil
}

func (p *SchemaParser) parseType() (*nodes.Type, error) {
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

func (p *SchemaParser) parseComment() {
	for p.currentText != "\n" {
		p.Next()
	}
	p.Next()
}
