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
	Finished         bool
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
			return
		}
		pregelStepValues.PregelState = PARTITION_GRAPH
	case PARTITION_GRAPH:
		// Partition graph and distribute to workers
        // pregelStepValues.Finished = false
        log.Println("Partitioning graph")
        master.lastCompletedSuperStep = master.lastCheckpointSuperStep
        if master.checkpointFrequency > 0 && master.lastCheckpointSuperStep >= 0 {
            log.Println("Reducing subgraphs from last checkpoint")
            pregelStepValues.Graph = master.reduceSubGraphsFromLastCheckpoint()
        }
		err := master.partitionGraph(pregelStepValues.Graph)
		if err != nil {
			pregelStepValues.PregelState = CHECK_WORKERS
			return
		}
		if master.debug {
			pregelStepValues.PregelState = WRITE_SUBGRAPHS
		} else {
			pregelStepValues.PregelState = EXECUTE_SUPERSTEP
		}
	case EXECUTE_SUPERSTEP:
		// Tell workers to execute superstep
		shouldStopPregel, err := master.orderWorkersToExecuteSuperStep()
		if err != nil {
			pregelStepValues.PregelState = CHECK_WORKERS
			return
		}
		master.lastCompletedSuperStep++
		if shouldStopPregel {
			pregelStepValues.Finished = true
			pregelStepValues.PregelState = WRITE_SUBGRAPHS
		}
		if master.debug {
			pregelStepValues.PregelState = WRITE_SUBGRAPHS
		}
		if master.lastCompletedSuperStep > 0 && master.checkpointFrequency > 0 {
			if master.lastCompletedSuperStep%master.checkpointFrequency == 0 {
				pregelStepValues.PregelState = WRITE_SUBGRAPHS
			}
		}
	case WRITE_SUBGRAPHS:
		// Tell workers to write subgraphs
		err := master.orderWorkersToWriteSubGraphs(pregelStepValues.Finished)
		if err != nil {
			pregelStepValues.PregelState = CHECK_WORKERS
			return
		}
		master.lastCheckpointSuperStep = master.lastCompletedSuperStep
		if pregelStepValues.Finished {
			pregelStepValues.PregelState = REDUCE_SUBGRAPHS
		} else {
			pregelStepValues.PregelState = EXECUTE_SUPERSTEP
		}
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
