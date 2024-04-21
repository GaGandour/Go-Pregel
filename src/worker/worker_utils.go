package worker

import (
	"pregel/customrpc"
	"pregel/graph_package"
	"pregel/remote_worker"
	"pregel/utils"
)

func (worker *Worker) getWorkerSubGraphFile(isPregelFinished bool) string {
	if isPregelFinished {
		return utils.GetSubGraphOutputFileName(worker.id)
	}
	return utils.GetSuperStepSubGraphOutputFileName(worker.id, worker.superStep)
}

func (worker *Worker) getRemoteWorkerByPartitionId(partitionId int) *remote_worker.RemoteWorker {
	return worker.remoteWorkersMap[partitionId]
}

func (worker *Worker) PassMessages() {
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
					receivingVertex.ReceiveMessage(worker.superStep+1, message)
				}
			}
		} else {
			remoteWorkerToReceive := worker.getRemoteWorkerByPartitionId(partitionId)
			args_remote := &customrpc.ReceiveMessagesArgs{
				SuperStep:  worker.superStep + 1,
				MessageMap: messageMap,
			}
			reply_remote := new(customrpc.ReceiveMessagesReply)
			worker.wg.Add(1)
			remoteWorkerToReceive.CallRemoteWorker("Worker.ReceiveMessages", args_remote, reply_remote, &worker.wg)
		}
	}
	worker.wg.Wait()
}
