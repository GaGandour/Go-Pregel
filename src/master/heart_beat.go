package master

import (
	"log"
	"pregel/remote_worker"
	"sync"
	"time"
)

func (master *Master) heartBeatCycle(failedChannel chan bool) {
	for {
		if master.workerHasFailed {
			return
		}
		time.Sleep(time.Duration(HEARTBEAT_INTERVAL_MS) * time.Millisecond)
		err := master.orderHeartBeats()
		if err != nil {
			log.Println("Error in heartbeat cycle:", err)
			failedChannel <- true
			return // Stop the heartbeat cycle
		}
	}
}

func (master *Master) orderHeartBeats() error {
	var wg sync.WaitGroup
	errs := make(chan error, len(master.workers))
	working_workers := make(chan *remote_worker.RemoteWorker, len(master.workers))
	for _, worker := range master.workers {
		wg.Add(1)
		go func(w *remote_worker.RemoteWorker) {
			err := master.orderHeartBeat(w, nil)
			if err != nil {
				errs <- err
			} else {
				working_workers <- w
			}
			wg.Done()
		}(worker)
	}
	wg.Wait()
	close(errs)
	close(working_workers)
	for err := range errs {
		if err != nil {
			log.Println("Error in heartbeat:", err)
			master.rearrangeWorkers(working_workers)
			return err
		}
	}
	return nil
}

func (master *Master) rearrangeWorkers(workingWorkersChan chan *remote_worker.RemoteWorker) {
	master.workers = make(map[int]*remote_worker.RemoteWorker, 0)
	i := 0
	for worker := range workingWorkersChan {
		worker.Id = i
		master.workers[i] = worker
		i++
	}
	master.numWorkingWorkers = len(master.workers)
	master.totalWorkers = len(master.workers)
}
