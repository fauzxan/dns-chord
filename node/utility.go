package node

import (
	"net/rpc"

	"github.com/fauzxan/dns-chord/v2/message"
	"github.com/rs/zerolog/log"
)

/*
***************************************************
		UTILITY FUNCTIONS FOR NODES
***************************************************
*/

/*
Node utility function to call RPC given a request message, and a destination IP address.
*/
func (node *Node) CallRPC(msg message.RequestMessage, IP string) message.ResponseMessage {
	log.Debug().Msgf("Nodeid: %d IP: %s is sending message %v to IP: %s", node.Nodeid, node.IP, msg, IP)
	clnt, err := rpc.Dial("tcp", IP)
	reply := message.ResponseMessage{}
	if err != nil {
		log.Error().Err(err).Msg(msg.Type)
		log.Debug().Msgf("Nodeid: %d IP: %s received reply %v from IP: %s", node.Nodeid, node.IP, reply, IP)
		reply.Type = EMPTY
		return reply
	}
	err = clnt.Call("Node.HandleIncomingMessage", msg, &reply)
	if err != nil {
		log.Error().Err(err).Msg("Error calling RPC")
		log.Debug().Msgf("Nodeid: %d IP: %s received reply %v from IP: %s", node.Nodeid, node.IP, reply, IP)
		reply.Type = EMPTY
		return reply
	}
	log.Debug().Msgf("Received reply from %s", IP)
	return reply
}

/*
Node utility function to print fingers
*/
func (node *Node) PrintFingers() {
	log.Info().Msg("Finger Table:")
	for i := 0; i < len(node.FingerTable); i++ {
		log.Info().Msgf("> Finger[%d]: Nodeid: %d IP: %s", i+1, node.FingerTable[i].Nodeid, node.FingerTable[i].IP)
	}
}

/*
Node utility function to print the successor
*/
func (node *Node) PrintSuccessor() {
	log.Info().Msg("Successor:")
	log.Info().Msgf(">Nodeid: %d Successor.IP: %s", node.Successor.Nodeid, node.Successor.IP)
}

/*
Node utility function to print predecessor
*/
func (node *Node) PrintPredecessor() {
	log.Info().Msg("Predecessor:")
	log.Info().Msgf(">Nodeid: %d Predecessor.IP: %s", node.Predecessor.Nodeid, node.Predecessor.IP)
}

func (node *Node) PrintStorage() {
	log.Info().Msg("STORAGE TABLE REQUESTED")
	log.Info().Msg("Storage:")
	for id, storage := range node.HashIPStorage {
		log.Info().Msgf(">id: %d", id)
		for _, value := range storage {
			log.Info().Msgf(">>value: %s", value)
		}
	}
}

func (node *Node) PrintCache() {
	log.Info().Msg("CACHE TABLE REQUESTED")
	for id, cache := range node.CachedQuery {
		log.Info().Msgf(">id: %d value: %s", id, cache.value)
	}
}

/*
Node utility function to check if an ID is in a given range (a, b].
*/
func belongsTo(id, a, b uint64) bool {
	if a == b {
		return true
	}
	if a < b {
		return a < id && id <= b
	} else {
		return a < id || id <= b
	}
}

/*
Node utility function to check if an ID is in a given range (a, b).
*/
func between(id, a, b uint64) bool {
	if a == b {
		return true
	}
	if a < b {
		return a < id && id < b
	} else {
		return a < id || id < b
	}
}
