package graph_package

func (vertex *Vertex) SuperStep() {
	vertex.HasSentMessages = false
	if vertex.ReceivedMessagesInSuperStep[vertex.GetSuperStepNumber()] != nil {
		vertex.Activate()
	}
	if vertex.IsHalted() {
		return
	}
	if vertex.GetSuperStepNumber() == 0 {
		vertex.ComputeInSuperStepZero()
	} else {
		vertex.Compute(vertex.ReceivedMessagesInSuperStep[vertex.GetSuperStepNumber()])
	}
	delete(vertex.ReceivedMessagesInSuperStep, vertex.GetSuperStepNumber())
}

func (vertex *Vertex) GetSuperStepNumber() int {
	return vertex.numSuperSteps
}

func (vertex *Vertex) IncreaseSuperStepNumber() {
	vertex.numSuperSteps++
}

func (vertex *Vertex) GetValue() VertexValue {
	return vertex.Value
}

func (vertex *Vertex) SetValue(value VertexValue) {
	vertex.Value = value
}

func (vertex *Vertex) GetEdgeValue(edgeId EdgeIdType) EdgeValue {
	edge := vertex.Edges[edgeId]
	return edge.Value
}

func (vertex *Vertex) SetEdgeValue(edgeId EdgeIdType, edgeValue EdgeValue) {
	if edge, ok := vertex.Edges[edgeId]; ok {
		edge.Value = edgeValue
		vertex.Activate()
	}
}

func (vertex *Vertex) GetOutEdges() map[EdgeIdType]*Edge {
	return vertex.Edges
}

func (vertex *Vertex) PrepareMessageToVertex(vertexId VertexIdType, message PregelMessage) {
	vertex.HasSentMessages = true
	vertex.MessagesToSend[vertexId] = append(vertex.MessagesToSend[vertexId], message)
}

func (vertex *Vertex) ReceiveMessage(superStepToReceive int, message PregelMessage) {
	vertex.messageMutex.Lock()
	receivedMessages := vertex.ReceivedMessagesInSuperStep[superStepToReceive]
	if receivedMessages == nil {
		receivedMessages = make([]PregelMessage, 0)
	}
	receivedMessages = append(receivedMessages, message)
	vertex.ReceivedMessagesInSuperStep[superStepToReceive] = receivedMessages
	vertex.messageMutex.Unlock()
}

func (vertex *Vertex) VoteToHalt() {
	vertex.VotedToHalt = true
}

func (vertex *Vertex) Activate() {
	vertex.VotedToHalt = false
}

func (vertex *Vertex) IsHalted() bool {
	return vertex.VotedToHalt
}

func (vertex *Vertex) IsActive() bool {
	return !vertex.VotedToHalt
}
