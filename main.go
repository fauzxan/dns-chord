package main

import (
	"bufio"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fauzxan/dns-chord/v2/utility"

	"github.com/fauzxan/dns-chord/v2/node"

	"github.com/fatih/color"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Color coded logs
var system = color.New(color.FgCyan).Add(color.BgBlack)

// Set the number of queries to pass in when 5 is pressed.
const numQueries = 100

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
	
	
	/*
	********************
	Logging setup
	********************
	*/
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	
	file, err := os.Create("./logs/" + strconv.Itoa(numQueries) + ".log")
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating log file")
	}
	defer file.Close()
	logWriter := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout}, file)
	log.Logger = zerolog.New(logWriter).With().Timestamp().Logger()
	/*
	********************
	End Logging Setup
	********************
	*/

	// get port from cli arguments (specified by user)
	var port string
	var helperIp string

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

	/*
	Load the data from the CSV file and store it in memory
	*/
	dataList, err := utility.ReadCSV("./website_data/" + "websites" + ".csv")
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
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
			log.Info().Msgf("Querying %v websites", numQueries)
			// Pause logging
			zerolog.SetGlobalLevel(zerolog.Disabled)
			// fmt.Scanln(&input)
			// Resume logging
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
			start := time.Now().UnixMilli()
			for _, query := range dataList[:numQueries]{
				// log.Info().Msg(query)
				go me.QueryDNS(query)
			}
			end := time.Now().UnixMilli()
			timeTaken := end - start
			log.Info().Msgf("TIME %v", timeTaken)
		case "6":
			me.PrintSuccList()
		case "m":
			showmenu()
		default:
			log.Warn().Msg("Invalid input...")
		}
	}
}
