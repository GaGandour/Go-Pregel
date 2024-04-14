package graph_package

/*
Vertex methods.
*/

func (vertex *Vertex) ComputeInSuperStepZero() {
	/*
		This method is called in the first superstep. Depending
		on the algorithm, the user might just call the Compute
		method with no messages (as no messages have been passed yet).

		However, if you want a different behavior in the first superstep,
		(most likely initializing fields and such), you can implement this
		method.
	*/
	vertex.Compute([]PregelMessage{})
}

func (vertex *Vertex) Compute(receivedMessages []PregelMessage) {
	/*
		This method is called in every superstep. The user should implement
		the algorithm here, and this method is the most important part of
		Pregel computation. Remember to interpret the receivedMessages and to
		call VoteToHalt() when the vertex is done.
	*/
	vertex.VoteToHalt()
}

/*
Combiner (not really a method, but a function)
*/

func CombinePregelMessages(messageList []PregelMessage) []PregelMessage {
	/*
		Implementing this function is completely optional and the only
		purpose of doing so is to reduce the number of messages that are
		sent between vertices (long story short, to reduce network traffic).

		You DON'T have to call this function anywhere in your code, as this
		is already called in the right place in the framework.

		You have to generate a new list of messages that is smaller than
		the original list, but that still contains the same information.

		For example, if your algorithm wants to find the maximum value
		among all vertexes, and if there are two messages to the same vertex
		with values 5 and 7, you can combine them into a single message
		with value 7.

		However, in smaller algorithms, the effort to implement this function
		does not compensate the benefits of using it.
	*/
	return messageList
}

/*
If any other methods are necessary, write them here
*/
