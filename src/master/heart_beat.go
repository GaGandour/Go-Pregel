package master

import (
	"log"
	"pregel/remote_worker"
	"sync"
	"time"
)

func (master *Master) heartBeatCycle() {
	for {
		time.Sleep(time.Duration(HEARTBEAT_INTERVAL_MS) * time.Millisecond)
		err := master.orderHeartBeats()
		if err != nil {
			log.Println("Error in heartbeat cycle:", err)
			return // Stop the heartbeat cycle
		}
	}
}

func (master *Master) orderHeartBeats() error {
	var wg sync.WaitGroup
	errs := make(chan error, len(master.workers))

	log.Println("Ordering workers to send heartbeat")
	for _, worker := range master.workers {
		wg.Add(1)
		go func(w *remote_worker.RemoteWorker) {
			err := master.orderHeartBeat(w, &wg)
			errs <- err
		}(worker)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			log.Println("Error in heartbeat:", err)
			return err
		}
	}
	return nil
}
