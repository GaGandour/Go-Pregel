package graph_package

func (vertex *Vertex) SuperStep() {
	vertex.Compute()
	if vertex.numSuperSteps == 0 {
		vertex.Activate()
	}
	vertex.numSuperSteps++
	vertex.ReceivedMessages = []PregelMessage{}
}

func (vertex *Vertex) GetValue() VertexValue {
	return vertex.Value
}

func (vertex *Vertex) SetValue(value VertexValue) {
	vertex.Value = value
}

func (vertex *Vertex) GetOutEdges() map[VertexIdType]*Edge {
	return vertex.Edges
}

func (vertex *Vertex) PrepareMessageToVertex(vertexId VertexIdType, message PregelMessage) {
	vertex.messageMutex.Lock()
	defer vertex.messageMutex.Unlock()
	vertex.MessagesToSend[vertexId] = append(vertex.MessagesToSend[vertexId], message)
}

func (vertex *Vertex) ReceiveMessage(message PregelMessage) {
	vertex.messageMutex.Lock()
	vertex.ReceivedMessages = append(vertex.ReceivedMessages, message)
	vertex.VotedToHalt = false
	vertex.messageMutex.Unlock()
}

func (vertex *Vertex) VoteToHalt() {
	vertex.VotedToHalt = true
}

func (vertex *Vertex) Activate() {
	vertex.VotedToHalt = false
}
