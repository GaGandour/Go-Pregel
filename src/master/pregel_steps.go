package master

import (
	"log"
	"pregel/graph_package"
	"pregel/remote_worker"
	"pregel/utils"
)

func (master *Master) checkWorkers() {
	log.Println("Checking workers")
	brokenWorkers := make(chan int, MAX_NUM_OF_WORKERS)

	for _, worker := range master.workers {
		master.wg.Add(1)
		go func(w *remote_worker.RemoteWorker) {
			err := master.checkWorker(w)
			if err != nil {
				brokenWorkers <- w.Id
			}
		}(worker)
	}
	master.wg.Wait()
	close(brokenWorkers)
	for workerId := range brokenWorkers {
		delete(master.workers, workerId)
		master.numWorkingWorkers--
	}

	newWorkersMap := make(map[int]*remote_worker.RemoteWorker, 0)
	i := 0
	for _, worker := range master.workers {
		newWorkersMap[i] = worker
		newWorkersMap[i].Id = i
		i++
	}
	master.workers = newWorkersMap
}

func (master *Master) partitionGraph(graph *graph_package.CommunicationGraph) error {
	log.Println("Partitioning graph between", master.numWorkingWorkers, "workers")

	var err error
	for _, worker := range master.workers {
		subGraph := graph_package.GetCommunicationSubGraphInPartition(master.numWorkingWorkers, graph, worker.Id)
		master.wg.Add(1)
		go func(w *remote_worker.RemoteWorker) {
			workerError := master.sendSubGraphToWorker(w, &subGraph)
			if workerError != nil {
				err = workerError
			}
		}(worker)
	}
	master.wg.Wait()
	return err
}

func (master *Master) orderWorkersToWriteSubGraphs(isPregelFinished bool) error {
	log.Println("Ordering workers to write subgraphs")

	var err error
	for _, worker := range master.workers {
		master.wg.Add(1)
		go func(w *remote_worker.RemoteWorker) {
			workerError := master.orderWriteSubGraph(w, isPregelFinished)
			if workerError != nil {
				err = workerError
			}
		}(worker)
	}
	master.wg.Wait()
	return err
}

func (master *Master) orderWorkersToExecuteSuperStep() (bool, error) {
	log.Println("Ordering workers to execute superstep")

	var err error
	master.votesToHaltChan = make(chan bool, MAX_NUM_OF_WORKERS)
	for _, worker := range master.workers {
		master.wg.Add(1)
		go func(w *remote_worker.RemoteWorker) {
			workerError := master.orderSuperStep(w)
			if workerError != nil {
				err = workerError
			}
		}(worker)
	}
	master.wg.Wait()
	close(master.votesToHaltChan)
	if err != nil {
		return false, err
	}
	for vote := range master.votesToHaltChan {
		if !vote {
			return false, nil
		}
	}
	return true, nil
}

func (master *Master) orderFinishOperations() error {
	log.Println("Ordering workers to finish operations")

	var err error
	for _, worker := range master.workers {
		master.wg.Add(1)
		go func(w *remote_worker.RemoteWorker) {
			workerError := master.orderFinishOperation(w)
			if workerError != nil {
				err = workerError
			}
		}(worker)
	}
	master.wg.Wait()
	return err
}

func (master *Master) reduceSubGraphsAndWriteToFile(output_file_path string) {
	log.Println("Reducing subgraphs and writing to file")
	fileNames := make([]string, 0)
	for _, worker := range master.workers {
		fileNames = append(fileNames, utils.GetSubGraphOutputFileName(worker.Id))
	}
	// assert that output_file_path begins with '/graphs/'
	if output_file_path[:8] != "/graphs/" {
		panic("Output file path must begin with '/graphs/'")
	}
	// remove '/graphs/' from the beginning of the output_file_path
	output_file_path = output_file_path[8:]

	communicationGraph := graph_package.ReduceSubGraphsToCommunicationGraph(fileNames)
	communicationGraph.WriteGraphToFile(utils.OUTPUT_FILES_DIR + output_file_path)
}

func (master *Master) reduceSubGraphsFromLastCheckpoint() *graph_package.CommunicationGraph {
	log.Println("Reducing subgraphs from last checkpoint")
	fileNames, err := utils.GetCheckpointFileNamesForSuperstep(master.lastCheckpointSuperStep + 1)
	if err != nil {
		log.Println("Error getting checkpoint file names")
		return nil
	}
	communicationGraph := graph_package.ReduceSubGraphsToCommunicationGraph(fileNames)
	return communicationGraph
}
