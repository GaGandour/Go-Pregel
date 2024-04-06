package graph_package

func (vertex *Vertex) SuperStep() {
	if vertex.IsHalted() {
		return
	}
	if vertex.GetSuperStepNumber() == 0 {
		vertex.ComputeInSuperStepZero()
	} else {
		vertex.Compute(vertex.ReceivedMessages)
	}
	vertex.numSuperSteps++
	vertex.ReceivedMessages = []PregelMessage{}
}

func (vertex *Vertex) GetSuperStepNumber() int {
	return vertex.numSuperSteps
}

func (vertex *Vertex) GetValue() VertexValue {
	return vertex.Value
}

func (vertex *Vertex) SetValue(value VertexValue) {
	vertex.Value = value
	vertex.Activate()
}

func (vertex *Vertex) GetOutEdges() map[VertexIdType]*Edge {
	return vertex.Edges
}

func (vertex *Vertex) PrepareMessageToVertex(vertexId VertexIdType, message PregelMessage) {
	vertex.messageMutex.Lock()
	defer vertex.messageMutex.Unlock()
	vertex.MessagesToSend[vertexId] = append(vertex.MessagesToSend[vertexId], message)
	vertex.Activate()
}

func (vertex *Vertex) ReceiveMessage(message PregelMessage) {
	vertex.messageMutex.Lock()
	vertex.ReceivedMessages = append(vertex.ReceivedMessages, message)
	vertex.Activate()
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
