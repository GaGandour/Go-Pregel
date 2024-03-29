package graph

// These are the types that might be changed by the user
type VertexValue struct {
	Value int
}

type EdgeValue struct {
	Cost int
}

type PregelMessage struct {
	VertexId VertexIdType
}
