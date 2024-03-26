package master

import (
	"log"
	"pregel/customrpc"
	"sync"
)

// runOperation start a single operation on a RemoteWorker and wait for it to return or fail.
func (master *Master) runOperation(remoteWorker *RemoteWorker, operation *Operation, wg *sync.WaitGroup) {
	var (
		err  error
		args *customrpc.RunArgs
	)

	args = &customrpc.RunArgs{}

	reply := interface{}(nil)
	err = remoteWorker.callRemoteWorker(operation.proc, args, &reply)

	if err != nil {
		log.Printf("Operation Failed. Error: %v\n", err)
		wg.Done()
	} else {
		wg.Done()
	}
}
