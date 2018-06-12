package graphkit

import (
	"fmt"
	"strings"
	"text/scanner"
)

type requestParser struct {
	request string
	stream  *scanner.Scanner

	currentTok  rune
	currentText string

	nextTok  rune
	nextText string
}

func newRequestParser(requestText string) *requestParser {
	reader := strings.NewReader(strings.TrimSpace(requestText))

	scan := scanner.Scanner{}
	scan.Init(reader)

	r := &requestParser{
		request: requestText,
		stream:  &scan,
	}

	r.Next()
	r.Next()
	return r
}

func (r *requestParser) Next() {
	r.currentText = r.nextText
	r.nextTok = r.currentTok

	r.nextTok = r.stream.Scan()
	r.nextText = r.stream.TokenText()
}

func (r *requestParser) accept(body string) (ok bool) {
	if ok = (r.nextText == body); ok {
		r.Next()
	}
	return
}

func (r *requestParser) expect(body string) error {
	if !r.accept(body) {
		return fmt.Errorf("expected '%s', got '%s'", body, r.nextText)
	}
	return nil
}

func (r *requestParser) Parse() (*request, error) {
	var req *request

	if r.currentText == "query" {
		req = &request{Name: "query"}
	} else if r.currentText == "mutation" {
		req = &request{Name: "mutation"}
	} else {
		return nil, fmt.Errorf("invalid root keyword: %s", r.currentText)
	}

	if r.accept("(") {
		r.parseParameters()
	}

	r.expect("{")
	children, err := r.parseBlock()
	if err != nil {
		return nil, err
	}

	req.Children = children

	return req, nil
}

func (r *requestParser) parseBlock() ([]*request, error) {
	block := []*request{}

	for !r.accept("}") && r.nextTok != scanner.EOF {
		r.Next()
		name := r.currentText

		if r.accept("(") {
			r.parseParameters()
		}

		children := []*request{}

		if r.accept("{") {
			var err error
			children, err = r.parseBlock()
			if err != nil {
				return nil, err
			}
		}
		block = append(block, &request{Name: name, Children: children})
	}
	return block, nil
}

func (r *requestParser) parseParameters() {
	for !r.accept(")") && r.nextTok != scanner.EOF {
		r.Next()
	}
}
