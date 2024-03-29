package graph

func (vertex *Vertex) InterpretMessages() {
	for _, message := range vertex.receivedMessages {
		// The user will implement this function
		println(message.VertexId)
	}
}

func (vertex *Vertex) Compute() {
	// The user will implement this function
}
