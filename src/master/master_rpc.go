package master

import (
	"log"
	"pregel/customrpc"
	"pregel/remote_worker"
	"pregel/utils"
)

// RPC - Register
// Procedure that will be called by workers to register within this master.
func (master *Master) Register(args *customrpc.RegisterArgs, reply *customrpc.RegisterReply) error {
	var (
		newWorker *remote_worker.RemoteWorker
	)
	log.Printf("Registering worker '%v' with hostname '%v'", master.numWorkingWorkers, args.WorkerHostname)

	master.workersMutex.Lock()

	newWorker = &remote_worker.RemoteWorker{
		Id:       master.numWorkingWorkers,
		Hostname: args.WorkerHostname,
		Status:   utils.WORKER_WAITING,
	}
	master.workers[newWorker.Id] = newWorker
	master.numWorkingWorkers++

	master.workersMutex.Unlock()

	*reply = customrpc.RegisterReply{WorkerId: newWorker.Id}
	return nil
}
