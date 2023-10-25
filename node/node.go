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
	Ip     string
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

var system = color.New(color.FgHiGreen).Add(color.BgBlack)

// TODO
func (node *Node) JoinNetwork(helper string) {
	// system.Println("Joining network...")
	if len(strings.Split(helper, ":")) == 1 { // I am the only node in this network
		system.Println("I am creating a new network...")
		node.Successor = Pointer{Nodeid: node.Nodeid, Ip: node.IP}
		node.Predecessor = Pointer{}
		node.FingerTable = make([]Pointer, 64)
		node.createFingerTable(node.Nodeid)
		system.Println("Finger table has been updated...", node.FingerTable)
	} else { // I am not the only one in this network, and I am joining using someone elses address-> "helper"
		system.Println("Contacting node in network at address", helper)
		reply := node.CallRPC(message.RequestMessage{Type: FIND_SUCCESSOR, TargetId: node.Nodeid}, helper)
		node.Successor = Pointer{Nodeid: reply.Nodeid, Ip: reply.IP}
		node.Predecessor = Pointer{}
		node.FingerTable = make([]Pointer, 64)
		node.createFingerTable(node.Nodeid)
		system.Println("Finger table has been updated...", node.FingerTable)
	}
}

func (node *Node) Sendping(nodeid string) {
	clnt, err := rpc.Dial("tcp", "127.0.0.1:3004")
	if err != nil {
		system.Println("Error Dialing RPC:", err)
	}
	msg := message.RequestMessage{Type: PING}
	reply := message.ResponseMessage{}

	err = clnt.Call("Node.HandleIncomingMessage", msg, &reply)
	if err != nil {
		system.Println("Error Calling RPC:", err)
	}
	if reply.Type == ACK {
		system.Println("Received an ACK")
	}
}

func (node *Node) HandleIncomingMessage(msg *message.RequestMessage, reply *message.ResponseMessage) error {
	system.Println("Message of type", msg.Type, "received.")

	for {
		switch msg.Type {
		case PING:
			// watever
		case ACK:
			// ...
		case FIND_SUCCESSOR:
			node.FindSuccessor(msg.TargetId)
		default:
			// system.Println("Client is alive and listening")
			time.Sleep(1000)
		}
	}
}

func (node *Node) FindSuccessor(id uint64) Pointer {
	if belongsTo(id, node.Nodeid, node.FingerTable[0].Nodeid) {
		return node.FingerTable[0]
	}
	p := node.ClosestPrecedingNode(id)
	reply := node.CallRPC(message.RequestMessage{Type: FIND_SUCCESSOR, TargetId: id}, p.Ip)
	return Pointer{Nodeid: reply.Nodeid, Ip: reply.IP}
}

func (node *Node) ClosestPrecedingNode(id uint64) Pointer {
	for i := 64; i >= 0; i-- {
		if belongsTo(node.FingerTable[i].Nodeid, node.Nodeid, id) {
			return node.FingerTable[i]
		}
	}
	return Pointer{Nodeid: id, Ip: node.IP}
}

func (node *Node) createFingerTable(nodeid uint64) {
	for i := 0; i < 64; i++ {
		node.FingerTable[i] = node.FindSuccessor(node.Nodeid + uint64(math.Pow(2, float64(i)))) // possible floating point error here.
	}
}

// func (node *Node) createFingerTable() {
// 	for i
// }

/*
***************************************
		UTILITY FUNCTIONS
***************************************
*/

// To check if an ID is in a given range (right-inclusive).
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
	clnt, err := rpc.Dial("tcp", IP)
	if err != nil {
		system.Println("Error Dialing RPC:", err)
	}
	reply := message.ResponseMessage{}
	err = clnt.Call("Node.HandleIncomingCommunication", msg, &reply)
	return reply
}
