package message

import "fmt"

// Sample message structure. To be replaced with a struct for protobuff
type Message struct{
	Content string
}

/*
***************************************
		UTILITY FUNCTIONS
***************************************	
*/
func (msg *Message) PrintContent(){
	fmt.Println("Message content:", msg.Content)
}