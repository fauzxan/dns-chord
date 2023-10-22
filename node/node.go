package node

/*
	Main logic for:
	1. node joining
	2. planned departure
	3. find_successor
	4. closest_preceding_node
*/
type Node struct{
	Nodeid int // ID of the node in the fingertable
	IP string // localhost or IP address AND port number. Can be set through environment variables
	FingerTable map[int]string
	Successor string // IP of it's direct successor.
}



/*
***************************************
		UTILITY FUNCTIONS
***************************************	
*/
// Write utility functions here