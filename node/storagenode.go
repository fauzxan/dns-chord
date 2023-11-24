/*
The QueryDNS function in the Chord-based DNS system first checks its cache, then local storage, 
and queries the Chord network for website records. If not found, it resorts to traditional DNS. 
It optimizes by caching recent queries and replicating data for resilience. The code handles queries 
efficiently in a distributed network, enhancing DNS performance.
*/

package node

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"core.com/message"
	"core.com/utility"
)

/*
Used for in-memory-storage. Used to maintain the list of recent queries, and improve query speed.
*/
type LRUCache struct {
	value   []string // List of values corresponding to websites records.
	cacheTime uint64   // Counter to indicate the timestamp of the entry. Used for kicking out Least Recently Used.
}

/*
Mother code for all DNS query logic. It executes one of the following logic pathways:

1. Query node -> check local cache -> return entry

2. Query node -> check local cache -> query local storage -> put in local cache -> return entry

3. Query node -> check local cache -> query local storage -> find successor, and send get -> put in local cache -> return entry

4. Query node -> check local cache -> query local storage -> find successor, and send get -> query legacy DNS -> send to appropriate node, or self -> put in local cache -> return entry
*/
func (node *Node) QueryDNS(website string) {
	if node.CachedQuery == nil {
		node.CachedQuery = make(map[uint64]LRUCache)
	}
	node.CacheTime += 1

	if strings.HasPrefix(website, "www.") {
		system.Println("> Removing Prefix")
		website = website[4:]
	}
	hashedWebsite := utility.GenerateHash(website)
	ip_addr, ok := node.CachedQuery[hashedWebsite]
	if ok {
		system.Println("> Retrieving from LRUCache")
		for _, ip_c := range ip_addr.value {
			system.Printf("> %s. IN A %s\n", website, ip_c)
		}
	} else {
		ip_addr, ok := node.HashIPStorage[node.Nodeid][hashedWebsite]
		system.Printf("> The Website %s has been hashed to %d\n", website, hashedWebsite)
		if ok {
			system.Println("> Retrieving from Local Storage")
			for _, ip_c := range ip_addr {
				system.Printf("> %s. IN A %s\n", website, ip_c)
			}
		} else {
			succPointer, hopCount := node.FindSuccessor(hashedWebsite, 0)
			system.Println("> Number of Hops:", hopCount)
			// log hopcount into the log file using the library
			system.Printf(">  The Website would be stored at it's succesor %d : %s\n", succPointer.Nodeid, succPointer.IP)
			msg := message.RequestMessage{Type: GET, TargetId: hashedWebsite}
			reply := node.CallRPC(msg, succPointer.IP)
			if reply.QueryResponse != nil {
				system.Println("> Retrieving from Chord Network")
				for _, ip_c := range reply.QueryResponse {
					system.Printf("> %s. IN A %s\n", website, ip_c)
				}
			} else {
				ips, err := net.LookupIP(website)
				if err != nil {
					system.Printf("> Could not get IPs: %v\n", err)
					return
				}
				ip_addresses := []string{}
				system.Println("IP ADDRESSES", ip_addresses)

				for _, ip := range ips {
					ip_addresses = append(ip_addresses, ip.String())
					system.Printf("> %s. IN A %s\n", website, ip.String())
				}
				node.CachedQuery[hashedWebsite] = LRUCache{value: ip_addresses, cacheTime: node.CacheTime}
				reply = node.CallRPC(message.RequestMessage{Type: PUT, TargetId: succPointer.Nodeid, Payload: map[uint64][]string{hashedWebsite: ip_addresses}}, succPointer.IP)

				if reply.Type == ACK {
					// finding the oldest one based on counter, and removing that key
					if len(node.CachedQuery) > CACHE_SIZE {
						var minKey uint64
						minValue := uint64(18446744073709551615)
						for key, value := range node.CachedQuery {
							if value.cacheTime < minValue {
								minKey = key
								minValue = value.cacheTime
							}
						}
						if minKey != 0 {
							delete(node.CachedQuery, minKey)
						}
					}
				} else {
					systemcommsin.Println("Put failed")
				}
			}
		}

	}
}

/*
Upon receiving a PUT message, or signal, it will simply
 1. Put the entry into local storage
 2. Call node.replicate(payload)
*/
func (node *Node) PutQuery(succesorId uint64, payload map[uint64][]string) bool {
	//systemcommsin.Println("Recieving a request to insert values into storage")
	_, ok := node.HashIPStorage[succesorId]
	if !ok {
		node.HashIPStorage[succesorId] = map[uint64][]string{}
	}
	for key, ip_cache := range payload {
		node.HashIPStorage[succesorId][key] = ip_cache
	}

	return true
}

/*
Replicate is called periodically to replicate all the storage entries to a new node. 
Replicated data is only sent to "REPLICATION_FACTOR" nodes
*/
func (node *Node) replicate () {
	for {
		time.Sleep(5 * time.Second)
		replicationSuccessor := make([]Pointer, REPLICATION_FACTOR)
		replicationSuccessor = append(replicationSuccessor, node.Successor)

		for i := 0; i < REPLICATION_FACTOR-1; i++ {
			succesor, _ := node.FindSuccessor(replicationSuccessor[len(replicationSuccessor)-1].Nodeid, 0)
			replicationSuccessor = append(replicationSuccessor, succesor)
		}

		for _, pointer := range replicationSuccessor {
			if (pointer.IP == node.IP || pointer == Pointer{}) {
				continue
			}
			msg := message.RequestMessage{Type: REPLICATE, TargetId: node.Nodeid, Payload: node.HashIPStorage[node.Nodeid]}
			go node.CallRPC(msg, pointer.IP)
		}
	}
}

/*
Processes the REPLICATE Type message received.
1. If the node's entry is not there, then dump the entire payload there, as it is the only entry. 
2. If the node's entry already exists, then add the new keys to it
*/
func (node *Node) processReplicate(senderId uint64, payload map[uint64][]string) bool {
	_, ok := node.HashIPStorage[senderId]
	if !ok {
		node.HashIPStorage[senderId] = payload
	}

	for key, ip_cache := range payload {
		node.HashIPStorage[senderId][key] = ip_cache
	}
	return true
}

/*
Given a hashed website name, return the records associated with it if it exists, else return nil. 
*/
func (node *Node) GetQuery(hashedId uint64) []string { // unused
	ip_addr, ok := node.HashIPStorage[node.Nodeid][hashedId]
	if ok {
		return ip_addr
	} else {
		return nil
	}
}

/*
Called when a SHIFT message is received. This means that there are new nodes in the network. The node will 
ask you to handover all the entries that falls between you and it. This method helps process this logic. 
*/
func (node *Node) GetShiftRecords(prececId uint64) map[uint64][]string {
	returnPayload := make(map[uint64][]string)
	nodeStorage, ok := node.HashIPStorage[node.Nodeid]
	if ok {
		for hashedWebsite := range nodeStorage {
			if prececId >= hashedWebsite {
				returnPayload[hashedWebsite] = nodeStorage[hashedWebsite]
				delete(nodeStorage, hashedWebsite)
			}
		}
		return returnPayload
	} else {
		return nil
	}
}

/*
Write the entry to persistent storage within the container.
It opens file in write or (create and write) mode.
*/
func (node *Node) writeToStorage() {
	filePath := fmt.Sprintf("/app/data/%s.json", node.IP)
	jsonData, err := json.Marshal(node.HashIPStorage[node.Nodeid])
	if err != nil {
		fmt.Println(err)
		return
	}
	if node.Logging {
		fmt.Printf("JSON data: %s\n", jsonData)
	}
	// Write to the file, create it if it doesn't exist
	// Append to the file or create it if it doesn't exist
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("Error opening or creating the file: %v\n", err)
		return
	}

	// Write the content to the file
	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Printf("Error writing to the file: %v\n", err)
		return
	}
	if node.Logging {
		fmt.Printf("JSON data written to file: %s\n", filePath)
	}
	defer file.Close()
}

/*
Reads file from local container storage.
It opens file in read or (create and read) mode.
*/
func (node *Node) readFromStorage() {
	filePath := fmt.Sprintf("/app/data/%s.json", node.IP)

	// Open the file for reading
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("Error opening or creating the file for reading: %v\n", err)
		return
	}
	defer file.Close()
	var storage map[uint64]map[uint64][]string
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&storage)
	if err != nil {
		fmt.Printf("Error decoding the JSON data: %v\n", err)
		return
	}

	if node.Logging {
		fmt.Printf("Data read from file\n")
	}

	for key, value := range storage {
		fmt.Printf("Key: %v, Value: %v\n", key, value)
	}
	node.HashIPStorage = storage
}
