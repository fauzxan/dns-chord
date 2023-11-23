package node

import (
	"net/rpc"

	"core.com/message"
)

/*
***************************************************
		UTILITY FUNCTIONS FOR NODES
***************************************************
*/

/*
Node utility function to call RPC given a request message, and a destination IP address
*/
func (node *Node) CallRPC(msg message.RequestMessage, IP string) message.ResponseMessage {
	systemcommsout.Println(node.Nodeid, node.IP, "is sending message", msg, "to", IP)
	clnt, err := rpc.Dial("tcp", IP)
	reply := message.ResponseMessage{}
	if err != nil {
		system.Println("Error Dialing RPC:", err)
		systemcommsin.Println("Received reply", reply)
		reply.Type = EMPTY
		return reply
	}
	err = clnt.Call("Node.HandleIncomingMessage", msg, &reply)
	if err != nil {
		system.Println("Faced an error trying to call RPC:", err)
		systemcommsin.Println("Received reply", reply)
		reply.Type = EMPTY
		return reply
	}
	systemcommsin.Println("Received reply", reply, "from", IP)
	return reply
}

// func (node *Node) PutRPC(msg message.RequestMessage, IP string) message.ResponseMessage {
// 	systemcommsout.Println(node.Nodeid, node.IP, "is sending PUT message", "to", IP)
// 	clnt, err := rpc.Dial("tcp", IP)
// 	reply := message.ResponseMessage{}
// 	if err != nil {
// 		system.Println("Error Dialing RPC:", err)
// 		systemcommsin.Println("Received reply", reply)
// 		return reply
// 	}
// 	err = clnt.Call("Node.PutQuery", msg.Payload, &reply)
// 	if err != nil {
// 		system.Println("Faced an error trying to call RPC:", err)
// 		systemcommsin.Println("Received reply", reply)
// 		return reply
// 	}
// 	systemcommsin.Println("Received reply to PUT:", reply, "from", IP)
// 	return reply
// }

/*
Node utility function to check if an ID is in a given range (a, b].
*/
func belongsTo(id, a, b uint64) bool {
	if a == b {
		return true
	}
	if a < b {
		return a < id && id <= b
	} else {
		return a < id || id <= b
	}
}

/*
Node utility function to check if an ID is in a given range (a, b).
*/
func between(id, a, b uint64) bool {
	if a == b {
		return true
	}
	if a < b {
		return a < id && id < b
	} else {
		return a < id || id < b
	}
}

/*
Node utility function to print fingers
*/
func (node *Node) PrintFingers() {
	system.Println("\n\nFINGER TABLE REQUESTED")
	for i := 0; i < len(node.FingerTable); i++ {
		system.Printf("> Finger[%d]: %d : %s\n", i+1, node.FingerTable[i].Nodeid, node.FingerTable[i].IP)
	}
}

/*
Node utility function to print the successor
*/
func (node *Node) PrintSuccessor() {
	system.Println(node.Successor)
}

/*
Node utility function to print predecessor
*/
func (node *Node) PrintPredecessor() {
	system.Println(node.Predecessor)
}

func (node *Node) PrintStorage() {
	system.Println("\n\nSTORAGE TABLE REQUESTED")
	system.Println(node.HashIPStorage)
}

func (node *Node) PrintCache() {
	system.Println("\n\nCACHE TABLE REQUESTED")
	for id, cache := range node.CachedQuery {
		system.Printf(">%d %s\n", id, cache.value)
	}
}
