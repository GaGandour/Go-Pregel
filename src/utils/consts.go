package utils

type WorkerStatus string

const (
	WORKER_STOPPED WorkerStatus = "stopped"
	WORKER_RUNNING WorkerStatus = "running"
	WORKER_WAITING WorkerStatus = "waiting"
	WORKER_FAILED  WorkerStatus = "failed"
)
