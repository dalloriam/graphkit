package nodes

type typeRepetition int

const (
	SINGLE typeRepetition = iota
	MULTIPLE
)

type Type struct {
	Name       string         `json:"name,omitempty"`
	Nullable   bool           `json:"nullable,omitempty"`
	Repetition typeRepetition `json:"repetition,omitempty"`
}
