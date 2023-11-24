package main

import (
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math"
	"net"
)

/*
***************************************************
		UTILITY FUNCTIONS FOR MAIN
***************************************************
*/

const M = 32

/*
	Function to generate the hash of the the input IP address
*/
func GenerateHash(input string) uint64{
	data := []byte(input)
	id := sha256.Sum256(data)
	unmoddedID := float64(binary.BigEndian.Uint64(id[:8]))
	modValue := float64(math.Pow(2, M))
	moddedID := math.Mod(unmoddedID, modValue)
	return uint64(moddedID)
}

/*
	Function to automatically get the outbound IP without user input in .env file
*/
func GetOutboundIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP
}

/*
	Function to get a port number that is currently not in use
*/
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
			return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
			return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}