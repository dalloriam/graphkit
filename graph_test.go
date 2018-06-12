package graphkit

import (
	"testing"
)

func Test_graph_AddEdge(t *testing.T) {
	t.Run("correctly adds the edge to the map", func(t *testing.T) {
		g := newGraph()
		g.AddEdge("A", "B", "C")

		aBlock, ok := g.internalMap["A"]
		if !ok {
			t.Error("from node not inserted correctly")
		}

		bBlock, ok := aBlock["B"]
		if !ok {
			t.Error("to node not inserted correctly")
		}

		_, ok = bBlock["C"]
		if !ok {
			t.Errorf("edge name not inserted correctly")
		}
	})
}

func Test_graph_HasEdge(t *testing.T) {
	g := newGraph()
	g.AddEdge("A", "B", "C")

	t.Run("returns true on existing edge", func(t *testing.T) {
		if !g.HasEdge("A", "B", "C") {
			t.Error("returns false on existing edge")
		}
	})

	t.Run("returns false when FromNode doesn't exist", func(t *testing.T) {
		if g.HasEdge("X", "B", "C") {
			t.Error("returns true on inexistent from node")
		}
	})

	t.Run("returns false when ToNode doesn't exist", func(t *testing.T) {
		if g.HasEdge("A", "X", "C") {
			t.Error("returns true on inexistent to node")
		}
	})

	t.Run("returns false when edge name doesn't exist", func(t *testing.T) {
		if g.HasEdge("A", "B", "X") {
			t.Error("returns true on inexistent edge")
		}
	})
}
