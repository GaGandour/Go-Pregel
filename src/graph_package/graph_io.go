package graph_package

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

func (graph *Graph) WriteGraphToFile(fileName string) error {
	ConvertGraphToCommunicationGraph(graph).WriteGraphToFile(fileName)
	return nil
}

func (graph *CommunicationGraph) WriteGraphToFile(fileName string) error {
	vertexesJson, err := json.MarshalIndent(graph.Vertexes, "", "\t")
	if err != nil {
		log.Println("Error marshalling vertexes")
		return err
	}
	os.WriteFile(fileName, vertexesJson, 0644)
	return nil
}

func ReadGraphFromFile(fileName string) *Graph {
	return ConvertCommunicationGraphToGraph(ReadCommunicationGraphFromFile(fileName))
}

func ReadCommunicationGraphFromFile(fileName string) *CommunicationGraph {
	graph := new(CommunicationGraph)
	vertexesJson, err := os.ReadFile(fileName)
	if err != nil {
		log.Println("Error reading file")
		return nil
	}
	err = json.Unmarshal(vertexesJson, &graph.Vertexes)
	if err != nil {
		log.Println("Error unmarshalling vertexes")
		return nil
	}
	graph.totalNumberOfVertexes = len(graph.Vertexes)
	for vertexId, communicationVertex := range graph.Vertexes {
		communicationVertex.Id = vertexId
		graph.Vertexes[vertexId] = communicationVertex
	}
	return graph
}

func ConvertCommunicationGraphToGraph(communicationGraph *CommunicationGraph) *Graph {
	graph := new(Graph)
	graph.totalNumberOfVertexes = communicationGraph.totalNumberOfVertexes
	graph.Vertexes = make(map[VertexIdType]*Vertex)
	for vertexId, communicationVertex := range communicationGraph.Vertexes {
		edges := make(map[VertexIdType]*Edge)
		for edgeId, communicationEdge := range communicationVertex.Edges {
			edges[edgeId] = &Edge{
				To:    communicationEdge.To,
				Value: communicationEdge.Value,
			}
		}
		graph.Vertexes[vertexId] = &Vertex{
			Id:                          vertexId,
			Value:                       communicationVertex.Value,
			Edges:                       edges,
			ReceivedMessagesInSuperStep: make(map[int][]PregelMessage),
			messageMutex:                sync.Mutex{},
			MessagesToSend:              make(map[VertexIdType][]PregelMessage),
			VotedToHalt:                 false,
			numSuperSteps:               0,
		}
	}
	return graph
}

func ConvertGraphToCommunicationGraph(graph *Graph) *CommunicationGraph {
	communicationGraph := new(CommunicationGraph)
	communicationGraph.totalNumberOfVertexes = graph.totalNumberOfVertexes
	communicationGraph.Vertexes = make(map[VertexIdType]CommunicationVertex)
	for vertexId, vertex := range graph.Vertexes {
		edges := make(map[VertexIdType]CommunicationEdge)
		for edgeId, edge := range vertex.Edges {
			edges[edgeId] = CommunicationEdge{
				To:    edge.To,
				Value: edge.Value,
			}
		}
		communicationGraph.Vertexes[vertexId] = CommunicationVertex{
			Id:    vertex.Id,
			Value: vertex.Value,
			Edges: edges,
		}
	}
	return communicationGraph
}

func ReduceSubGraphsToCommunicationGraph(fileNames []string) *CommunicationGraph {
	communicationGraph := new(CommunicationGraph)
	communicationGraph.Vertexes = make(map[VertexIdType]CommunicationVertex)
	for _, fileName := range fileNames {
		graph := ReadCommunicationGraphFromFile(fileName)
		for vertexId, communicationVertex := range graph.Vertexes {
			communicationGraph.Vertexes[vertexId] = communicationVertex
		}
	}
	communicationGraph.totalNumberOfVertexes = len(communicationGraph.Vertexes)
	return communicationGraph
}
