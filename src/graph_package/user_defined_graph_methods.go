package graph_package

func (vertex *Vertex) Compute() {
	// The user will implement this function
	if vertex.numSuperSteps == 0 {
		for receivingId := range vertex.GetOutEdges() {
			vertex.PrepareMessageToVertex(receivingId, PregelMessage{OriginVertexId: vertex.Id, Value: vertex.Value.Value})
		}
		return
	}

	vertex.VoteToHalt()
	for _, message := range vertex.ReceivedMessages {
		if message.Value > vertex.Value.Value {
			newValue := VertexValue{
				Value: message.Value,
			}
			vertex.SetValue(newValue)
			vertex.Activate()
		} else {
			vertex.PrepareMessageToVertex(message.OriginVertexId, PregelMessage{OriginVertexId: vertex.Id, Value: vertex.Value.Value})
		}
	}

	if !vertex.VotedToHalt {
		for receivingId := range vertex.GetOutEdges() {
			vertex.PrepareMessageToVertex(receivingId, PregelMessage{OriginVertexId: vertex.Id, Value: vertex.Value.Value})
		}
	}
}

func CombinePregelMessages(messageList []PregelMessage) []PregelMessage {
	// The user can implement this function
	return messageList
}
