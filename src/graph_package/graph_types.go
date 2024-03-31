package graph_package

import "sync"

type VertexIdType string

type Graph struct {
	totalNumberOfVertexes int
	Vertexes              map[VertexIdType]Vertex
}

type Vertex struct {
	Id               VertexIdType
	Value            VertexValue
	Edges            map[VertexIdType]Edge
	ReceivedMessages []PregelMessage
	messageMutex     *sync.Mutex
	MessagesToSend   map[VertexIdType][]PregelMessage
	VotedToHalt      bool
}

type Edge struct {
	To    VertexIdType
	Value EdgeValue
}
