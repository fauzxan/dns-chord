package node

import "sync"

type FingerTable struct{
	Lock sync.Mutex // Used to lock the local copy of the fingertable
	Clientlist map[int]string // List of clients with the structure: {id: ip_addr. id: ip_addr ...}
}



/*
***************************************
		UTILITY FUNCTIONS
***************************************	
*/
// Write utility functions here