package graph_package

import (
	"encoding/json"
	"log"
	"os"
)

func (graph *Graph) WriteGraphToFile(fileName string) error {
	vertexesJson, err := json.MarshalIndent(graph.Vertexes, "", "\t")
	if err != nil {
		log.Println("Error marshalling vertexes")
		return err
	}
	os.WriteFile(fileName, vertexesJson, 0644)
	return nil
}

func ReadGraphFromFile(fileName string) *Graph {
	graph := new(Graph)
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
	return graph
}
