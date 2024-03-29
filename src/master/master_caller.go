package master

import (
	"log"
	"pregel/customrpc"
	"pregel/graph_package"
	"pregel/remote_worker"
)

func (master *Master) sendSubGraphToWorker(remoteWorker *remote_worker.RemoteWorker, subGraph *graph_package.Graph) {
	var (
		err  error
		args *customrpc.RegisterSubGraphArgs
	)

	args = &customrpc.RegisterSubGraphArgs{}
	reply := interface{}(nil)

	args.SubGraph = *subGraph

	err = remoteWorker.CallRemoteWorker("Worker.RegisterSubGraph", args, &reply, &master.wg)

	if err != nil {
		log.Printf("Failed to send subgraph to worker. Error: %v\n", err)
	}
}

func (master *Master) orderMessagePassing(remoteWorker *remote_worker.RemoteWorker) {
	var (
		err  error
		args *customrpc.PassMessagesArgs
	)

	args = &customrpc.PassMessagesArgs{}
	reply := interface{}(nil)

	err = remoteWorker.CallRemoteWorker("Worker.PassMessages", args, &reply, &master.wg)

	if err != nil {
		log.Printf("Failed to order message passing to worker. Error: %v\n", err)
	}
}

func (master *Master) orderSuperStep(remoteWorker *remote_worker.RemoteWorker) {
	var (
		err  error
		args *customrpc.RunSuperStepArgs
	)

	args = &customrpc.RunSuperStepArgs{}
	reply := customrpc.RunSuperStepReply{}

	err = remoteWorker.CallRemoteWorker("Worker.RunSuperStep", args, &reply, &master.wg)

	if err != nil {
		log.Printf("Failed to order superstep to worker. Error: %v\n", err)
	}
	master.votesToHaltChan <- reply.VoteToHalt
}

func (master *Master) orderWriteSubGraph(remoteWorker *remote_worker.RemoteWorker) {
	var (
		err  error
		args *customrpc.WriteSubGraphToFileArgs
	)

	args = &customrpc.WriteSubGraphToFileArgs{}
	reply := interface{}(nil)

	err = remoteWorker.CallRemoteWorker("Worker.WriteSubGraphToFile", args, &reply, &master.wg)

	if err != nil {
		log.Printf("Failed to order write subgraph to worker. Error: %v\n", err)
	}
}
