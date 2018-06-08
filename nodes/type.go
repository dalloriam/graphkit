package nodes

type typeRepetition int

const (
	SINGLE typeRepetition = iota
	MULTIPLE
)

type Type struct {
	Name       string
	Nullable   bool
	Repetition typeRepetition
}
