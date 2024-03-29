package master

import (
	"net/rpc"
	"pregel/utils"
)

type RemoteWorker struct {
	id       int
	hostname string
	status   utils.WorkerStatus
}

// Call a RemoteWork with the procedure specified in parameters. It will also handle connecting
// to the server and closing it afterwards.
func (worker *RemoteWorker) callRemoteWorker(proc string, args interface{}, reply *interface{}) error {
	var (
		err    error
		client *rpc.Client
	)

	client, err = rpc.Dial("tcp", worker.hostname)

	if err != nil {
		return err
	}

	defer client.Close()
	err = client.Call(proc, args, reply)

	for err != nil {
		var tmpClient *rpc.Client
		tmpClient, err = rpc.Dial("tcp", worker.hostname)
		if err != nil {
			return err
		}
		defer tmpClient.Close()

		err = tmpClient.Call(proc, args, reply)

	}

	return nil
}
