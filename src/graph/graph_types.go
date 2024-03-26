package graph

type Vertex struct {
	Id    int
	Value VertexValue
	Edges []Edge
}

type Edge struct {
	To    int
	Value EdgeValue
}

type VertexValue struct {
	Value int
}

type EdgeValue struct {
	Cost int
}
