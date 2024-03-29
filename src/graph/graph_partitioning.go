package graph

import "hash/fnv"

func getSubGraphInPartition(numberOfPartitions int, graph Graph, partitionId int) Graph {
	var (
		subGraph Graph
	)

	subGraph = Graph{
		Vertexes: make(map[VertexIdType]*Vertex),
	}

	for vertexId, vertex := range graph.Vertexes {
		if getPartitionIdFromVertex(numberOfPartitions, vertex) == partitionId {
			subGraph.Vertexes[vertexId] = vertex
		}
	}

	return subGraph
}

func getPartitionIdFromVertex(numberOfPartitions int, vertex *Vertex) int {
	return customHash(vertex.Id) % numberOfPartitions
}

func customHash(s VertexIdType) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	hashValue := int(h.Sum32())
	return hashValue
}
