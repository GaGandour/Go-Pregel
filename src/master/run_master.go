package master

import (
	"log"
	"net"
	"net/rpc"
	"pregel/graph_package"
	"pregel/utils"
	"time"
)

const (
	OUTPUT_FILE_NAME = "./output_graphs/output_graph.json"
)

// RunMaster will start a master node on the map reduce operations.
// In the distributed model, a Master should serve multiple workers and distribute
// the operations to be executed in order to complete the task.
//   - task: the Task object that contains the mapreduce operation.
//   - hostname: the tcp/ip address on which it will listen for connections.
func RunMaster(hostname string, inputFile string) {
	var (
		err          error
		master       *Master
		newRpcServer *rpc.Server
		listener     net.Listener
	)

	log.Println("Running Master on", hostname)

	master = newMaster(hostname)

	newRpcServer = rpc.NewServer()
	err = newRpcServer.Register(master)

	if err != nil {
		log.Panicln("Failed to register RPC server. Error:", err)
	}

	master.rpcServer = newRpcServer

	listener, err = net.Listen("tcp", master.address)

	if err != nil {
		log.Panicln("Failed to start TCP server. Error:", err)
	}

	master.listener = listener

	master.getConnectionsFromWorkers()

	// Ler JSON
	log.Println(inputFile)
	graph := graph_package.ReadCommunicationGraphFromFile(inputFile)
	if graph == nil {
		log.Println("Error reading graph")
		return
	}

	// Particionar Grafo
	master.partitionGraph(graph)
	// Comandar Superstep 0
	shouldStopPregel := false
	for !shouldStopPregel {
		// Comandar Supersteps até todos os workers terminarem
		shouldStopPregel = master.orderWorkersToExecuteSuperStep()
	}
	// Comandar Escrita do Grafo
	master.orderWorkersToWriteSubGraphs()
	// Juntar os SubGrafos
	master.reduceSubGraphsAndWriteToFile(OUTPUT_FILE_NAME)
	master.orderFinishOperations()
}

func (master *Master) getConnectionsFromWorkers() {
	log.Println("Waiting for workers to connect")
	go master.acceptMultipleConnections()
	time.Sleep(time.Duration(5) * time.Second)
	master.numWorkingWorkers = len(master.workers)
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
	master.votesToHaltChan = make(chan bool, VOTE_TO_HALT_CHANNEL_BUFFER_SIZE)
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
