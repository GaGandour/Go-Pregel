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

func (vertex *Vertex) IncreaseSuperStepNumber() {
	vertex.numSuperSteps++
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

func (vertex *Vertex) Activate() {
	vertex.VotedToHalt = false
}

func (vertex *Vertex) IsHalted() bool {
	return vertex.VotedToHalt
}

func (vertex *Vertex) IsActive() bool {
	return !vertex.VotedToHalt
}
