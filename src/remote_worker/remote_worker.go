package remote_worker

import (
	"net/rpc"
	"pregel/utils"
	"sync"
)

type RemoteWorker struct {
	Id       int
	Hostname string
	Status   utils.WorkerStatus
}

// Call a RemoteWork with the procedure specified in parameters. It will also handle connecting
// to the server and closing it afterwards.
func (worker *RemoteWorker) CallRemoteWorker(proc string, args interface{}, reply interface{}, wg *sync.WaitGroup) error {
	var (
		err    error
		client *rpc.Client
	)
	if wg != nil {
		defer wg.Done()
	}

	client, err = rpc.Dial("tcp", worker.Hostname)

	if err != nil {
		return err
	}

	defer client.Close()
	err = client.Call(proc, args, reply)

	for err != nil {
		var tmpClient *rpc.Client
		tmpClient, err = rpc.Dial("tcp", worker.Hostname)
		if err != nil {
			return err
		}
		defer tmpClient.Close()

		err = tmpClient.Call(proc, args, reply)
	}

	return nil
}
