package worker

import (
	"pregel/remote_worker"
	"pregel/utils"
)

func (worker *Worker) getWorkerSubGraphFile() string {
	return utils.GetSubGraphOutputFileName(worker.id)
}

func (worker *Worker) getRemoteWorkerByPartitionId(partitionId int) *remote_worker.RemoteWorker {
	return worker.remoteWorkersMap[partitionId]
}
