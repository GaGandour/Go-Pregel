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
	worker.remoteWorkersMap = make(map[int]*remote_worker.RemoteWorker)
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
	outputFileName := worker.getWorkerSubGraphFile(args.IsPregelFinished)
	worker.graph.WriteGraphToFile(outputFileName)
	return nil
}

// RPC - RunSuperStep
func (worker *Worker) RunSuperStep(args *customrpc.RunSuperStepArgs, reply *customrpc.RunSuperStepReply) error {
	log.Println("Running SuperStep")
    // Work failure
    if worker.graph.SuperStep == worker.failureStep {
        log.Println("Worker failed at superstep", worker.graph.SuperStep)
        panic("Worker failed!")
    }
	var workerVoteToHalt = true
	for _, vertex := range worker.graph.Vertexes {
		vertex.SuperStep()
		vertex.IncreaseSuperStepNumber()
		workerVoteToHalt = workerVoteToHalt && vertex.IsHalted() && !vertex.HasSentMessages
	}
	reply.VoteToHalt = workerVoteToHalt
	worker.PassMessages()
	worker.graph.SuperStep++
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
	log.Println("HeartBeat")
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
