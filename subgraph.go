package main

type Subgraph struct {
	internalMap map[string]map[string]map[string]struct{}
}

func newSubgraph() *Subgraph {
	return &Subgraph{make(map[string]map[string]map[string]struct{})}
}

func (s *Subgraph) AddEdge(fromNode, toNode, edgeName string) {
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

func (s *Subgraph) HasEdge(fromNode, toNode, edgeName string) bool {
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
