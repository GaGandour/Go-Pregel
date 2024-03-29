package master

import (
	"log"
	"net"
	"net/rpc"
	"sync"
)

type Master struct {
	// Network
	address   string
	rpcServer *rpc.Server
	listener  net.Listener

	// Workers handling
	workersMutex      sync.Mutex
	workers           map[int]*RemoteWorker
	totalWorkers      int // Used to generate unique ids for new workers
	numWorkingWorkers int
}

type Operation struct {
	proc string
}

// Construct a new Master struct
func newMaster(address string) (master *Master) {
	master = new(Master)
	master.address = address
	master.workers = make(map[int]*RemoteWorker, 0)
	// master.ii = invertedindex.InvertedIndex{}
	master.totalWorkers = 0
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
