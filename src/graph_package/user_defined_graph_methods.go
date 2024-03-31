package graph_package

func (vertex *Vertex) InterpretSingleMessage(message PregelMessage) {
	// The user will implement this function
	if message.Value > vertex.Value.Value {
		newValue := VertexValue{message.Value}
		vertex.SetValue(newValue)
	}
}

func (vertex *Vertex) Compute() {
	// The user will implement this function
	for _, edge := range vertex.GetOutEdges() {
		vertex.PrepareMessageToVertex(edge.To, PregelMessage{vertex.Id, vertex.Value.Value})
	}
}

func CombinePregelMessages(messageList []PregelMessage) []PregelMessage {
	// The user can implement this function
	return messageList
}
