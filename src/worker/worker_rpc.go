package worker

import (
	"log"
	"pregel/customrpc"
	"pregel/graph_package"
	"pregel/remote_worker"
)

// RPC - RegisterSubGraph
func (worker *Worker) RegisterSubGraph(args *customrpc.RegisterSubGraphArgs, reply *customrpc.RegisterSubGraphReply) error {
	log.Println("Registering SubGraph")
	worker.graph = graph_package.ConvertCommunicationGraphToGraph(&args.SubGraph)
	worker.id = args.WorkerId
	// Build map
	for key, value := range args.RemoteWorkersMap {
		worker.remoteWorkersMap[key] = &remote_worker.RemoteWorker{
			Id:       value.Id,
			Hostname: value.Hostname,
			Status:   value.Status,
		}
	}
	worker.numberOfPartitions = len(worker.remoteWorkersMap)
	return nil
}

// RPC - WriteSubGraphToFile
func (worker *Worker) WriteSubGraphToFile(args *customrpc.WriteSubGraphToFileArgs, reply *customrpc.WriteSubGraphToFileReply) error {
	log.Println("Writing SubGraph to file")
	outputFileName := worker.getWorkerSubGraphFile()
	worker.graph.WriteGraphToFile(outputFileName)
	return nil
}

// RPC - RunSuperStep
func (worker *Worker) RunSuperStep(args *customrpc.RunSuperStepArgs, reply *customrpc.RunSuperStepReply) error {
	log.Println("Running SuperStep")
	var workerVoteToHalt = true
	for _, vertex := range worker.graph.Vertexes {
		vertex.SuperStep()
		vertex.IncreaseSuperStepNumber()
		workerVoteToHalt = workerVoteToHalt && vertex.VotedToHalt
	}
	reply.VoteToHalt = workerVoteToHalt
	worker.PassMessages()
	worker.superStep++
	return nil
}

// RPC - ReceiveMessages
func (worker *Worker) ReceiveMessages(args *customrpc.ReceiveMessagesArgs, reply *customrpc.ReceiveMessagesReply) error {
	for receiverId, messageList := range args.MessageMap {
		vertex := worker.graph.Vertexes[receiverId]
		for _, message := range messageList {
			vertex.ReceiveMessage(args.SuperStep, message)
		}
	}
	return nil
}

// RPC - HeartBeat
func (worker *Worker) HeartBeat(_ *struct{}, _ *struct{}) error {
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
