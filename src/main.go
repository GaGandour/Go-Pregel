package main

import (
	"flag"
	"log"
	"pregel/master"
	"pregel/worker"
	"strconv"
)

var (
	// Run mode settings
	nodeType = flag.String("type", "worker", "Node type: master or worker")

	// Network settings
	addr                = flag.String("addr", "localhost", "IP address to listen on")
	port                = flag.Int("port", 5000, "TCP port to listen on")
	masterAddr          = flag.String("master", "localhost:5000", "Master address")
	graphInputFile      = flag.String("graph_file", "../graphs/graph1.json", "Graph input file")
	debug               = flag.Bool("debug", false, "Enable debug mode")
	checkPointFrequency = flag.Int("checkpoint_frequency", 0, "Frequency of checkpointing. 0 means no checkpointing.")
	failureStep         = flag.Int("failure_step", -1, "Step at which worker should fail. Negative value means no failure.")
)

// Code Entry Point
func main() {
	var (
		hostname string
	)

	flag.Parse()

	// Distributed runs the map and reduce operations in remote workers
	// that are registered with a master.
	switch *nodeType {
	case "master":
		log.Println("NodeType:", *nodeType)
		log.Println("Address:", *addr)
		log.Println("Port:", *port)
		log.Println("Graph File:", *graphInputFile)
		log.Println("Debug:", *debug)
		log.Println("Checkpoint Frequency:", *checkPointFrequency)

		hostname = *addr + ":" + strconv.Itoa(*port)

		masterArgs := master.MasterArguments{
			Hostname:       hostname,
			GraphInputFile: *graphInputFile,
			Debug:          *debug,
		}

		// Create fan in and out channels for mapreduce.Tas
		master.RunMaster(masterArgs)

	case "worker":
		log.Println("NodeType:", *nodeType)
		log.Println("Address:", *addr)
		log.Println("Port:", *port)
		log.Println("Master:", *masterAddr)
		log.Println("Failure Step:", *failureStep)

		hostname = *addr + ":" + strconv.Itoa(*port)
		log.Println("Hostname:", hostname)

		workerArgs := worker.WorkerArguments{
			Hostname:       hostname,
			MasterHostname: *masterAddr,
            FailureStep:    *failureStep,
		}
		worker.RunWorker(workerArgs)
	}
}
