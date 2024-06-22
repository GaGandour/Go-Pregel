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

type DoneArgs struct {
	WorkerId int
}

type DoneReply struct {
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

type ReceiveMessagesArgs struct {
	SuperStep  int
	MessageMap map[graph_package.VertexIdType][]graph_package.PregelMessage
}

type ReceiveMessagesReply struct {
	WorkerId int
}

type RegisterSubGraphArgs struct {
	WorkerId         int
	NumberOfWorkers  int
	RemoteWorkersMap map[int]remote_worker.RemoteWorker
	SubGraph         graph_package.CommunicationGraph
}

type RegisterSubGraphReply struct {
	WorkerId int
}

type WriteSubGraphToFileArgs struct {
	WorkerId         int
	IsPregelFinished bool
}

type WriteSubGraphToFileReply struct {
	WorkerId int
}
