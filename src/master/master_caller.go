package master

import (
	"log"
	"pregel/customrpc"
	"pregel/graph_package"
	"pregel/remote_worker"
)

func (master *Master) checkWorker(remoteWorker *remote_worker.RemoteWorker) error {
	var (
		err   error
		args  *customrpc.HeartBeatArgs
		reply *customrpc.HeartBeatReply
	)

	args = new(customrpc.HeartBeatArgs)
	reply = new(customrpc.HeartBeatReply)

	err = remoteWorker.CallRemoteWorker("Worker.HeartBeat", args, reply, &master.wg)

	if err != nil {
		log.Printf("Failed to check worker. Error: %v\n", err)
	}
	return err
}

func (master *Master) sendSubGraphToWorker(remoteWorker *remote_worker.RemoteWorker, subGraph *graph_package.CommunicationGraph) error {
	var (
		err   error
		args  *customrpc.RegisterSubGraphArgs
		reply *customrpc.RegisterSubGraphReply
	)

	args = new(customrpc.RegisterSubGraphArgs)
	reply = new(customrpc.RegisterSubGraphReply)

	args.WorkerId = remoteWorker.Id
	args.NumberOfWorkers = master.numWorkingWorkers
	args.RemoteWorkersMap = make(map[int]remote_worker.RemoteWorker)
	for _, worker := range master.workers {
		args.RemoteWorkersMap[worker.Id] = *worker
	}
	args.SubGraph = *subGraph

	err = remoteWorker.CallRemoteWorker("Worker.RegisterSubGraph", args, reply, &master.wg)

	if err != nil {
		log.Printf("Failed to send subgraph to worker. Error: %v\n", err)
	}
	return err
}

func (master *Master) orderSuperStep(remoteWorker *remote_worker.RemoteWorker) error {
	var (
		err   error
		args  *customrpc.RunSuperStepArgs
		reply *customrpc.RunSuperStepReply
	)

	args = new(customrpc.RunSuperStepArgs)
	reply = new(customrpc.RunSuperStepReply)

	err = remoteWorker.CallRemoteWorker("Worker.RunSuperStep", args, reply, &master.wg)

	if err != nil {
		log.Printf("Failed to order superstep to worker. Error: %v\n", err)
	}
	master.votesToHaltChan <- reply.VoteToHalt
	return err
}

func (master *Master) orderWriteSubGraph(remoteWorker *remote_worker.RemoteWorker) error {
	var (
		err   error
		args  *customrpc.WriteSubGraphToFileArgs
		reply *customrpc.WriteSubGraphToFileReply
	)

	args = new(customrpc.WriteSubGraphToFileArgs)
	reply = new(customrpc.WriteSubGraphToFileReply)

	err = remoteWorker.CallRemoteWorker("Worker.WriteSubGraphToFile", args, reply, &master.wg)

	if err != nil {
		log.Printf("Failed to order write subgraph to worker. Error: %v\n", err)
	}
	return err
}

func (master *Master) orderFinishOperation(remoteWorker *remote_worker.RemoteWorker) error {

	var (
		err   error
		args  *customrpc.DoneArgs
		reply *customrpc.DoneReply
	)

	args = new(customrpc.DoneArgs)
	reply = new(customrpc.DoneReply)
	err = remoteWorker.CallRemoteWorker("Worker.Done", args, reply, &master.wg)

	if err != nil {
		log.Printf("Failed to order finish operation to worker. Error: %v\n", err)
	}
	return err
}
