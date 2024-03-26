// common_rpc.go defined all the parameters used in RPC between
// master and workers
package customrpc

import (
	"pregel/graph"
)

type RegisterArgs struct {
	WorkerHostname string
}

type RegisterReply struct {
	WorkerId int
}

type RunArgs struct {
	// Id   int
	Vertexes []graph.Vertex
}

type RunSuperStepArgs struct {
	workerId int
}

type RunSuperStepReply struct {
	workerId int
}
