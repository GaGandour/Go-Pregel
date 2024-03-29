package graph

import (
	"encoding/json"
	"log"
	"os"
)

func (graph *Graph) WriteGraphToFile(fileName string) error {
	vertexesJson, err := json.MarshalIndent(graph, "", "\t")
	if err != nil {
		log.Println("Error marshalling vertexes")
		return err
	}
	os.WriteFile(fileName, vertexesJson, 0644)
	return nil
}

func (graph *Graph) ReadGraphFromFile(fileName string) error {
	vertexesJson, err := os.ReadFile(fileName)
	if err != nil {
		log.Println("Error reading file")
		return err
	}
	err = json.Unmarshal(vertexesJson, graph)
	if err != nil {
		log.Println("Error unmarshalling vertexes")
		return err
	}
	return nil
}
