/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package server

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type RPCFunc struct{}

func (*RPCFunc) Echo(arg *string, result *string) error {
	log.Print("Arg passed: " + *arg)
	*result = ">" + *arg + "<"
	return nil
}

func Serve() {
	log.Print("Starting server...")
	l, err := net.Listen("tcp", ":1234")
	defer l.Close()
	if err != nil {
		log.Fatal("Failed to start server: %s\n", err.Error())
	}
	log.Print("Listening on: ", l.Addr())
	rpc.Register(new(RPCFunc))
	for {
		log.Printf("Waiting for connection...")
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("Failed connection acceptance: ", err.Error())
		}
		log.Printf("Connection started: %v", conn.RemoteAddr())
		go jsonrpc.ServeConn(conn)
	}
}
