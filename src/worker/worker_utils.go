package worker

import "fmt"

func (worker *Worker) getWorkerSubGraphFile() string {
	return "SubGraph-" + fmt.Sprint(worker.id) + ".json"
}
