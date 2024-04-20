package graph_package

import "sync"

// These are the types used in computing

type Graph struct {
	totalNumberOfVertexes int
	Vertexes              map[VertexIdType]*Vertex
}

type Vertex struct {
	Id                          VertexIdType
	Value                       VertexValue
	Edges                       map[VertexIdType]*Edge
	ReceivedMessagesInSuperStep map[int][]PregelMessage
	messageMutex                sync.Mutex
	MessagesToSend              map[VertexIdType][]PregelMessage
	VotedToHalt                 bool
	numSuperSteps               int
	HasSentMessages             bool
}

type Edge struct {
	To    VertexIdType
	Value EdgeValue
}

// These are the types used in communication

type CommunicationGraph struct {
	totalNumberOfVertexes int
	Vertexes              map[VertexIdType]CommunicationVertex
}

type CommunicationVertex struct {
	Id    VertexIdType
	Value VertexValue
	Edges map[VertexIdType]CommunicationEdge
}

type CommunicationEdge struct {
	To    VertexIdType
	Value EdgeValue
}
