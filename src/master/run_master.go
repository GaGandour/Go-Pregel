package master

import (
	"log"
	"net"
	"net/rpc"
	"time"
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

	master.runFaultTolerantPregel(inputFile)
}

func (master *Master) getConnectionsFromWorkers() {
	log.Println("Waiting for workers to connect")
	go master.acceptMultipleConnections()
	time.Sleep(time.Duration(5) * time.Second)
	master.numWorkingWorkers = len(master.workers)
}

func (master *Master) runFaultTolerantPregel(inputFile string) {
	heartBeatFailedChan := make(chan bool, 1)
	success := false
	for !success {
		master.workerHasFailed = false
		log.Println("Num workers:", len(master.workers))
		log.Println("Starting Pregel")
		master.orderHeartBeats()
		log.Println("Received HeartBeats")
		go master.heartBeatCycle(heartBeatFailedChan)
		success = master.executePregel(inputFile, heartBeatFailedChan)
		time.Sleep(5 * time.Second)
	}
}
