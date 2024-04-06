package graph_package

/*
Vertex methods
*/

func (vertex *Vertex) ComputeInSuperStepZero() {
	neighbors := NewVertexIdSet()
	for receivingId := range vertex.GetOutEdges() {
		neighbors.Add(receivingId)
		vertex.PrepareMessageToVertex(receivingId, PregelMessage{OriginVertexId: vertex.Id, Value: vertex.Value.Value})
	}
	vertex.SetValue(VertexValue{Value: vertex.Value.Value, neighbors: *neighbors})
}

func (vertex *Vertex) Compute(receivedMessages []PregelMessage) {
	// The user will implement this function
	currentValue := vertex.Value.Value
	neighbors := vertex.Value.neighbors
	for _, message := range receivedMessages {
		neighbors.Add(message.OriginVertexId)
		if message.Value > currentValue {
			currentValue = message.Value
		}
	}

	if currentValue != vertex.Value.Value || !VertexIdSetsAreEqual(&neighbors, &vertex.Value.neighbors) {
		vertex.SetValue(VertexValue{Value: currentValue, neighbors: neighbors})
		for _, receivingId := range neighbors.ToSlice() {
			vertex.PrepareMessageToVertex(receivingId, PregelMessage{OriginVertexId: vertex.Id, Value: vertex.Value.Value})
		}
	} else {
		vertex.VoteToHalt()
	}
}

func CombinePregelMessages(messageList []PregelMessage) []PregelMessage {
	// The user can implement this function
	return messageList
}
