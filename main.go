package main

import (
	"core/node"
	"net"
	"net/rpc"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

var system = color.New(color.FgCyan).Add(color.BgBlack)

const PING = "ping"

func main() {
	// get port from cli arguments (specified by user)
	err := godotenv.Load()
	if err != nil {
		system.Println("Error getting env variables...")
	}

	var IPADDRESS = os.Getenv("IPADDRESS")

	var port string
	var joinerPort string
	for i, arg := range os.Args {
		if arg == "-p" {
			if i+1 > len(os.Args) {
				panic("Enter a valid port number for self!!")
			}
			system.Println("Port number specified is", os.Args[i+1])
			port = os.Args[i+1]
		} else if arg == "-u" {
			if i+1 > len(os.Args) {
				panic("Enter a valid port number that you are going to use!!")
			}
			system.Println("Client to join using has port number", os.Args[i+1])
			joinerPort = ":" + os.Args[i+1]
		}
	}

	// Create new Node object for yourself
	var me = node.Node{}
	var addr = IPADDRESS + ":" + port
	me.IP = addr
	me.Nodeid = GenerateHash(addr)
	system.Println("My id is:", me.Nodeid)

	// Bind yourself to a port and listen to it
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		system.Println("Error resolving TCP address", err)
	}
	inbound, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		system.Println("Could not listen to TCP address", err)
	}

	// Register RPC methods and accept incoming requests
	rpc.Register(&me)
	system.Println("Node is runnning at IP address:", tcpAddr)
	go rpc.Accept(inbound)

	me.JoinNetwork(IPADDRESS + joinerPort)

	// Keep the parent thread alive
	for {
		time.Sleep(1000)
	}

}
