/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"flag"
	"fmt"
	"github.com/bgmerrell/geto/lib/config"
	"github.com/bgmerrell/geto/lib/ssh"
	"os"
)

/* Set by command line parsing */
var configPath string

var conf config.Config

func parseCommandLine() {
	/* TODO: look for a system-wide config file in a portable manner */
	flag.StringVar(&configPath, "config-path", "geto.ini", "Configuration file path")
	flag.Parse()
}

func testConnection() {
	fmt.Println("Testing SSH connectivity...")
	for _, host := range conf.Hosts {
		fmt.Printf("%s@%s:%d : ", host.Username, host.Addr, host.PortNum)
		if err := ssh.TestConnection(host.Addr, host.Username, host.Password, conf.PrivKeyPath, host.PortNum); err == nil {
			fmt.Printf("PASS\n")
		} else {
			fmt.Printf("FAIL (%s)\n", err.Error())
		}
	}
}

func testRemoteEcho() {
	var stdout, stderr string
	var command string = "echo -n test"
	var err error
	fmt.Println("Testing remote echo...")
	for _, host := range conf.Hosts {
		fmt.Printf("%s@%s:%d : ", host.Username, host.Addr, host.PortNum)
		stdout, stderr, err = ssh.Run(
			host.Addr,
			host.Username,
			host.Password,
			conf.PrivKeyPath,
			host.PortNum,
			command)
		if err != nil {
			fmt.Printf("FAIL (%s)\n", err.Error())
		} else if stdout == "test" && stderr == "" {
			fmt.Printf("PASS\n")
		} else {
			fmt.Printf("FAIL (stdout: %s, stderr: %s)\n", stdout, stderr)
		}
	}
}

func main() {
	parseCommandLine()
	var err error
	if conf, err = config.ParseConfig(configPath); err != nil {
		os.Exit(1)
	}
	testConnection()
	fmt.Println("")
	testRemoteEcho()
}
