package graphkit

// This map used for O(1) insertions & lookups was proving a bit clunky.
// This graph struct should make it a bit friendlier to use.
type graph struct {
	internalMap map[string]map[string]map[string]struct{}
}

func newGraph() *graph {
	return &graph{make(map[string]map[string]map[string]struct{})}
}

func (s *graph) AddEdge(fromNode, toNode, edgeName string) {
	fromMap, ok := s.internalMap[fromNode]
	if !ok {
		s.internalMap[fromNode] = make(map[string]map[string]struct{})
		fromMap = s.internalMap[fromNode]
	}

	toMap, ok := fromMap[toNode]
	if !ok {
		fromMap[toNode] = make(map[string]struct{})
		toMap = fromMap[toNode]
	}

	toMap[edgeName] = struct{}{}
}

func (s *graph) HasEdge(fromNode, toNode, edgeName string) bool {
	fromMap, ok := s.internalMap[fromNode]
	if !ok {
		return false
	}

	toMap, ok := fromMap[toNode]
	if !ok {
		return false
	}

	_, ok = toMap[edgeName]
	return ok
}
