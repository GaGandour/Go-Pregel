package graph_package

import "hash/fnv"

func GetSubGraphInPartition(numberOfPartitions int, graph *Graph, partitionId int) Graph {
	subGraph := Graph{
		totalNumberOfVertexes: graph.totalNumberOfVertexes,
		Vertexes:              make(map[VertexIdType]Vertex),
	}

	for vertexId, vertex := range graph.Vertexes {
		if GetPartitionIdFromVertex(numberOfPartitions, vertex.Id) == partitionId {
			subGraph.Vertexes[vertexId] = vertex
		}
	}

	return subGraph
}

func GetPartitionIdFromVertex(numberOfPartitions int, vertexId VertexIdType) int {
	return customHash(vertexId) % numberOfPartitions
}

func customHash(s VertexIdType) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	hashValue := int(h.Sum32())
	return hashValue
}
