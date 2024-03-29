package graph_package

func (vertex *Vertex) InterpretMessages() {
	for _, message := range vertex.ReceivedMessages {
		// The user will implement this function
		println(message.VertexId)
	}
}

func (vertex *Vertex) Compute() {
	// The user will implement this function
	vertex.VotedToHalt = true
}

func CombinePregelMessages(messageList []PregelMessage) []PregelMessage {
	// The user can implement this function
	return messageList
}
