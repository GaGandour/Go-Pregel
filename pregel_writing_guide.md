# How to use Pregel

In this guide, we'll learn briefly how to convert a sequential algorithm to a distributed pregel version. We'll first study a high-level overview of the conversion process and then dive into the details of the Go-Pregel API.

## High-level Overview

### Understanding the original algorithm

The first step in converting a sequential algorithm to the distributed version is understanding the algorithm itself. Let's take the example of the Single Source Shortest Path (SSSP) algorithm. The most famous graph algorithm to solve this problem is the Dijkstra algorithm. The sequential version of the Dijkstra algorithm is as follows (in python pseudocode):

```python
def dijkstra(graph, source):
    # This is the distance from the source to all other vertices.
    # In the end, this is what we want to return!
    dist = {v: inf for v in graph.vertices}
    dist[source] = 0

    # Priority queue to get the vertex with the smallest distance.
    pq = PriorityQueue()
    pq.push((source, 0))

    # Main loop of the algorithm.
    while not pq.empty():
        u, d = pq.pop()
        for v, w in graph.adj[u]:
            if dist[v] > dist[u] + w:
                dist[v] = dist[u] + w
                pq.push((v, dist[v]))

    # Return the distances.
    return dist
```
To understand the necessary parts of the algorithm, we'll make the following questions:

1. Is there an initialization step? **Answer: Yes, we initialize the distances to infinity and the source to 0.**
2. Are there any inputs to the algorithm, besides the graph itself? **Answer: Yes, we need the source vertex.**
3. What is the algorithm output, and what does it represent? **Answer: The output is the distance from the source to all other vertices.**
4. What happens to the algorithm from the point of view of a single vertex? This question can be divided in the followings:
    1. When does a vertex's data is accessed or changed by the algorithm? **Answer: when we update the distance of a vertex. This happens if and only if the vertex is the first vertex in the priority queue, and it never happens again. A vertex is the first vertex in a priority queue if and only if it is the vertex with the shortest distance registered in it.**
    2. does a vertex stop being accessed by the algorithm? **Answer: when the vertex is the first vertex in the priority queue, that is, after the first time it is explored.**
    3. What happens in each vertex access? What information do we need to do that? **Answer: we change the current distance for each neighbour of the current vertex and store the values to the priority queue.**
    4. is the relationship between the order of vertex accesses and the topology of the graph (i.e., the edges)? **Answer: The order of vertex accesses is determined by the priority queue. The only relationship that exists is that a vertex to be explored is, for sure, a neighbour of an already explored graph.**
    5. the order of exploration matters for the algorithm? If so, is there any specific order change that wouldn't change the result? **Answer: the order of the algorithm does matter and cannot be changed. The order of the algorithm depends heavily on the first element of the priority queue, which is determinedby the topology of the graph and by the weighs of the edges. There isn't any way to change the order of vertex exploration while guaranteeing the same results.**

After question 4.5, if we realize that the order of exploration matters and there is no way to change it without changing the result, the algorithm as it is is not suitable for Pregel (note that this is ***usually*** true for DFS-based algorithms - but not always). We have two choices then: the first is being creative and changing the algorithm in a way that we can change the order of operations without changing the result. The second is to use a different algorithm that is suitable for Pregel. In the end, both options are the same: use another algorithm. In the case of the Dijkstra algorithm, we can use the Bellman-Ford algorithm, which is a more suitable algorithm for Pregel. Let's see the python pseudocode for the Bellman-Ford algorithm and repeat the questions:

```python
def bellman_ford(graph, source):
    # This is the distance from the source to all other vertices.
    # In the end, this is what we want to return!
    dist = {v: inf for v in graph.vertices}
    dist[source] = 0

    # Main loop of the algorithm.
    for _ in range(len(graph.vertices) - 1):
        for u, v, w in graph.edges:
            if dist[v] > dist[u] + w:
                dist[v] = dist[u] + w

    # Return the distances.
    return dist
```
1. Is there an initialization step? **Answer: Yes, we initialize the distances to infinity and the source to 0.**
2. Are there any inputs to the algorithm, besides the graph itself? **Answer: Yes, we need the source vertex.**
3. What is the algorithm output, and what does it represent? **Answer: The output is the distance from the source to all other vertices.**
4. What happens to the algorithm from the point of view of a single vertex? This question can be divided in the followings:
    1. When does a vertex's data is accessed or changed by the algorithm? **Answer: a vertex's data is accessed several times regularly in a loop**
    2. Does a vertex stop being accessed by the algorithm? **Answer: after N-1 iterations of the loop, where N is the number of vertices.**
    3. What happens in each vertex access? What information do we need to do that? **Answer: we check if the distance to a vertex can be updated, based on the current distance of its neighbours and the weight of the respective edges. For that, we need the weight of the edges pointing to the vertex, as well as the distance stored in the other vertex. Or even better, we actually only need the sum of those two numbers.**
    4. Is the relationship between the order of vertex accesses and the topology of the graph (i.e., the edges)? **Answer: None. The vertexes are accessed in order inside a for loop, regardless of topology.**
    5. The order of exploration matters for the algorithm? If so, is there any specific order change that wouldn't change the result? **Answer: as long as, at the end of each iteration, each edge of a vertex has been used once, we can run it in any order and the result won't change.**

Good! It seems that the Bellman-Ford algorithm is suitable for Pregel. Let's see how we can convert it to a distributed version.

### Mapping the initialization step

The pregel algorithm has a initialization step that is executed before the main loop. Here, we recall the answers from questions 1 and 2 from the previous algorithm:

1. Is there an initialization step? **Answer: Yes, we initialize the distances to infinity and the source to 0.**
2. Are there any inputs to the algorithm, besides the graph itself? **Answer: Yes, we need the source vertex.**

The initialization step is directly mapped to the `ComputeInSuperStepZero` method, which we'll see in details later. The other inputs to the algorithm normally are converted to global variables in the pregel algorithm. We'll see a practical example of this later.

### Mapping the main loop

Pregel is based on a big loop in which each vertex executes a previously programmed function defined by YOU, the user. Here, it's important that the order of operations between two different vertices does not matter. In Pregel, each vertex acts as a single independent machine, so in each iteration of the main loop, we need to send messages to the oher vertices with the necessary data for them to update their state (recall the answer from question 4.3).

### Writing the halting condition

Also, we need to make sure the algorithm ends, so each vertex needs to vote to halt the algorithm when it has nothing more to do. We need to check the answer for question 4.2 for that. In the Bellman-Ford algorithm, we can see that a vertex stops being accessed after N-1 iterations of the loop, where N is the number of vertices. So, we can use this as a halting condition.

### Deciding the output format

Finally, we need to decide how we'll output the result. In Pregel, the ouput is ALWAYS a graph, instead of a number, a dictionary or any other thing. But we can store this data in the graph in a clever way. A vertex pregel can store any previously defined data with any schema, so we'll use that to our advantage. In the SSSP problem, we can store the distance from the source to each vertex in the vertex itself. Our output, than, will be the graph itself, but with each updated vertex distance stored in the vertex itself.

## Go-Pregel API

### The User-Defined Graph Types

The files we can change in the Go-Pregel are the ones that start with "user_defined". The first one is the `./src/graph_package/user_defined_graph_types.go` file. This file contains the schemas used to store the data. There are five types you can define here:

1. `VertexIdType`: the type of the vertex id. It can be any hasheable type.
2. `EdgeIdType`: the type of the edge id. It can be any hasheable type.
3. `VertexValue`: this is a struct that will hold the values stored in each vertex. You can put any thing you want here, even other structs. In the SSSP problem, we can store the distance from the source to the vertex here.
4. `EdgeValue`: this is a struct that will hold the values stored in each edge. You can put any thing you want here, even other structs. In the SSSP problem, we can store the weight of the edge here.
5. `PregelMessage`: this is a struct that will hold the messages sent between vertices. You can put any thing you want here, even other structs. In the SSSP problem, we can store the distance from the source to the vertex here together with the edge's weight to the target vertex.

### The Graph Methods

The second file we can change is the `./src/graph_package/user_defined_graph_methods.go` file. This file contains the methods that will be executed in each superstep. The methods are:

1. `ComputeInSuperStepZero`: this method is executed in the first superstep. Here, we can initialize the graph and the vertices.
2. `Compute`: this method is executed in all other supersteps. Here, we can update the vertices and send messages to the other vertices. We can also vote to halt the algorithm here.
3. `CombinePregelMessages`: totally optional, but don't use it if you don't know what you're doing. This method is used to reduce the number of messages sent to a vertex by combining two of them into a single equivalent message.

We can also store the global variables in the top of the file. In the SSSP problem, we can store the source vertex here.

The implementation of the functions 1 and 2 are the most important part of the algorithm, and for that we'll need to use the available methods (not to be implemented nor changed!) in the `./src/graph_package/graph_auxiliary_methods.go` file. The available methods are:

1. `GetSuperStepNumber`: returns the current superstep number.
2. `GetValue`: returns the value stored in the vertex.
3. `SetValue`: sets the value stored in the vertex to a given one.
4. `GetEdgeValue`: returns the value stored in an edge, given the edge ID.
5. `SetEdgeValue`: sets the value stored in an edge, given the edge ID, to a given value.
6. `GetOutEdges`: returns a map of the edges that go out of the vertex. The keys are the edge IDs.
7. `PrepareMessageToVertex`: receives a vertex ID and a message and prepares it to be sent to the vertex in due time. (All messages are sent together, but you don't have o worry about that).
8. `VoteToHalt`: votes to halt the algorithm. If all vertices vote to halt, the algorithm ends.

There are other methods in the file, but those are the only ones that you'll need.

### Other Important Graph Types

Another important type is the `Edge` type. It represents an edge in the graph, and the `Vertex` method `GetOutEdges` returns a map of edges (of type `Edge`) that go out of the vertex. The keys are the edge IDs. The `Edge` type has three fields:

1. `Id`: the edge ID. (type `EdgeIdType`)
2. `To`: the vertex ID of the target vertex. (type `VertexIdType`)
3. `Value`: the value stored in the edge. (type `EdgeValue`)

## The Bellman-Ford Algorithm in Go-Pregel

Now that we know how to convert an algorithm to Pregel, let's see the Bellman-Ford algorithm in Go-Pregel. The first thing we need to do is to define the graph types. We'll use the following types:

```go
type VertexIdType int
type EdgeIdType int

type VertexValue struct {
    Distance int
}

type EdgeValue struct {
    Weight int
}

type PregelMessage struct {
    NewDistance int
}
```

Now, we need to define the global variables and the methods, both in the methods file. We'll use the following global variables:

```go
var source VertexIdType
```

And the following methods:

```go
func (vertex *Vertex) ComputeInSuperStepZero() {
	source = "0"
	if vertex.Id == source {
		vertex.SetValue(VertexValue{Distance: 0})
	} else {
		vertex.SetValue(VertexValue{Distance: -1})
	}

	if vertex.Id == source {
		for _, edge := range vertex.GetOutEdges() {
			pregelMessage := PregelMessage{NewDistance: edge.Value.Weight}
			vertex.PrepareMessageToVertex(edge.To, pregelMessage)
		}
	}
}

func (vertex *Vertex) Compute(receivedMessages []PregelMessage) {
	currentDistance := vertex.GetValue().Distance

	hasChanged := false
	for _, message := range receivedMessages {
		if currentDistance == -1 || message.NewDistance < currentDistance {
			hasChanged = true
			currentDistance = message.NewDistance
		}
	}

	if hasChanged {
		vertex.SetValue(VertexValue{Distance: currentDistance})
		for _, edge := range vertex.GetOutEdges() {
			pregelMessage := PregelMessage{NewDistance: currentDistance + edge.Value.Weight}
			vertex.PrepareMessageToVertex(edge.To, pregelMessage)
		}
	} else {
		vertex.VoteToHalt()
	}
}
```

And that's it! We have the Bellman-Ford algorithm in Go-Pregel. Now we can run it in a distributed way.

After running pregel, we want to customize the output visualization. We can do that on the `./visualization/user_defined_value_displaying.py` file. We can rewrite the functions `vertex_value_to_display` and `edge_value_to_display` to display the information we want. In the SSSP problem, we can display the distance from the source to the vertex in the vertex itself. We also would like to display the weight of the edge in the edge itself. The functions mentioned already come with two parameters: the vertex or edge ID and the value stored in the vertex or edge, as a dictionary. The dictionary follows the same schema as in the Go-Pregel algorithm. The following code is an example of customization:

```python
def vertex_value_to_display(vertex_id, value_dict) -> str:
    return f"""{vertex_id}\nDistance: {value_dict["Distance"]}"""


def edge_value_to_display(edge_id, value_dict) -> str:
    return f"""{edge_id}\nWeight: {value_dict["Weight"]}"""
```

Now, the visualization tool will display the graph with the distances and weights in the vertices and edges, respectively.

![Graph output image](./assets/sssp\ output\ graph.png)
