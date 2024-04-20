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
			ok := master.checkWorker(w)
			if !ok {
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

func (master *Master) partitionGraph(graph *graph_package.CommunicationGraph) {
	log.Println("Partitioning graph")
	// Particionar o grafo
	for partitionId := 0; partitionId < master.numWorkingWorkers; partitionId++ {
		// Enviar subgrafo para worker
		subGraph := graph_package.GetCommunicationSubGraphInPartition(master.numWorkingWorkers, graph, partitionId)
		master.wg.Add(1)
		go master.sendSubGraphToWorker(master.workers[partitionId], &subGraph)
	}
	master.wg.Wait()
}

func (master *Master) orderWorkersToWriteSubGraphs() {
	log.Println("Ordering workers to write subgraphs")
	for _, worker := range master.workers {
		master.wg.Add(1)
		go master.orderWriteSubGraph(worker)
	}
	master.wg.Wait()
}

func (master *Master) orderWorkersToExecuteSuperStep() bool {
	log.Println("Ordering workers to execute superstep")
	master.votesToHaltChan = make(chan bool, MAX_NUM_OF_WORKERS)
	for _, worker := range master.workers {
		master.wg.Add(1)
		go master.orderSuperStep(worker)
	}
	master.wg.Wait()
	close(master.votesToHaltChan)
	for vote := range master.votesToHaltChan {
		if !vote {
			return false
		}
	}
	return true
}

func (master *Master) orderFinishOperations() {
	log.Println("Ordering workers to finish operations")
	for _, worker := range master.workers {
		master.wg.Add(1)
		go master.orderFinishOperation(worker)
	}
	master.wg.Wait()
}

func (master *Master) reduceSubGraphsAndWriteToFile(outputFile string) {
	log.Println("Reducing subgraphs and writing to file")
	fileNames := make([]string, 0)
	for _, worker := range master.workers {
		fileNames = append(fileNames, utils.GetSubGraphOutputFileName(worker.Id))
	}
	// Reduzir os subgrafos
	communicationGraph := graph_package.ReduceSubGraphsToCommunicationGraph(fileNames)
	// Escrever o grafo final
	communicationGraph.WriteGraphToFile(outputFile)
}
