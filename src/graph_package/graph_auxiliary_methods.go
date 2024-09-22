package graph_package

func (vertex *Vertex) GetSuperStepNumber() int {
	/*
	   This method returns the current superstep number.
	*/
	return vertex.numSuperSteps
}

func (vertex *Vertex) GetValue() VertexValue {
	/*
	   This method returns the value of the vertex.
	*/
	return vertex.Value
}

func (vertex *Vertex) SetValue(value VertexValue) {
	/*
	   This method receives a `VertexValue` and sets the value of the vertex to it.
	*/
	vertex.Value = value
}

func (vertex *Vertex) GetEdgeValue(edgeId EdgeIdType) EdgeValue {
	/*
	   This method receives an `EdgeIdType` and returns the value of the edge with that id.
	*/
	edge := vertex.Edges[edgeId]
	return edge.Value
}

func (vertex *Vertex) SetEdgeValue(edgeId EdgeIdType, edgeValue EdgeValue) {
	/*
	   This method receives an `EdgeIdType` and an `EdgeValue` and sets the value of the edge with that id to it.
	*/
	if edge, ok := vertex.Edges[edgeId]; ok {
		edge.Value = edgeValue
		vertex.Activate()
	}
}

func (vertex *Vertex) GetOutEdges() map[EdgeIdType]*Edge {
	/*
	   This method returns a map of all the out edges of the vertex.
	   The key is the edge id and the value is a pointer to the edge itself.
	*/
	return vertex.Edges
}

func (vertex *Vertex) PrepareMessageToVertex(vertexId VertexIdType, message PregelMessage) {
	/*
	   This method receives a `VertexIdType` and a `PregelMessage`.
	   The pregel message will be sent to the vertex with the given ID at the
	   end of the current superstep, together will all other messages that were prepared using
	   this method.
	*/
	vertex.HasSentMessages = true
	vertex.MessagesToSend[vertexId] = append(vertex.MessagesToSend[vertexId], message)
}

func (vertex *Vertex) VoteToHalt() {
	/*
	   This method should be called when the vertex is done with its computation.
	   If all vertices vote to halt, the computation will stop. If a vertex receives a
	   pregel message, it will be reactivated, and it will continue participating in the supersteps
	   even if it has voted to halt before, until this method is called again.
	*/
	vertex.VotedToHalt = true
}
