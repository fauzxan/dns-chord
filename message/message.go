package message

import (
	"fmt"
)

// Sample message structure. To be replaced with a struct for protobuff
type RequestMessage struct{
	Type string // PING | SYNC | ACK | FIND_SUCCESSOR | CLOSEST_PRECEDING_NODE
	TargetId uint64
}

type ResponseMessage struct {
	Type string
	Nodeid uint64
	IP string
}

/*
***************************************
		UTILITY FUNCTIONS
***************************************	
*/
func (msg *RequestMessage) PrintContent(){
	fmt.Println("Message content:", msg)
}

func (msg *ResponseMessage) PrintContent(){
	fmt.Println("Message content:", msg)
}