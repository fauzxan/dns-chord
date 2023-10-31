package main

import (
	"core/node"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

var system = color.New(color.FgCyan).Add(color.BgBlack)

const PING = "ping"

func showmenu() {
	system.Println("********************************")
	system.Println("\t\tMENU")
	system.Println("Press 1 to see the fingertable")
	system.Println("Press m to see the fingertable")
	system.Println("********************************")
}

func main() {
	// get port from cli arguments (specified by user)
	err := godotenv.Load()
	if err != nil {
		system.Println("Error getting env variables...")
	}

	var IPADDRESS = os.Getenv("IPADDRESS")

	var port string
	var joinerPort string
	var ClientNode bool
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
		} else if arg == "-c" {
			system.Println("This is a clientnode!")
			ClientNode = true
		}
	}

	var meNormal *node.Node
	var meClient *node.ClientNode
	// Create new Node object for yourself
	if ClientNode {
		me := node.ClientNode{}
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

		meClient = &me
	} else {
		me := node.Node{}
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

		meNormal = &me
	}

	showmenu()
	// Keep the parent thread alive
	for {
		time.Sleep(1000)
		var input string
		fmt.Scanln(&input)
		if input == "1" && meClient != nil {
			meClient.ShowFingers()
		} else if input == "1" && meNormal != nil {
			meNormal.ShowFingers()
		} else if strings.ToLower(input) == "m" {
			showmenu()
		}
	}
}
