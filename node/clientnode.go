package node

/* 
	Inherits methods and properties of a normal chord node, but with added functionalities. 
	1. Be able to lookup hierarchical DNS
	2. Be able to receive and make queries. 
	TODO: modify code in main to check if the node joining is a client or not. We can have a flag for this in the arguments passed. 
*/

type ClientNode struct{
	Node // <-- inheriting from Node
	// extra properties
}




/*
***************************************
		UTILITY FUNCTIONS
***************************************
*/