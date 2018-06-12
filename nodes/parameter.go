package nodes

// Parameter represents a parameter to a graphql type
type Parameter struct {
	Name         string `json:"name,omitempty"`
	Type         *Type  `json:"type,omitempty"`
	DefaultValue string `json:"default_value,omitempty"`
}
