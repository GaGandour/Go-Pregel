package worker

import (
	"log"
	"net"
	"net/rpc"
	"pregel/customrpc"
	"pregel/graph"
)

type Worker struct {
	id int

	// Network
	hostname       string
	masterHostname string
	listener       net.Listener
	rpcServer      *rpc.Server

	done chan bool

	// SubGraph
	graph graph.Graph
}

// Call RPC Register on Master to notify that this worker is ready to receive operations.
func (worker *Worker) register() error {
	var (
		err   error
		args  *customrpc.RegisterArgs
		reply *customrpc.RegisterReply
	)

	log.Println("Registering with Master")

	args = new(customrpc.RegisterArgs)
	args.WorkerHostname = worker.hostname

	reply = new(customrpc.RegisterReply)

	err = worker.callMaster("Master.Register", args, reply)

	if err == nil {
		worker.id = reply.WorkerId
		log.Printf("Registered. WorkerId: %v\n", worker.id)
	}

	return err
}

// acceptMultipleConnections will handle the connections from multiple workers.
func (worker *Worker) acceptMultipleConnections() error {
	var (
		err     error
		newConn net.Conn
	)

	log.Printf("Accepting connections on %v\n", worker.listener.Addr())

	for {
		newConn, err = worker.listener.Accept()

		if err == nil {
			go worker.handleConnection(&newConn)
		} else {
			log.Println("Failed to accept connection. Error: ", err)
			break
		}
	}

	log.Println("Stopped accepting connections.")
	return nil
}

// Handle a single connection until it's done, then closes it.
func (worker *Worker) handleConnection(conn *net.Conn) error {
	worker.rpcServer.ServeConn(*conn)
	(*conn).Close()
	return nil
}

// Connect to Master and call remote procedure.
func (worker *Worker) callMaster(proc string, args interface{}, reply interface{}) error {
	var (
		err    error
		client *rpc.Client
	)

	client, err = rpc.Dial("tcp", worker.masterHostname)
	if err != nil {
		return err
	}

	defer client.Close()

	err = client.Call(proc, args, reply)
	if err != nil {
		return err
	}

	return nil
}
