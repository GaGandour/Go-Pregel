package master

import (
	"log"
	"pregel/customrpc"
	"pregel/utils"
)

// RPC - Register
// Procedure that will be called by workers to register within this master.
func (master *Master) Register(args *customrpc.RegisterArgs, reply *customrpc.RegisterReply) error {
	var (
		newWorker *RemoteWorker
	)
	log.Printf("Registering worker '%v' with hostname '%v'", master.totalWorkers, args.WorkerHostname)

	master.workersMutex.Lock()

	newWorker = &RemoteWorker{master.totalWorkers, args.WorkerHostname, utils.WORKER_WAITING}
	master.workers[newWorker.id] = newWorker
	master.totalWorkers++

	master.workersMutex.Unlock()

	*reply = customrpc.RegisterReply{WorkerId: newWorker.id}
	return nil
}
