package node

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"core.com/message"
	"core.com/utility"
)

func (node *Node) replicate(payload map[uint64][]string) {
	replicationSuccessor := make([]Pointer, REPLICATION_FACTOR)
	replicationSuccessor = append(replicationSuccessor, node.Successor)

	for i := 0; i < REPLICATION_FACTOR; i++ {
		succesor, _ := node.FindSuccessor(replicationSuccessor[len(replicationSuccessor)-1].Nodeid, 0)
		replicationSuccessor = append(replicationSuccessor, succesor)
	}

	for _, pointer := range replicationSuccessor {
		if (pointer.IP == node.IP || pointer == Pointer{}) {
			continue
		}
		msg := message.RequestMessage{Type: REPLICATE, TargetId: node.Nodeid, Payload: payload}
		node.CallRPC(msg, pointer.IP)
	}
}

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

func (node *Node) PutQuery(succesorid uint64, payload map[uint64][]string) bool {
	//systemcommsin.Println("Recieving a request to insert values into storage")
	_, ok := node.HashIPStorage[succesorid]
	if !ok {
		node.HashIPStorage[succesorid] = map[uint64][]string{}
	}
	for key, ip_cache := range payload {
		node.HashIPStorage[succesorid][key] = ip_cache
	}

	node.replicate(payload)
	return true
}

func (node *Node) GetQuery(hashedId uint64) []string { // unused
	ip_addr, ok := node.HashIPStorage[node.Nodeid][hashedId]
	if ok {
		return ip_addr
	} else {
		return nil
	}
}

// 1
// 1000, 2000, 3000
func (node *Node) QueryDNS(website string) {
	if node.CachedQuery == nil {
		node.CachedQuery = make(map[uint64]Cache)
	}
	node.Counter += 1

	if strings.HasPrefix(website, "www.") {
		system.Println("> Removing Prefix")
		website = website[4:]
	}
	hashedWebsite := utility.GenerateHash(website)
	ip_addr, ok := node.CachedQuery[hashedWebsite]
	if ok {
		system.Println("> Retrieving from Cache")
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
				node.CachedQuery[hashedWebsite] = Cache{value: ip_addresses, counter: node.Counter}
				reply = node.CallRPC(message.RequestMessage{Type: PUT, TargetId: succPointer.Nodeid, Payload: map[uint64][]string{hashedWebsite: ip_addresses}}, succPointer.IP)
				// system.Println(reply)
				// system.Println("Node HashIPStorage: ", node.HashIPStorage)

				if reply.Type == ACK {
					if len(node.CachedQuery) > CACHE_SIZE {
						var minKey uint64
						minValue := uint64(18446744073709551615)
						for key, value := range node.CachedQuery {
							if value.counter < minValue {
								minKey = key
								minValue = value.counter
							}
						}
						if minKey != 0 {
							delete(node.CachedQuery, minKey)
						}

					}
					// for key, value := range node.CachedQuery
					// 	system.Printf("Key: %d, Value: %s, %d\n", key, value.value[0], value.counter)
					// }
				} else {
					systemcommsin.Println("Put failed")
				}
			}
		}

	}
	// node.CachedQuery[website] = ip.String();

}

func (node *Node) writeToStorage() {

	filePath := fmt.Sprintf("/app/data/%s.json", node.IP)
	// content := fmt.Sprintf("%d : %v\n", hashedWebsite, ip_addresses)
	jsonData, err := json.Marshal(node.HashIPStorage)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("JSON data: %s\n", jsonData)
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

	fmt.Printf("JSON data written to file: %s\n", filePath)
	// _, err = file.Seek(0, 0)
	// if err != nil {
	// 	fmt.Printf("Error seeking to the beginning of the file: %v\n", err)
	// 	return
	// }
	// var storage map[uint64][]string
	// decoder := json.NewDecoder(file)
	// err = decoder.Decode(&storage)
	// if err != nil {
	// 	fmt.Printf("Error decoding the JSON data: %v\n", err)
	// 	return
	// }

	// fmt.Printf("Data read from file\n")
	// for key, value := range storage {
	// 	fmt.Printf("Key: %v, Value: %v\n", key, value)
	// }
	defer file.Close()
}

func (node *Node) readFromStorage() {
	filePath := fmt.Sprintf("/app/data/%s.json", node.IP)

	// Open the file for reading
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("Error opening the file for reading: %v\n", err)
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

	fmt.Printf("Data read from file\n")

	// When node crashes, node.HashIPStorage = storage

	for key, value := range storage {
		fmt.Printf("Key: %v, Value: %v\n", key, value)
	}
	node.HashIPStorage = storage
	defer file.Close()

}