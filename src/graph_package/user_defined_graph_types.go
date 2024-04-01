package graph_package

// These are the types that might be changed by the user
type VertexValue struct {
	Value int
}

type EdgeValue struct {
	Cost int
}

type PregelMessage struct {
	OriginVertexId VertexIdType
	Value          int
}
