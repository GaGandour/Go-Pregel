package graph

import "sync"

type VertexIdType string

type Graph struct {
	Vertexes map[VertexIdType]*Vertex
}

type Vertex struct {
	Id               VertexIdType
	Value            VertexValue
	Edges            map[VertexIdType]Edge
	receivedMessages []PregelMessage
	messageMutex     sync.Mutex
	messagesToSend   map[VertexIdType][]PregelMessage
	votedToHalt      bool
}

type Edge struct {
	To    VertexIdType
	Value EdgeValue
}
