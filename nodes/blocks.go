package nodes

type Block struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`

	Fields []*Field `json:"fields,omitempty"`
}
