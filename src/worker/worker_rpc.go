package worker

import (
	"encoding/json"
	"log"
	"os"
	"pregel/customrpc"
)

// RPC - RegisterSubGraph
func (worker *Worker) RegisterSubGraph(args *customrpc.RegisterSubGraphArgs, reply *customrpc.RegisterSubGraphReply) error {
	worker.graph = args.SubGraph
	return nil
}

// RPC - WriteSubGraphToFile
func (worker *Worker) WriteSubGraphToFile(args *customrpc.WriteSubGraphToFileArgs, reply *customrpc.WriteSubGraphToFileReply) error {
	vertexesJson, error := json.MarshalIndent(worker.graph, "", "\t")
	if error != nil {
		log.Println("Error marshalling vertexes")
		return error
	}
	// Write to file
	outputFileName := worker.getWorkerSubGraphFile()
	os.WriteFile(outputFileName, vertexesJson, 0644)
	return nil
}

// RPC - RunSuperStep
func (worker *Worker) RunSuperStep(args *customrpc.RunSuperStepArgs, reply *customrpc.RunSuperStepReply) error {
	return nil
}

// RPC - Done
// Will be called by Master when the task is done.
func (worker *Worker) Done(_ *struct{}, _ *struct{}) error {
	log.Println("Done.")
	defer func() {
		close(worker.done)
	}()
	return nil
}
