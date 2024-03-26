package worker

import (
	"log"
	"net"
	"net/rpc"
	"time"
)

// RunWorker will run a instance of a worker. It'll initialize and then try to register with
// master.
func RunWorker(hostname string, masterHostname string) {
	var (
		err           error
		worker        *Worker
		rpcs          *rpc.Server
		listener      net.Listener
		retryDuration time.Duration
	)

	log.Println("Running Worker on", hostname)

	worker = new(Worker)
	worker.hostname = hostname
	worker.masterHostname = masterHostname
	worker.done = make(chan bool)

	rpcs = rpc.NewServer()
	rpcs.Register(worker)

	worker.rpcServer = rpcs

	listener, err = net.Listen("tcp", worker.hostname)

	if err != nil {
		log.Panic("Starting RPC listener failed. Error:", err)
	}

	worker.listener = listener
	defer worker.listener.Close()

	retryDuration = time.Duration(2) * time.Second
	for {
		err = worker.register()

		if err == nil {
			break
		}

		log.Printf("Registration failed. Retrying in %v seconds...\n", retryDuration)
		time.Sleep(retryDuration)
	}

	go worker.acceptMultipleConnections()

	<-worker.done
}
