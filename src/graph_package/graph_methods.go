package graph_package

func (vertex *Vertex) SuperStep() {
	vertex.InterpretMessages()
	if !vertex.VotedToHalt {
		vertex.Compute()
	}
	if vertex.numSuperSteps > 0 {
		vertex.VoteToHalt()
	}
	vertex.numSuperSteps++
	vertex.ReceivedMessages = []PregelMessage{}
}

func (vertex *Vertex) GetValue() VertexValue {
	return vertex.Value
}

func (vertex *Vertex) SetValue(value VertexValue) {
	vertex.Activate()
	vertex.Value = value
}

func (vertex *Vertex) GetOutEdges() []*Edge {
	var edges []*Edge
	for _, edge := range vertex.Edges {
		edges = append(edges, edge)
	}
	return edges
}

func (vertex *Vertex) PrepareMessageToVertex(vertexId VertexIdType, message PregelMessage) {
	vertex.Activate()
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

func (vertex *Vertex) InterpretMessages() {
	for _, message := range vertex.ReceivedMessages {
		vertex.InterpretSingleMessage(message)
		vertex.Activate()
	}
}
