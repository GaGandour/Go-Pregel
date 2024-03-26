package worker

import (
	"log"
	"pregel/customrpc"
)

// RPC - RunIntersect
func (worker *Worker) RunIntersect(args *customrpc.RunArgs, reply *interface{}) error {
	return nil
}

// RPC - Done
// Will be called by Master when the task is done.
func (worker *Worker) Done(_ *struct{}, _ *struct{}) error {
	log.Println("Done.")
	defer func() {
		close(worker.done)
	}()
	return nil
}
