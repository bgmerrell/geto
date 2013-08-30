/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

/*
Start a raw JSON RPC server

A client may call (via RPC) any of the GetoRPC functions exported here.
*/
package server

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type GetoRPC struct{}

// Echo a string
// The incoming parameter is the string to echo
// The echoed string is returned via the "outgoing" parameter
// This RPC is useful to validate that the RPC server is working
func (*GetoRPC) Echo(incoming *string, outgoing *string) error {
	log.Print("Echoing: ", *incoming)
	*outgoing = *incoming
	return nil
}

// Start the JSON RPC server
func Serve() bool {

	log.Print("Starting server...")
	listener, err := net.Listen("tcp", ":11102")
	if err != nil {
		log.Fatal("Failed to start server: %s\n", err.Error())
		return false
	}
	defer listener.Close()
	log.Print("Listening on: ", listener.Addr())
	rpc.Register(new(GetoRPC))
	for {
		log.Printf("Waiting for connection...")
		if conn, err := listener.Accept(); err == nil {
			log.Printf("Connection started: %v", conn.RemoteAddr())
			go jsonrpc.ServeConn(conn)
		} else {
			log.Fatal("Failed connection acceptance: ", err.Error())
		}
	}
}
