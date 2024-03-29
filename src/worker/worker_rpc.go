package worker

import (
	"log"
	"pregel/customrpc"
	"pregel/graph_package"
)

// RPC - RegisterSubGraph
func (worker *Worker) RegisterSubGraph(args *customrpc.RegisterSubGraphArgs, reply *customrpc.RegisterSubGraphReply) error {
	worker.graph = args.SubGraph
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

// RPC - ReceiveMessage
func (worker *Worker) ReceiveMessage(args *customrpc.ReceiveMessageArgs, reply *customrpc.ReceiveMessageReply) error {
	vertexId := args.VertexId
	vertex := worker.graph.Vertexes[vertexId]
	vertex.MessageMutex.Lock()
	vertex.ReceivedMessages = append(vertex.ReceivedMessages, args.Message)
	vertex.MessageMutex.Unlock()
	return nil
}

// RPC - PassMessages
func (worker *Worker) PassMessages(args *customrpc.PassMessagesArgs, reply *customrpc.PassMessagesReply) error {
	for _, sendingVertex := range worker.graph.Vertexes {
		for receiverId, messageList := range sendingVertex.MessagesToSend {
			combinedMessageList := graph_package.CombinePregelMessages(messageList)
			for _, message := range combinedMessageList {
				partitionToReceiveMessage := graph_package.GetPartitionIdFromVertex(worker.numberOfPartitions, receiverId)
				if partitionToReceiveMessage == worker.id {
					// register message in vertex
					receivingVertex := worker.graph.Vertexes[receiverId]
					receivingVertex.MessageMutex.Lock()
					receivingVertex.ReceivedMessages = append(receivingVertex.ReceivedMessages, message)
					receivingVertex.MessageMutex.Unlock()
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
