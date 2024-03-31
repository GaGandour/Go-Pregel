package worker

import (
	"log"
	"pregel/customrpc"
	"pregel/graph_package"
)

// RPC - RegisterSubGraph
func (worker *Worker) RegisterSubGraph(args *customrpc.RegisterSubGraphArgs, reply *customrpc.RegisterSubGraphReply) error {
	log.Println("Registering SubGraph")
	worker.graph = graph_package.ConvertCommunicationGraphToGraph(&args.SubGraph)
	worker.id = args.WorkerId
	// Build map
	for key, value := range args.RemoteWorkersMap {
		worker.remoteWorkersMap[key] = &value
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
		if !vertex.VotedToHalt {
			vertex.SuperStep()
		}
		workerVoteToHalt = workerVoteToHalt && vertex.VotedToHalt
	}
	reply.VoteToHalt = workerVoteToHalt
	return nil
}

// RPC - ReceiveMessage
func (worker *Worker) ReceiveMessage(args *customrpc.ReceiveMessageArgs, reply *customrpc.ReceiveMessageReply) error {
	log.Println("Receiving Message")
	vertexId := args.VertexId
	vertex := worker.graph.Vertexes[vertexId]
	vertex.ReceiveMessage(args.Message)
	return nil
}

// RPC - PassMessages
func (worker *Worker) PassMessages(args *customrpc.PassMessagesArgs, reply *customrpc.PassMessagesReply) error {
	log.Println("Passing Messages")
	for _, sendingVertex := range worker.graph.Vertexes {
		for receiverId, messageList := range sendingVertex.MessagesToSend {
			combinedMessageList := graph_package.CombinePregelMessages(messageList)
			for _, message := range combinedMessageList {
				partitionToReceiveMessage := graph_package.GetPartitionIdFromVertex(worker.numberOfPartitions, receiverId)
				if partitionToReceiveMessage == worker.id {
					// register message in vertex
					receivingVertex := worker.graph.Vertexes[receiverId]
					receivingVertex.ReceiveMessage(message)
				} else {
					remoteWorkerToReceive := worker.getRemoteWorkerByPartitionId(partitionToReceiveMessage)
					args_remote := customrpc.ReceiveMessageArgs{
						Message:  message,
						VertexId: receiverId,
					}
					reply_remote := customrpc.ReceiveMessageReply{}
					worker.wg.Add(1)
					go remoteWorkerToReceive.CallRemoteWorker("Worker.ReceiveMessage", args_remote, reply_remote, &worker.wg)
				}
			}
		}
	}
	worker.wg.Wait()
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
