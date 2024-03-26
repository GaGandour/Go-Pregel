package master

import (
	"log"
	"net"
	"net/rpc"
)

// RunMaster will start a master node on the map reduce operations.
// In the distributed model, a Master should serve multiple workers and distribute
// the operations to be executed in order to complete the task.
//   - task: the Task object that contains the mapreduce operation.
//   - hostname: the tcp/ip address on which it will listen for connections.
func RunMaster(hostname string) {
	var (
		err          error
		master       *Master
		newRpcServer *rpc.Server
		listener     net.Listener
	)

	log.Println("Running Master on", hostname)

	master = newMaster(hostname)

	newRpcServer = rpc.NewServer()
	newRpcServer.Register(master)

	if err != nil {
		log.Panicln("Failed to register RPC server. Error:", err)
	}

	master.rpcServer = newRpcServer

	listener, err = net.Listen("tcp", master.address)

	if err != nil {
		log.Panicln("Failed to start TCP server. Error:", err)
	}

	master.listener = listener

	go master.acceptMultipleConnections()

	// Ler JSON
	// Particionar Grafo
	// Comandar Superstep 0
	// Esperar retorno
	// Comandar Supersteps até todos os workers terminarem
	// Comandar Escrita do Grafo
	// Esperar retorno
	// Juntar os subgrafos
	// Escrever grande grafo
}
