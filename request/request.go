package request

type Request struct {
	Name     string
	Children []*Request
}

func NewRequest(body string) (*Request, error) {
	parser := newRequestParser(body)
	return parser.Parse()
}
