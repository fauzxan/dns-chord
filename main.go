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

var menu = color.New(color.FgCyan).Add(color.BgBlack)

const PING = "ping"

func main() {
	// get port from cli arguments (specified by user)
	err := godotenv.Load()
	if err != nil {
		menu.Println("Error getting env variables...")
	}
	var port string
	if len(os.Args) > 1 && os.Args[1] == "-p" {
		menu.Println("Port number specified is", os.Args[2])
		port = os.Args[2]
	}

	// Create new Node object for yourself
	var me = node.Node{}
	var addr = os.Getenv("IPADDRESS") + ":" + port
	me.IP = addr
	me.Nodeid = GenerateHash(addr)

	// Bind yourself to a port and listen to it
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		menu.Println("Error resolving TCP address")
	}
	inbound, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		menu.Println("Could not listen to TCP address")
	}

	rpc.Register(&me)
	menu.Println("Node is runnning at IP address:", tcpAddr)
	go rpc.Accept(inbound)
	me.JoinNetwork()

	// Keep the parent thread alive
	for {
		time.Sleep(1000)
	}
}
