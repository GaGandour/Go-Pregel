// common_rpc.go defined all the parameters used in RPC between
// master and workers
package customrpc

import (
	"pregel/graph"
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
	workerId int
}

type RunSuperStepReply struct {
	workerId int
}

type PassMessagesArgs struct {
	workerId int
}

type PassMessagesReply struct {
	workerId int
}

type RegisterSubGraphArgs struct {
	workerId              int
	numberOfWorkers       int
	totalNumberOfVertexes int
	SubGraph              graph.Graph
}

type RegisterSubGraphReply struct {
	workerId int
}

type WriteSubGraphToFileArgs struct {
	workerId int
}

type WriteSubGraphToFileReply struct {
	workerId int
}
