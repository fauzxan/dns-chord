package main

import (
	"bufio"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strings"
	"time"

	"github.com/fauzxan/dns-chord/utility"

	"github.com/fauzxan/dns-chord/node"

	"github.com/fatih/color"
	"github.com/joho/godotenv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Color coded logs
var system = color.New(color.FgCyan).Add(color.BgBlack)

/*
Show a list of options to choose from.
*/
func showmenu() {
	system.Println("********************************")
	system.Println("\t\tMENU")
	system.Println("Press 1 to see the fingertable")
	system.Println("Press 2 to see the successor and predecessor")
	system.Println("Press 3 to see the node storage")
	system.Println("Press 4 to see the cache")
	system.Println("Press 5 to query a website")
	system.Println("Press m to see the menu")
	system.Println("********************************")
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// get port from cli arguments (specified by user)
	err := godotenv.Load()
	if err != nil {
		log.Error().Msg("Error getting env variables...")
	}

	var port string
	var helperIp string

	// Read your own port number and also the IP address of the other node, if new network
	myIpAddress := utility.GetOutboundIP().String()
	reader := bufio.NewReader(os.Stdin)
	// read input from user
	system.Println("Enter your port number:")
	port, err = reader.ReadString('\n')
	if err != nil {
		log.Error().Err(err).Msg("Error reading input")
	}
	system.Println("Enter IP address and port used to join network:")
	// read input from user
	helperIp, err = reader.ReadString('\n')
	if err != nil {
		log.Error().Err(err).Msg("Error reading input")
	}

	var addr = myIpAddress + ":" + port

	// Create new Node object for yourself
	me := node.Node{
		Nodeid:        utility.GenerateHash(addr),
		IP:            addr[:len(addr)-1],
		CachedQuery:   make(map[uint64]node.LRUCache, 69),
		HashIPStorage: make(map[uint64]map[uint64][]string, 69),
	}

	log.Info().Str("Address", addr)
	log.Info().Uint64("My id is", me.Nodeid)

	// Bind yourself to a port and listen to it
	tcpAddr, err := net.ResolveTCPAddr("tcp", me.IP)
	if err != nil {
		log.Error().Err(err).Msg("Error resolving TCP address")
	}
	inbound, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Error().Err(err).Msg("Could not listen to TCP address")
	}

	// Register RPC methods and accept incoming requests
	rpc.Register(&me)
	log.Info().Msgf("Node is running at IP address: %s", tcpAddr.String())
	go rpc.Accept(inbound)

	helperIp = helperIp[:len(helperIp)-1]

	/*
		When a node first joins, it checks if it is the first node, then creates a new
		chord network, or joins an existing chord network accordingly.
	*/
	if len(strings.Split(helperIp, ":")) == 1 { // I am the only node in this network
		me.CreateNetwork()
	} else {
		me.JoinNetwork(helperIp)
	}

	showmenu()
	// Keep the parent thread alive
	for {
		time.Sleep(1000)
		var input string
		system.Println("********************************")
		system.Println("     Enter 1, 2, 3, 4, 5, m:    ")
		system.Println("********************************")
		fmt.Scanln(&input)

		switch input {
		case "1":
			system.Println("Printing Fingertable:")
			me.PrintFingers()
		case "2":
			system.Println("Printing Successor and Predecessor:")
			system.Println("Successor:")
			me.PrintSuccessor()
			system.Println("Predecessor:")
			me.PrintPredecessor()
		case "3":
			system.Println("Printing Node Storage:")
			me.PrintStorage()
		case "4":
			system.Println("Printing Cache:")
			me.PrintCache()
		case "5":
			log.Info().Msg("Querying website:")
			system.Println("Please type the website:")
			// Pause logging
			zerolog.SetGlobalLevel(zerolog.Disabled)
			fmt.Scanln(&input)
			// Resume logging
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
			me.QueryDNS(input)
		case "m":
			showmenu()
		default:
			log.Warn().Msg("Invalid input...")
		}
	}
}
