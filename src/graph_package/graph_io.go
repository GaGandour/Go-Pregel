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
	graph.AutoGenerated = true
	graphJson, err := json.MarshalIndent(graph, "", "\t")
	if err != nil {
		log.Println("Error marshalling graph")
		return err
	}
	os.WriteFile(fileName, graphJson, 0644)
	return nil
}

func ReadGraphFromFile(fileName string) *Graph {
	return ConvertCommunicationGraphToGraph(ReadCommunicationGraphFromFile(fileName))
}

func ReadCommunicationGraphFromFile(fileName string) *CommunicationGraph {
	graph := new(CommunicationGraph)
	graphJson, err := os.ReadFile(fileName)
	if err != nil {
		log.Println("Error reading file")
		return nil
	}
	err = json.Unmarshal(graphJson, &graph)
	if err != nil {
		log.Println("Error unmarshalling graph")
		return nil
	}
	if !graph.AutoGenerated {
		graph.TotalNumberOfVertexes = len(graph.Vertexes)
		for vertexId, communicationVertex := range graph.Vertexes {
			communicationVertex.Id = vertexId
			graph.Vertexes[vertexId] = communicationVertex
			for edgeId, communicationEdge := range communicationVertex.Edges {
				communicationEdge.Id = edgeId
				communicationVertex.Edges[edgeId] = communicationEdge
			}
		}
	}
	return graph
}

func ConvertCommunicationGraphToGraph(communicationGraph *CommunicationGraph) *Graph {
	graph := new(Graph)
	graph.totalNumberOfVertexes = communicationGraph.TotalNumberOfVertexes
	graph.Vertexes = make(map[VertexIdType]*Vertex)
	graph.SuperStep = communicationGraph.NextSuperStep
	for vertexId, communicationVertex := range communicationGraph.Vertexes {
		edges := make(map[EdgeIdType]*Edge)
		for edgeId, communicationEdge := range communicationVertex.Edges {
			edges[edgeId] = &Edge{
				Id:    communicationEdge.Id,
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
			numSuperSteps:               communicationGraph.NextSuperStep,
		}
		if len(communicationVertex.IncomingMessages) > 0 {
			graph.Vertexes[vertexId].ReceivedMessagesInSuperStep[communicationGraph.NextSuperStep] = communicationVertex.IncomingMessages
		}
	}
	return graph
}

func ConvertGraphToCommunicationGraph(graph *Graph) *CommunicationGraph {
	communicationGraph := new(CommunicationGraph)
	communicationGraph.TotalNumberOfVertexes = graph.totalNumberOfVertexes
	communicationGraph.NextSuperStep = graph.SuperStep
	communicationGraph.Vertexes = make(map[VertexIdType]CommunicationVertex)
	for vertexId, vertex := range graph.Vertexes {
		edges := make(map[EdgeIdType]CommunicationEdge)
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
            IncomingMessages: vertex.ReceivedMessagesInSuperStep[graph.SuperStep],
		}
	}
	return communicationGraph
}

func ReduceSubGraphsToCommunicationGraph(fileNames []string) *CommunicationGraph {
	communicationGraph := new(CommunicationGraph)
	communicationGraph.Vertexes = make(map[VertexIdType]CommunicationVertex)
	communicationGraphNextSuperStep := 0 // this is actually impossible and will be overwritten
	for _, fileName := range fileNames {
		graph := ReadCommunicationGraphFromFile(fileName)
		for vertexId, communicationVertex := range graph.Vertexes {
			communicationGraph.Vertexes[vertexId] = communicationVertex
		}
		if communicationGraphNextSuperStep == 0 {
			communicationGraphNextSuperStep = graph.NextSuperStep
		} else {
			if communicationGraphNextSuperStep != graph.NextSuperStep {
				log.Println("Warning: superstep mismatch when reducing subgraphs!")
			}
		}
	}
	communicationGraph.TotalNumberOfVertexes = len(communicationGraph.Vertexes)
	return communicationGraph
}
