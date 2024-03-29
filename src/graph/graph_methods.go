package graph

func (vertex *Vertex) SuperStep() {
	vertex.InterpretMessages()
	vertex.Compute()
	vertex.receivedMessages = []PregelMessage{}
}

func (vertex *Vertex) GetValue() VertexValue {
	return vertex.Value
}

func (vertex *Vertex) SetValue(value VertexValue) {
	vertex.Value = value
}

func (vertex *Vertex) GetOutEdges() []Edge {
	var edges []Edge
	for _, edge := range vertex.Edges {
		edges = append(edges, edge)
	}
	return edges
}

func (vertex *Vertex) PrepareMessageToVertex(vertexId VertexIdType, message PregelMessage) {
	vertex.messageMutex.Lock()
	defer vertex.messageMutex.Unlock()
	vertex.messagesToSend[vertexId] = append(vertex.messagesToSend[vertexId], message)
}

func (vertex *Vertex) VoteToHalt() {
	vertex.votedToHalt = true
}
