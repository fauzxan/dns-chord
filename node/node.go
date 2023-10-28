package node

import (
	"core/message"
	"math"
	"net/rpc"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Pointer struct {
	Nodeid uint64
	IP     string
}

type Node struct {
	Nodeid      uint64    // ID of the node
	IP          string    // localhost or IP address AND port number. Can be set through environment variables.
	FingerTable []Pointer // id mapping to ip address
	Successor   Pointer   // Nodeid of it's direct successor.
	Predecessor Pointer   // Nodeid of it's direct predecessor.
}

// Message types
const PING = "ping"
const ACK = "ack"
const FIND_SUCCESSOR = "find_successor"
const CLOSEST_PRECEDING_NODE = "closest_preceding_node"
const GET_PREDECESSOR = "get_predecessor"
const NOTIFY = "notify"

var system = color.New(color.FgHiGreen).Add(color.BgBlack)
var systemcommsin = color.New(color.FgHiMagenta).Add(color.BgBlack)
var systemcommsout = color.New(color.FgHiYellow).Add(color.BgBlack)


func (node *Node) HandleIncomingMessage(msg *message.RequestMessage, reply *message.ResponseMessage) error {
	systemcommsin.Println("Message of type", msg.Type, "received.")
	switch msg.Type {
		case PING:
			// watever
			reply.Type = ACK
		case ACK:
			// ...
		case FIND_SUCCESSOR:
			systemcommsin.Println("Received a message to find successor of", msg.TargetId)
			pointer := node.FindSuccessor(msg.TargetId)
			reply.Type = ACK
			reply.Nodeid = pointer.Nodeid
			reply.IP = pointer.IP
		case NOTIFY:
			systemcommsin.Println("Received a message to notify me about a new predecessor")
			status := node.Notify(Pointer{Nodeid: msg.TargetId, IP: msg.IP})
			if status {
				reply.Type = ACK
			}
		case GET_PREDECESSOR:
			systemcommsin.Println("Received a message to get predecessor")
			reply.Nodeid = node.Predecessor.Nodeid
			reply.IP = node.Predecessor.IP
		default:
			// system.Println("Client is alive and listening")
			time.Sleep(1000)
	}
	return nil
}

// TODO
func (node *Node) JoinNetwork(helper string) {
	if len(strings.Split(helper, ":")) == 1 { // I am the only node in this network
		system.Println("I am creating a new network...")
		node.Successor = Pointer{Nodeid: node.Nodeid, IP: node.IP}
		node.Predecessor = Pointer{}
		node.FingerTable = make([]Pointer, 64)
		node.createFingerTable(node.Nodeid)
		system.Println("Finger table has been updated...", node.FingerTable)
	} else { // I am not the only one in this network, and I am joining using someone elses address-> "helper"
		system.Println("Contacting node in network at address", helper)
		reply := node.CallRPC(message.RequestMessage{Type: FIND_SUCCESSOR, TargetId: node.Nodeid}, helper)
		node.Successor = Pointer{Nodeid: reply.Nodeid, IP: reply.IP}
		system.Println("My successor id is:", node.Successor.Nodeid)
		node.Predecessor = Pointer{}
		node.FingerTable = make([]Pointer, 64)
		node.createFingerTable(node.Nodeid)
		system.Println("Finger table has been updated...", node.FingerTable)
	}
	go node.stabilize()
}


// UNUSED FUNCTION
func (node *Node) Sendping(nodeid string) {
	clnt, err := rpc.Dial("tcp", "127.0.0.1:3004")
	if err != nil {
		system.Println("Error Dialing RPC:", err)
	}
	reply := message.ResponseMessage{}

	err = clnt.Call("Node.HandleIncomingMessage", message.RequestMessage{Type: PING}, &reply)
	if err != nil {
		system.Println("Error Calling RPC:", err)
	}
	if reply.Type == ACK {
		system.Println("Received an ACK")
	}
}

func (node *Node) FindSuccessor(id uint64) Pointer {

	if belongsTo(id, node.Nodeid, node.Successor.Nodeid) {
		// system.Println("Id being found", id) // Debugging output
		// system.Println("A", node.Nodeid)
		// system.Println("B", node.Successor.Nodeid)
		return Pointer{Nodeid: node.Successor.Nodeid, IP: node.Successor.IP} // Case when this is the first node.
	}
	p := node.ClosestPrecedingNode(id)
	system.Println("p nodeid and node nodeid", p.Nodeid, node.Nodeid)
	if p.Nodeid != node.Nodeid {
		reply := node.CallRPC(message.RequestMessage{Type: FIND_SUCCESSOR, TargetId: id}, p.IP)
		return Pointer{Nodeid: reply.Nodeid, IP: reply.IP}
	} else {
		return node.Successor
	}
}

func (node *Node) ClosestPrecedingNode(id uint64) Pointer {
	for i := 63; i >= 0; i-- {
		if belongsTo(node.FingerTable[i].Nodeid, node.Nodeid, id) {
			// system.Println("Closes Preceding node in fingertable:", node.FingerTable[i]) // Debugging output
			return node.FingerTable[i]
		}
	}
	system.Println("Closes Preceding node outside fingertable:", Pointer{Nodeid: node.Nodeid, IP: node.IP})
	return Pointer{Nodeid: node.Nodeid, IP: node.IP}
}

func (node *Node) createFingerTable(nodeid uint64) {
	for i := 0; i < 64; i++ {
		nodePlusTwoI := node.Nodeid + uint64(math.Pow(2, float64(i)))
		// system.Println("Nodeid + 2^i:", nodePlusTwoI) // Debugging output
		successor := node.FindSuccessor(nodePlusTwoI)
		// system.Println("Successor:", successor) // Debugging output
		node.FingerTable[i] = successor // possible floating point error here.
	}
}

func (node *Node) stabilize() {
	for {
		time.Sleep(5 * time.Second)
		if (node.Successor.IP == node.IP) {continue} // Don't need to call your own 
		reply := node.CallRPC(
			message.RequestMessage{Type: GET_PREDECESSOR, TargetId: node.Successor.Nodeid, IP: node.Successor.IP},
			node.Successor.IP,
		)
		sucessorsPredecessor := Pointer{Nodeid: reply.Nodeid, IP: reply.IP}
		if (sucessorsPredecessor != Pointer{}) { // Only execute this block if the successorsPredecessor is not nil
			if belongsTo(sucessorsPredecessor.Nodeid, node.Nodeid, node.Successor.Nodeid) {
				node.Successor = Pointer{Nodeid: sucessorsPredecessor.Nodeid, IP: sucessorsPredecessor.IP}
			}
			reply = node.CallRPC(
				message.RequestMessage{Type: NOTIFY, TargetId: node.Nodeid, IP: node.IP},
				node.Successor.IP,
			)
			if reply.Type == ACK {
				system.Println("Successfully notified successor of it's new predecessor")
			}
		}
	}

}

func (node *Node) Notify(x Pointer) bool {

	if (node.Predecessor == Pointer{} || belongsTo(x.Nodeid, node.Predecessor.Nodeid, node.Nodeid)) {
		node.Predecessor = Pointer{Nodeid: x.Nodeid, IP: x.IP}
	}
	return true
}

/*
***************************************
		UTILITY FUNCTIONS
***************************************
*/

// To check if an ID is in a given range (a, b].
func belongsTo(id, a, b uint64) bool {
	if a == b {
		return true
	}
	if a < b {
		return a < id && id <= b
	}
	return a < id || id <= b
}

func (node *Node) CallRPC(msg message.RequestMessage, IP string) message.ResponseMessage {
	systemcommsout.Println("Sending message", msg)
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
	systemcommsin.Println("Received reply", reply)
	return reply
}
