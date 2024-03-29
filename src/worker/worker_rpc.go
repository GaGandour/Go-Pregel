package worker

import (
	"log"
	"pregel/customrpc"
)

// RPC - RegisterSubGraph
func (worker *Worker) RegisterSubGraph(args *customrpc.RegisterSubGraphArgs, reply *customrpc.RegisterSubGraphReply) error {
	worker.graph = args.SubGraph
	return nil
}

// RPC - WriteSubGraphToFile
func (worker *Worker) WriteSubGraphToFile(args *customrpc.WriteSubGraphToFileArgs, reply *customrpc.WriteSubGraphToFileReply) error {
	outputFileName := worker.getWorkerSubGraphFile()
	worker.graph.WriteGraphToFile(outputFileName)
	return nil
}

// RPC - RunSuperStep
func (worker *Worker) RunSuperStep(args *customrpc.RunSuperStepArgs, reply *customrpc.RunSuperStepReply) error {
	for _, vertex := range worker.graph.Vertexes {
		vertex.SuperStep()
	}
	return nil
}

// RPC - PassMessages
func (worker *Worker) PassMessages(args *customrpc.PassMessagesArgs, reply *customrpc.PassMessagesReply) error {
	// TODO: Implement
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
