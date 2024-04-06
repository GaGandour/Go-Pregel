package graph_package

// These are the types that might be changed by the user

type VertexIdType string

type VertexValue struct {
	Value     int
	neighbors VertexIdSet
}

type EdgeValue struct {
	Cost int
}

type PregelMessage struct {
	OriginVertexId VertexIdType
	Value          int
}
