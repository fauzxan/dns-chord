package main

import (
	"bufio"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"time"

	"core.com/node"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

// Color coded logs
var system = color.New(color.FgCyan).Add(color.BgBlack)

/*
	Show a list of options to choose from.
*/
func showmenu(){
	system.Println("********************************")
	system.Println("\t\tMENU")
	system.Println("Press 1 to see the fingertable")
	system.Println("Press 2 to see the successor and predecessor")
	system.Println("Press m to see the menu")
	system.Println("********************************")
}

func main() {
	// get port from cli arguments (specified by user)
	err := godotenv.Load()
	if err != nil {
		system.Println("Error getting env variables...")
	}

	// Get the some free port number
	var port, porterr = GetFreePort()
	if porterr != nil{
		system.Println("Error finding new port number")
	}
	var helperIp string

	// Read your own port number and also the IP address of the other node, if new network
	myIpAddress := GetOutboundIP().String() 
	reader := bufio.NewReader(os.Stdin)
	system.Println("Enter IP address and port used to join network:")
	// read input from user
	helperIp, err = reader.ReadString('\n')
	if err != nil {
		system.Fprintln(os.Stderr, "Error reading input:", err)
	}


	// Create new Node object for yourself
	me := node.Node{}
	var addr = myIpAddress + ":" + strconv.Itoa(port)
	system.Println(addr)
	me.IP = addr
	me.Nodeid = GenerateHash(addr)
	system.Println("My id is:", me.Nodeid)

	// Bind yourself to a port and listen to it
	tcpAddr, err := net.ResolveTCPAddr("tcp", me.IP)
	if err != nil {
		system.Println("Error resolving TCP address", err)
	}
	inbound, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		system.Println("Could not listen to TCP address", err)
	}
	system.Println("********************************")
	system.Println("Node is running at:", tcpAddr)
	system.Println("********************************")

	// Register RPC methods and accept incoming requests
	rpc.Register(&me)
	go rpc.Accept(inbound)

	// Join the network using helperIp
	me.JoinNetwork(helperIp[:len(helperIp) - 1])

	showmenu()
	// Keep the parent thread alive
	for {
		time.Sleep(1000)
		var input string
		fmt.Scanln(&input)
		switch input{
		case "1":
			me.PrintFingers()
		case "2":
			system.Println("\n\nSuccessor")
			me.PrintSuccessor()
			system.Println("Predecessor")
			me.PrintPredecessor()
		default:
			system.Println("Enter valid number")
			showmenu()
		}
	}
}