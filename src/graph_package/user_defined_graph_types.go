package graph_package

// These are the types that might be changed by the user
type VertexIdType string

type VertexValue struct {
	Value int
}

type VertexState struct {
	neighbors map[VertexIdType]bool
}

type EdgeValue struct {
	Cost int
}

type PregelMessage struct {
	OriginVertexId VertexIdType
	Value          int
}
