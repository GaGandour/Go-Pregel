package worker

import (
	"log"
	"net"
	"net/rpc"
	"time"
)

// RunWorker will run a instance of a worker. It'll initialize and then try to register with
// master.
func RunWorker(args WorkerArguments) {
	var (
		err           error
		worker        *Worker
		rpcs          *rpc.Server
		listener      net.Listener
		retryDuration time.Duration
	)

	log.Println("Running Worker on", args.Hostname)

	worker = newWorker(args)

	rpcs = rpc.NewServer()
	rpcs.Register(worker)

	worker.rpcServer = rpcs
	listener, err = net.Listen("tcp", worker.hostname)

	if err != nil {
		log.Panic("Starting RPC listener failed. Error:", err)
	}

	worker.listener = listener
	defer worker.listener.Close()

	retryDuration = time.Duration(500) * time.Millisecond
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
