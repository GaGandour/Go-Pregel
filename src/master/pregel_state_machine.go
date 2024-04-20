package master

import (
	"log"
	"pregel/graph_package"
	"pregel/utils"
)

type PregelState int

const (
	READ_GRAPH_FROM_FILE PregelState = iota
	CHECK_WORKERS
	PARTITION_GRAPH
	EXECUTE_SUPERSTEP
	CHECK_HALT_SUPERSTEP
	WRITE_SUBGRAPHS
	REDUCE_SUBGRAPHS
	FINISH_OPERATIONS
	END_PREGEL
)

type PregelStepValues struct {
	ShouldStopPregel bool
	InputFile        string
	PregelState      PregelState
	Graph            *graph_package.CommunicationGraph
}

func (master *Master) executePregelStep(pregelStepValues *PregelStepValues) {
	switch pregelStepValues.PregelState {
	case READ_GRAPH_FROM_FILE:
		// Read JSON file
		inputFile := pregelStepValues.InputFile
		graph := graph_package.ReadCommunicationGraphFromFile(inputFile)
		if graph == nil {
			log.Println("Error reading graph from file")
		}
		pregelStepValues.Graph = graph
		pregelStepValues.PregelState = CHECK_WORKERS
	case CHECK_WORKERS:
		// Check workers
		master.checkWorkers()
		if master.numWorkingWorkers == 0 {
			log.Println("No workers available. Ending Pregel.")
			pregelStepValues.PregelState = END_PREGEL
		}
		pregelStepValues.PregelState = PARTITION_GRAPH
	case PARTITION_GRAPH:
		// Partition graph and distribute to workers
		err := master.partitionGraph(pregelStepValues.Graph)
		if err != nil {
			pregelStepValues.PregelState = CHECK_WORKERS
			return
		}
		pregelStepValues.PregelState = EXECUTE_SUPERSTEP
	case EXECUTE_SUPERSTEP:
		// Tell workers to execute superstep
		shouldStopPregel, err := master.orderWorkersToExecuteSuperStep()
		if err != nil {
			pregelStepValues.PregelState = CHECK_WORKERS
			return
		}
		if shouldStopPregel {
			pregelStepValues.PregelState = WRITE_SUBGRAPHS
		}
	case WRITE_SUBGRAPHS:
		// Tell workers to write subgraphs
		err := master.orderWorkersToWriteSubGraphs()
		if err != nil {
			pregelStepValues.PregelState = CHECK_WORKERS
			return
		}
		pregelStepValues.PregelState = REDUCE_SUBGRAPHS
	case REDUCE_SUBGRAPHS:
		// Reduce subgraphs and write to file
		master.reduceSubGraphsAndWriteToFile(utils.OUTPUT_FILE_NAME)
		pregelStepValues.PregelState = FINISH_OPERATIONS
	case FINISH_OPERATIONS:
		// Tell workers to shut down
		master.orderFinishOperations()
		pregelStepValues.PregelState = END_PREGEL
	}
}
