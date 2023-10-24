package node

import (
	"core/message"
	"fmt"
	"net/rpc"
)

/*
Main logic for:
1. node joining
2. planned departure
3. find_successor
4. closest_preceding_node
*/
type Node struct {
	Nodeid      uint64            // ID of the node in the fingertable
	IP          string            // localhost or IP address AND port number. Can be set through environment variables
	FingerTable map[string]string // id mapping to ip address
	Successor   string            // IP of it's direct successor.
}

// Message types
const PING = "ping"
const ACK = "ack"

// TODO
func (node *Node) JoinNetwork() {
	fmt.Println("Joining network...")
}

func (node *Node) Send(nodeid string, msg message.Message) {
	clnt, err := rpc.Dial("tcp", "127.0.0.1:3004")
	if err != nil {
		fmt.Println("Error Dialing RPC:", err)
	}

	reply := message.Message{}

	err = clnt.Call("Node.HandleIncomingMessage", msg, &reply)
	if err != nil {
		fmt.Println("Error Calling RPC:", err)
	}
	if reply.Type == ACK {
		fmt.Println("Received an ACK")
	}
}

func (node *Node) HandleIncomingMessage(msg *message.Message, reply *message.Message) error {
	fmt.Println("Message of type", msg.Type, "received.")
	if msg.Type == PING {
		reply.Type = ACK
	}
	return nil
}

/*
***************************************
		UTILITY FUNCTIONS
***************************************
*/

/*
Version 1:
fingerTable []Node <-- list of Nodes. ith entry is the closest successor

for i in range 0 to m
[i] --> closest_successor(nodeid + 2^i)

Resulting structure: [Node{Nodeid, IP}]

Version 2:

fingerTable map[uint64][]Node <-- mapping nodeid+2^i to copy of node structure rerpresenting successor of nodeid+2^i

for i in range 0 to m
(nodeid + 2^i) : closest_successor(nodeid + 2^i)

Resulting structure: {nodeid+2^i : Node{Nodeid, IP}}

*/

// 
