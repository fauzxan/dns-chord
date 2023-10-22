package main

import (
	"core/node"
	"encoding/json"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"time"
)


func extract_fingertable() map[int]string{
	data, err := os.ReadFile("fingertable.json")
	if err != nil {
		fmt.Println("Error reading file fingertable.json")
	}
	var extracteddata []map[string]string
	err = json.Unmarshal(data, &extracteddata)
	if (err!=nil){
		panic(err)
	}
	// fmt.Println(extracteddata)
	var fingertable = make(map[int]string)
	for _, node := range extracteddata{
		id, err := strconv.Atoi(node["Nodeid"])
		if err!=nil{
			panic(err) //  we will be changing ids to hash functions later. Need to remove this. 
		}
		fingertable[id] = node["IP"]
	}
	return fingertable
}


func main(){
	// fmt.Println("Yo")
	var me = node.Node{
		FingerTable: make(map[int]string),
	}
	me.FingerTable = extract_fingertable()
	fmt.Println("Finger", me.FingerTable)
	// TODO 1: Write code to set own id, that is not already in the me.fingertable. Then update fingertable.json with me.fingertable
	address, err := net.ResolveTCPAddr("tcp", me.FingerTable[1]) // Change this to own id once TODO 1 is done
	if err != nil {
		fmt.Println("Error resolving TCP address")
	}
	inbound, err := net.ListenTCP("tcp", address)
	if err != nil {
		fmt.Println("Could not listen to TCP address")
	}
	rpc.Register(&me)
	fmt.Println("Node is runnning at IP address", address)
	go rpc.Accept(inbound)
	for {
		time.Sleep(1000)
	}
}