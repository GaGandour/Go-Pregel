package graph

type VertexIdType string

type Graph struct {
	Vertexes map[VertexIdType]Vertex
}

type Vertex struct {
	Id    VertexIdType
	Value VertexValue
	Edges map[VertexIdType]Edge
}

type Edge struct {
	To    VertexIdType
	Value EdgeValue
}

type VertexValue struct {
	Value int
}

type EdgeValue struct {
	Cost int
}
