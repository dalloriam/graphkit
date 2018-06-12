package nodes

type Field struct {
	Name       string       `json:"name,omitempty"`
	Parameters []*Parameter `json:"parameters,omitempty"`
	Type       *Type        `json:"type,omitempty"`
}
