package graphkit

type request struct {
	Name     string
	Children []*request
}

func newRequest(body string) (*request, error) {
	parser := newRequestParser(body)
	return parser.Parse()
}
