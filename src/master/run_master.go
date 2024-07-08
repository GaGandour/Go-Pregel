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
func RunMaster(hostname string, inputFile string, debug bool) {
	var (
		err          error
		master       *Master
		newRpcServer *rpc.Server
		listener     net.Listener
	)

	log.Println("Running Master on", hostname)

	master = newMaster(hostname, debug)

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

	master.executePregel(inputFile)
}

// getConnectionsFromWorkers will wait for workers to connect to the master.
// It will ignore connections that happen after 5 seconds.
func (master *Master) getConnectionsFromWorkers() {
	log.Println("Waiting for workers to connect")
	go master.acceptMultipleConnections()
	time.Sleep(time.Duration(TIME_TO_WAIT_FOR_WORKER_REGISTER_IN_SECONDS) * time.Second)
}

func (master *Master) executePregel(inputFile string) bool {
	pregelStepValues := &PregelStepValues{
		ShouldStopPregel: false,
		InputFile:        inputFile,
		PregelState:      READ_GRAPH_FROM_FILE,
		Graph:            nil,
		Finished:         false,
	}

	for {
		if pregelStepValues.PregelState == END_PREGEL {
			log.Println("Pregel finished")
			return true
		}
		master.executePregelStep(pregelStepValues)
	}
}
