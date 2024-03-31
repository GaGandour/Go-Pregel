package master

import (
	"log"
	"net"
	"net/rpc"
	"pregel/graph_package"
	"time"
)

// RunMaster will start a master node on the map reduce operations.
// In the distributed model, a Master should serve multiple workers and distribute
// the operations to be executed in order to complete the task.
//   - task: the Task object that contains the mapreduce operation.
//   - hostname: the tcp/ip address on which it will listen for connections.
func RunMaster(hostname string) {
	var (
		err          error
		master       *Master
		newRpcServer *rpc.Server
		listener     net.Listener
	)

	log.Println("Running Master on", hostname)

	master = newMaster(hostname)

	newRpcServer = rpc.NewServer()
	newRpcServer.Register(master)

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
	graph := graph_package.ReadCommunicationGraphFromFile("../graphs/graph1.json")
	graph.WriteGraphToFile("grafodoidomaster.json")
	if graph == nil {
		log.Println("Error reading graph")
		return
	}

	// Particionar Grafo
	master.partitionGraph(graph)
	// Comandar Superstep 0
	shouldStopPregel := master.orderWorkersToExecuteSuperStep()
	for !shouldStopPregel {
		// Comandar Passagem de Mensagens
		master.orderWorkersToPassMessages()
		// Comandar Supersteps até todos os workers terminarem
		master.orderWorkersToExecuteSuperStep()
	}
	// Comandar Escrita do Grafo
	master.orderWorkersToWriteSubGraphs()
	// Juntar os SubGrafos
}

func (master *Master) getConnectionsFromWorkers() {
	go master.acceptMultipleConnections()
	time.Sleep(time.Duration(5) * time.Second)
	master.numWorkingWorkers = len(master.workers)
}

func (master *Master) partitionGraph(graph *graph_package.CommunicationGraph) {
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
	for _, worker := range master.workers {
		master.wg.Add(1)
		go master.orderWriteSubGraph(worker)
	}
	master.wg.Wait()
}

func (master *Master) orderWorkersToExecuteSuperStep() bool {
	for _, worker := range master.workers {
		master.wg.Add(1)
		go master.orderSuperStep(worker)
	}
	master.wg.Wait()
	for vote := range master.votesToHaltChan {
		if !vote {
			return false
		}
	}
	return true
}

func (master *Master) orderWorkersToPassMessages() {
	for _, worker := range master.workers {
		master.wg.Add(1)
		go master.orderMessagePassing(worker)
	}
	master.wg.Wait()
}
