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
		workerVoteToHalt = workerVoteToHalt && vertex.VotedToHalt
	}
	reply.VoteToHalt = workerVoteToHalt
	return nil
}

// RPC - ReceiveMessages
func (worker *Worker) ReceiveMessages(args *customrpc.ReceiveMessagesArgs, reply *customrpc.ReceiveMessagesReply) error {
	log.Println("Receiving Messages")
	for receiverId, messageList := range args.MessageMap {
		vertex := worker.graph.Vertexes[receiverId]
		for _, message := range messageList {
			vertex.ReceiveMessage(message)
		}
	}
	return nil
}

// RPC - PassMessages
func (worker *Worker) PassMessages(args *customrpc.PassMessagesArgs, reply *customrpc.PassMessagesReply) error {
	log.Println("Passing Messages")
	messagesToSend := make(map[int]map[graph_package.VertexIdType][]graph_package.PregelMessage)

	for _, sendingVertex := range worker.graph.Vertexes {
		for receiverId, messageList := range sendingVertex.MessagesToSend {
			combinedMessageList := graph_package.CombinePregelMessages(messageList)
			partitionToReceiveMessages := graph_package.GetPartitionIdFromVertex(worker.numberOfPartitions, receiverId)
			if messagesToSend[partitionToReceiveMessages] == nil {
				messagesToSend[partitionToReceiveMessages] = make(map[graph_package.VertexIdType][]graph_package.PregelMessage)
			}
			messagesToSend[partitionToReceiveMessages][receiverId] = append(messagesToSend[partitionToReceiveMessages][receiverId], combinedMessageList...)
		}
		sendingVertex.MessagesToSend = make(map[graph_package.VertexIdType][]graph_package.PregelMessage)
	}

	for partitionId, messageMap := range messagesToSend {
		if partitionId == worker.id {
			// register message in vertex
			for receiverVertexId, messageList := range messageMap {
				receivingVertex := worker.graph.Vertexes[receiverVertexId]
				for _, message := range messageList {
					receivingVertex.ReceiveMessage(message)
				}
			}
		} else {
			remoteWorkerToReceive := worker.getRemoteWorkerByPartitionId(partitionId)
			args_remote := &customrpc.ReceiveMessagesArgs{
				MessageMap: messageMap,
			}
			reply_remote := new(customrpc.ReceiveMessagesReply)
			worker.wg.Add(1)
			remoteWorkerToReceive.CallRemoteWorker("Worker.ReceiveMessages", args_remote, reply_remote, &worker.wg)
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
