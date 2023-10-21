package client

/*
	This file contains the main logic of the client.
	1. Receiving requests
	2. Making requests
	3. Updating fingertables
	4. Handling node failures, 
	5. etc.
*/

type Client struct{
	Clientid int // ID of the client in the fingertable
	IP string // localhost or IP address AND port number. Can be set through environment variables
	FingerTable map[int]string
}



/*
***************************************
		UTILITY FUNCTIONS
***************************************	
*/
// Write utility functions here