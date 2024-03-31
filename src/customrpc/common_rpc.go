// common_rpc.go defined all the parameters used in RPC between
// master and workers
package customrpc

import (
	"pregel/graph_package"
	"pregel/remote_worker"
)

type HeartBeatArgs struct {
	WorkerId int
}

type HeartBeatReply struct {
	WorkerId int
}

type RegisterArgs struct {
	WorkerHostname string
}

type RegisterReply struct {
	WorkerId int
}

type RunSuperStepArgs struct {
	WorkerId int
}

type RunSuperStepReply struct {
	WorkerId   int
	VoteToHalt bool
}

type PassMessagesArgs struct {
	WorkerId int
}

type PassMessagesReply struct {
	WorkerId int
}

type ReceiveMessageArgs struct {
	Message  graph_package.PregelMessage
	VertexId graph_package.VertexIdType
}

type ReceiveMessageReply struct {
	WorkerId int
}

type RegisterSubGraphArgs struct {
	WorkerId         int
	NumberOfWorkers  int
	RemoteWorkersMap map[int]remote_worker.RemoteWorker
	SubGraph         graph_package.Graph
}

type RegisterSubGraphReply struct {
	WorkerId int
}

type WriteSubGraphToFileArgs struct {
	WorkerId int
}

type WriteSubGraphToFileReply struct {
	WorkerId int
}
