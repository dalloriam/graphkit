package nodes

// Parameter represents a parameter to a graphql type
type Parameter struct {
	Name         string
	Type         *Type
	DefaultValue string
}
