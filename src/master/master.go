package master

import (
	"log"
	"net"
	"net/rpc"
	"pregel/remote_worker"
	"sync"
)

type Master struct {
	// Network
	address   string
	rpcServer *rpc.Server
	listener  net.Listener

	// Workers handling
	workersMutex      sync.Mutex
	workers           map[int]*remote_worker.RemoteWorker
	numWorkingWorkers int
	wg                sync.WaitGroup

	// Pregel Specific
	votesToHaltChan chan bool

	// Utils
	debug bool
}

type MasterArguments struct {
	Hostname       string
	GraphInputFile string
	Debug          bool
}

// Construct a new Master struct
func newMaster(args MasterArguments) (master *Master) {
	master = new(Master)
	master.address = args.Hostname
	master.workers = make(map[int]*remote_worker.RemoteWorker, 0)
	master.votesToHaltChan = make(chan bool, MAX_NUM_OF_WORKERS)
	master.numWorkingWorkers = 0
	master.debug = args.Debug
	return
}

// acceptMultipleConnections will handle the connections from multiple workers.
func (master *Master) acceptMultipleConnections() {
	var (
		err     error
		newConn net.Conn
	)

	log.Printf("Accepting connections on %v\n", master.listener.Addr())

	for {
		newConn, err = master.listener.Accept()

		if err == nil {
			go master.handleConnection(&newConn)
		} else {
			log.Println("Failed to accept connection. Error: ", err)
			break
		}
	}

	log.Println("Stopped accepting connections.")
}

// Handle a single connection until it's done, then closes it.
func (master *Master) handleConnection(conn *net.Conn) error {
	master.rpcServer.ServeConn(*conn)
	(*conn).Close()
	return nil
}
