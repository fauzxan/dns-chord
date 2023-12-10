/*
This repository implements a simplified Chord protocol: a decentralized peer-to-peer (P2P) distributed hash table (DHT)
for distributed data storage and lookup. The system orchestrates nodes forming a ring-based network structure where each node
maintains information about its successor, predecessor, and a portion of the network keyspace. It includes functionalities for
node joining, stabilizing the network, and updating finger tables, enabling efficient decentralized lookup of key-value pairs
across a distributed system where each node manages a segment of the overall keyspace. The implementation involves periodic
checks, such as node stabilization, finger table fixing, predecessor checks, and message handling for essential network operations
like finding successors and notifying or updating neighboring nodes.
*/
package node

import (
	"math"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/fauzxan/dns-chord/v2/message"
	"github.com/rs/zerolog/log"
)

// Colour coded logs
var system = color.New(color.FgHiGreen).Add(color.BgBlack)
var systemcommsin = color.New(color.FgHiMagenta).Add(color.BgBlack)
var systemcommsout = color.New(color.FgHiYellow).Add(color.BgBlack)

// Mutex to prevent race condition when accessing SuccList
var mu sync.Mutex

type Pointer struct {
	Nodeid uint64 // ID of the pointed Node
	IP     string // IP of the pointed Node
}

/*
Represents everything that a node in the chord network needs to take care of.
*/
type Node struct {
	Nodeid        uint64                         // ID of the node
	IP            string                         // localhost or IP address AND port number. Can be set through environment variables.
	FingerTable   []Pointer                      // id mapping to ip address
	Successor     Pointer                        // Nodeid of it's direct successor.
	Predecessor   Pointer                        // Nodeid of it's direct predecessor.
	CachedQuery   map[uint64]LRUCache            // caching queries on the node locally
	HashIPStorage map[uint64]map[uint64][]string // storage for hashed ips associated with the node
	CacheTime     uint64                         // To keep track of scalar timestamp to assign to LRUCache
	SuccList      []Pointer                      // Maintain a list of successors for fault tolerance
}

// Constants
const (
	M                  = 32
	CACHE_SIZE         = 5
	REPLICATION_FACTOR = 2
)

// Message types.
const (
	PING                   = "ping"                   // Used to check predecessor.
	ACK                    = "ack"                    // Used for general acknowledgements.
	GET_SUCCESSOR          = "get_successor"          // Used in RPC call to get node.Successor
	FIND_SUCCESSOR         = "find_successor"         // Used to find successor.
	CLOSEST_PRECEDING_NODE = "closest_preceding_node" // Used to find the closest preceding node, given a successor id.
	GET_PREDECESSOR        = "get_predecessor"        // Used to get the predecessor of some node.
	NOTIFY                 = "notify"                 // Used to notify a node about a new predecessor.
	PUT                    = "put"                    // Used to insert a DNS query.
	GET                    = "get"                    // Used to retrieve a DNS record.
	SHIFT                  = "shift"               	  // Used to shift entries.
	EMPTY                  = "empty"                  // Placeholder or undefined message type or errenous communications.
	REPLICATE              = "replicate"              // Used to replicate data.
)

/*
The default method called by all RPCs. This method receives different
types of requests, and calls the appropriate functions.
*/
func (node *Node) HandleIncomingMessage(msg *message.RequestMessage, reply *message.ResponseMessage) error {
	log.Debug().Msgf("Message of type %s received.", msg.Type)
	switch msg.Type {
	case PING:
		log.Debug().Msg("Received PING message")
		reply.Type = ACK
	case GET_SUCCESSOR:
		log.Debug().Msgf("Received a message to GET SUCCESSOR of %d", node.Nodeid)
		reply.Nodeid = node.Successor.Nodeid
		reply.IP = node.Successor.IP
	case FIND_SUCCESSOR:
		log.Debug().Msgf("Received a message to FIND SUCCESSOR of %d", msg.TargetId)
		pointer, _ := node.FindSuccessor(msg.TargetId, msg.HopCount)
		reply.Type = ACK
		reply.Nodeid = pointer.Nodeid
		reply.IP = pointer.IP
	case NOTIFY:
		log.Debug().Msgf("Received a message to NOTIFY me about a new predecessor %d", msg.TargetId)
		status := node.Notify(Pointer{Nodeid: msg.TargetId, IP: msg.IP})
		if status {
			reply.Type = ACK
		}
	case GET_PREDECESSOR:
		log.Debug().Msg("Received a message to GET PREDECESSOR")
		reply.Nodeid = node.Predecessor.Nodeid
		reply.IP = node.Predecessor.IP
	case GET:
		log.Debug().Msg("Received a message to GET DNS record")
		reply.QueryResponse = node.GetQuery(msg.TargetId)
	case SHIFT:
		log.Debug().Msg("Received a message to GET SOME DNS records")
		reply.Payload = node.GetShiftRecords(msg.TargetId)
	case PUT:
		log.Debug().Msg("Received a message to INSERT a query")
		status := node.PutQuery(msg.TargetId, msg.Payload)
		if status {
			reply.Type = ACK
		}
	case REPLICATE:
		log.Debug().Msg("Received a message to REPLICATE data")
		node.processReplicate(msg.TargetId, msg.Payload)
		reply.Type = ACK
	default:
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

// Create new network (genesis node)
func (node *Node) CreateNetwork() {
	log.Info().Msg("> Creating a new network...")
	node.Successor = Pointer{Nodeid: node.Nodeid, IP: node.IP}
	node.Predecessor = Pointer{}
	node.FingerTable = make([]Pointer, M)
	go node.FixFingers()
	log.Info().Msg("> Finger table has been updated...")
	for i := 0; i < len(node.FingerTable); i++ {
		log.Info().Msgf("> Finger[%d]: Nodeid: %d IP: %s", i+1, node.FingerTable[i].Nodeid, node.FingerTable[i].IP)
	}

	// Initialize SuccList with self.
	myPointer := Pointer{node.Nodeid, node.IP}
	node.SuccList = []Pointer{myPointer}

	go node.stabilize()
	go node.CheckPredecessor()
	go node.replicate()
}

// Join existing chord network
func (node *Node) JoinNetwork(helper string) {
	log.Info().Msgf("Contacting node in existing network at address: %s", helper)
	reply := node.CallRPC(message.RequestMessage{Type: FIND_SUCCESSOR, TargetId: node.Nodeid}, helper)
	node.Successor = Pointer{Nodeid: reply.Nodeid, IP: reply.IP}
	log.Info().Msgf("My successor is: Nodeid: %d IP: %s", node.Successor.Nodeid, node.Successor.IP)
	node.Predecessor = Pointer{}
	node.FingerTable = make([]Pointer, M)
	go node.FixFingers()
	log.Info().Msg("> Finger table has been updated...")
	for i := 0; i < len(node.FingerTable); i++ {
		log.Info().Msgf("> Finger[%d]: Nodeid: %d IP: %s", i+1, node.FingerTable[i].Nodeid, node.FingerTable[i].IP)
	}

	log.Info().Msg("Performing key re-distribution")
	reply = node.CallRPC(message.RequestMessage{Type: SHIFT, TargetId: node.Successor.Nodeid}, node.Successor.IP)
	_, ok := node.HashIPStorage[node.Nodeid]
	if !ok {
		node.HashIPStorage[node.Nodeid] = map[uint64][]string{}
	}
	for hashedWebsite := range reply.Payload {
		node.HashIPStorage[node.Nodeid][hashedWebsite] = reply.Payload[hashedWebsite]
	}

	// Initialize SuccList with self.
	myPointer := Pointer{node.Nodeid, node.IP}
	node.SuccList = []Pointer{myPointer}

	go node.stabilize()
	go node.CheckPredecessor()
	go node.replicate()
}

/*
If id falls between its successor, find successor is finished and node
n returns its successor. Otherwise, n searches its finger table for the
node whose ID most immediately precedes id, and then invokes find successor
at that ID
*/
func (node *Node) FindSuccessor(id uint64, hopCount int) (Pointer, int) {
	hopCount++
	if belongsTo(id, node.Nodeid, node.Successor.Nodeid) {
		return Pointer{Nodeid: node.Successor.Nodeid, IP: node.Successor.IP}, hopCount // Case when this is the first node.
	}
	p := node.ClosestPrecedingNode(id)
	if (p != Pointer{} && p.Nodeid != node.Nodeid) {

		reply := node.CallRPC(message.RequestMessage{Type: FIND_SUCCESSOR, TargetId: id, HopCount: hopCount}, p.IP)
		return Pointer{Nodeid: reply.Nodeid, IP: reply.IP}, hopCount
	} else {
		return node.Successor, hopCount
	}
}

/*
Works jointly with FindSuccessor(id). If id doesn't fall between
my id, and my immediate successors id, then we find the closest
preceding node, so we can call find successor on that node.
*/
func (node *Node) ClosestPrecedingNode(id uint64) Pointer {
	for i := M - 1; i >= 0; i-- {
		if belongsTo(node.FingerTable[i].Nodeid, node.Nodeid, id) {
			return node.FingerTable[i]
		}
	}
	log.Info().Msgf("Closest Preceding node outside fingertable: Nodeid: %d IP: %s", node.Nodeid, node.IP)
	return Pointer{Nodeid: node.Nodeid, IP: node.IP}
}

/*
Each node periodically calls fix fingers to make sure its finger
table entries are correct; this is how new nodes initialize
their finger tables, and it is how existing nodes incorporate
new nodes into their finger tables.
*/
func (node *Node) FixFingers() {

	for {
		time.Sleep(1 * time.Second)
		log.Debug().Msg("Fixing fingers...")
		for id := range node.FingerTable {
			nodePlusTwoI := (node.Nodeid + uint64(math.Pow(2, float64(id))))
			power := uint64(math.Pow(2, float64(M)))
			if nodePlusTwoI > power {
				nodePlusTwoI -= power
			}
			node.FingerTable[id], _ = node.FindSuccessor(uint64(nodePlusTwoI), 0)
		}
		// it has just restarted, so it needs to read from storage
		if len(node.HashIPStorage) == 0 {
			go node.readFromStorage()
		}

		go node.writeToStorage()

	}
}

/*
Every node runs stabilize() periodically to learn about newly
joined nodes. Each time node n runs stabilize(), it asks its successor
for the successor’s predecessor p, and decides whether p
should be n’s successor instead. This would be the case if node p
recently joined the system. In addition, stabilize() notifies node
n’s successor of n’s existence, giving the successor the chance
to change its predecessor to n. The successor does this only if it
knows of no closer predecessor than n.
*/
func (node *Node) stabilize() {
	for {
		time.Sleep(1 * time.Second)
		reply := node.CallRPC(
			message.RequestMessage{Type: GET_PREDECESSOR, TargetId: node.Successor.Nodeid, IP: node.Successor.IP},
			node.Successor.IP,
		)

		// [3000, 3001, 3000]

		// Current successor is dead. Look at successor list for next successor.
		if reply.Type == EMPTY {
			// get next successor from SuccList and make it your successor
			for _, pointer := range node.SuccList[1:] {
				if node.checkSuccessorAlive(pointer) {
					node.Successor = pointer
				}
			}

			// Current successor is alive. Check if it's predecessor lies between you and your current successor. If yes, node.Successor = the middle fella
		} else {
			sucessorsPredecessor := Pointer{Nodeid: reply.Nodeid, IP: reply.IP}
			if (sucessorsPredecessor != Pointer{}) {
				// The new dude in between you and your successor is not dead, then my true successor is the new dude. Or you're the only dude.
				if between(sucessorsPredecessor.Nodeid, node.Nodeid, node.Successor.Nodeid) {
					node.Successor = Pointer{Nodeid: sucessorsPredecessor.Nodeid, IP: sucessorsPredecessor.IP}
				}
			}
		}

		// Notify your new successor (whoever it is) that you are it's predecessor
		reply = node.CallRPC(
			message.RequestMessage{Type: NOTIFY, TargetId: node.Nodeid, IP: node.IP},
			node.Successor.IP,
		)
		if reply.Type == ACK {
			log.Debug().Msgf("Successfully notified successor of it's new predecessor Nodeid: %d IP: %s\n", node.Nodeid, node.IP)
		}

		// Recompute SuccList
		node.maintainSuccList()
		// replicate
	}
}

/*
x thinks it might be nodes predecessor
*/
func (node *Node) Notify(x Pointer) bool {
	if (node.Predecessor == Pointer{} || between(x.Nodeid, node.Predecessor.Nodeid, node.Nodeid)) {
		node.Predecessor = Pointer{Nodeid: x.Nodeid, IP: x.IP}
		return true
	}
	return false
}

/*
Each node also runs check predecessor periodically, to clear the node’s
predecessor pointer if n.predecessor has failed; this allows it to accept
a new predecessor in notify.
*/
func (node *Node) CheckPredecessor() {
	for {
		time.Sleep(1 * time.Second)
		if (node.Predecessor == Pointer{}) {
			continue
		}
		reply := node.CallRPC(message.RequestMessage{Type: PING}, node.Predecessor.IP)
		if reply.Type == EMPTY {
			hashMap, ok := node.HashIPStorage[node.Predecessor.Nodeid]
			if ok {
				for id, ip_cache := range hashMap {
					_, ok := node.HashIPStorage[node.Nodeid]
					if !ok {
						node.HashIPStorage[node.Nodeid] = make(map[uint64][]string)
					}
					node.HashIPStorage[node.Nodeid][id] = ip_cache
				}
				delete(node.HashIPStorage, node.Predecessor.Nodeid)
				node.Predecessor = Pointer{}
			}

		} else {
			log.Debug().Msgf("Predecessor Nodeid: %d IP: %s is alive", node.Predecessor.Nodeid, node.Predecessor.IP)
		}
	}
}

func (node *Node) maintainSuccList() {

	mu.Lock()
	myPointer := Pointer{Nodeid: node.Nodeid, IP: node.IP}
	node.SuccList = []Pointer{myPointer}
	for i := 0; i < REPLICATION_FACTOR; i++ {
		lastSucc := node.SuccList[len(node.SuccList)-1]
		reply := node.CallRPC(message.RequestMessage{Type: GET_SUCCESSOR}, lastSucc.IP)
		nextSucc := Pointer{Nodeid: reply.Nodeid, IP: reply.IP}
		node.SuccList = append(node.SuccList, nextSucc)
	}
	mu.Unlock()
}

func (node *Node) checkSuccessorAlive(pointer Pointer) bool {
	reply := node.CallRPC(message.RequestMessage{Type: PING}, pointer.IP)
	return reply.Type == ACK
}
