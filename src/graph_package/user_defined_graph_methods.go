package graph_package

func (vertex *Vertex) Compute() {
	// The user will implement this function
	for _, message := range vertex.ReceivedMessages {
		if message.Value > vertex.Value.Value {
			newValue := VertexValue(message)
			vertex.SetValue(newValue)
			vertex.Activate()
		}
	}

	if !vertex.VotedToHalt {
		for receivingId := range vertex.GetOutEdges() {
			vertex.PrepareMessageToVertex(receivingId, PregelMessage{vertex.Value.Value})
		}
	}
	vertex.VoteToHalt()
}

func CombinePregelMessages(messageList []PregelMessage) []PregelMessage {
	// The user can implement this function
	return messageList
}
