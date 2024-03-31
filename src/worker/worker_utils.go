package worker

import (
	"fmt"
	"pregel/remote_worker"
)

func (worker *Worker) getWorkerSubGraphFile() string {
	return "./output_graphs/SubGraph-" + fmt.Sprint(worker.id) + ".json"
}

func (worker *Worker) getWorkerHostnameByPartitionId(partitionId int) string {
	return worker.remoteWorkersMap[partitionId].Hostname
}

func (worker *Worker) getRemoteWorkerByPartitionId(partitionId int) *remote_worker.RemoteWorker {
	return worker.remoteWorkersMap[partitionId]
}
