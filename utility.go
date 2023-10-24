package main

import (
	"crypto/sha256"
	"encoding/binary"
)

/*
***************************************
		UTILITY FUNCTIONS FOR MAIN
***************************************
*/

func GenerateHash(input string) uint64{
	data := []byte(input)
	id := sha256.Sum256(data)
	return binary.BigEndian.Uint64(id[:8]) 
}