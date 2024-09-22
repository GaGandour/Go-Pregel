package utils

type WorkerStatus string

const (
	WORKER_STOPPED WorkerStatus = "stopped"
	WORKER_RUNNING WorkerStatus = "running"
	WORKER_WAITING WorkerStatus = "waiting"
	WORKER_FAILED  WorkerStatus = "failed"

	OUTPUT_FILES_DIR    string = "./output_graphs/"
	JSON_FILE_EXTENSION string = ".json"
)
