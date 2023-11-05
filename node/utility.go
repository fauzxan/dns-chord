package node

import (
	"core/message"
	"net/rpc"
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
		return reply
	}
	err = clnt.Call("Node.HandleIncomingMessage", msg, &reply)
	if err != nil {
		system.Println("Faced an error trying to call RPC:", err)
		systemcommsin.Println("Received reply", reply)
		return reply
	}
	systemcommsin.Println("Received reply", reply, "from", IP)
	return reply
}

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
	system.Println(node.FingerTable)
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