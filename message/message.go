/*
Defines the structure for the Request and Response messages.
*/
package message

import (
	"fmt"
)

// Sample message structure. To be replaced with a struct for protobuff
type RequestMessage struct {
	Type     string // PING | SYNC | FIND_SUCCESSOR | CLOSEST_PRECEDING_NODE | PUT
	TargetId uint64 // ID of the parameter node passed to the destination
	IP       string // IP of the parameter node passed to the destination
	Payload  map[uint64][]string
	HopCount int
}

type ResponseMessage struct {
	Type   string // PING | SYNC | ACK | FIND_SUCCESSOR | CLOSEST_PRECEDING_NODE
	Nodeid uint64 // ID of the node in the response message
	IP     string // IP of the node in the response message
	QueryResponse []string
}

/*
***************************************
		UTILITY FUNCTIONS
***************************************
*/
// Message utility function to print the request message
func (msg *RequestMessage) PrintContent() {
	fmt.Println("Message content:", msg)
}

// Message utility function to print the response message
func (msg *ResponseMessage) PrintContent() {
	fmt.Println("Message content:", msg)
}
