/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"flag"
	"fmt"
	"github.com/bgmerrell/geto/lib/config"
	"github.com/bgmerrell/geto/lib/remote"
	"github.com/bgmerrell/geto/lib/remote/ssh"
	"os"
)

const SCP_TEST_PATH = "/tmp/geto-scp-test.txt"

/* Set by command line parsing */
var configPath string

var conf config.Config
var r remote.Remote

func parseCommandLine() {
	/* TODO: look for a system-wide config file in a portable manner */
	flag.StringVar(&configPath, "config-path", "geto.ini", "Configuration file path")
	flag.Parse()
}

func testConnection() {
	fmt.Println("Testing SSH connectivity...")
	for _, host := range conf.Hosts {
		fmt.Printf("%s@%s:%d : ", host.Username, host.Addr, host.PortNum)
		if err := r.TestConnection(host); err == nil {
			fmt.Printf("PASS\n")
		} else {
			fmt.Printf("FAIL (%s)\n", err.Error())
		}
	}
}

func testCopyToRemote() (err error) {
	fmt.Println("Testing SCP of local file to remote host(s)...")
	f, err := os.Create(SCP_TEST_PATH)
	defer f.Close()
	defer os.Remove(f.Name())
	if err != nil {
		fmt.Printf("Failed to open %s: %s", SCP_TEST_PATH, err.Error())
		return
	}
	if _, err = f.Write([]byte("Testing 1, 2, 3\n")); err != nil {
		fmt.Printf("Failed to write to %s: %s", SCP_TEST_PATH, err.Error())
		return
	}
	for _, host := range conf.Hosts {
		fmt.Printf("%s@%s:%d : ", host.Username, host.Addr, host.PortNum)
		err = r.CopyTo(
			host,
			false,
			SCP_TEST_PATH,
			SCP_TEST_PATH)
		if err != nil {
			fmt.Printf("FAIL (%s)\n", err.Error())
		} else {
			fmt.Printf("PASS\n")
		}
	}
	return nil
}

func testCopyFromRemote() (err error) {
	fmt.Println("Testing SCP of remote file to localhost...")
	for _, host := range conf.Hosts {
		fmt.Printf("%s@%s:%d : ", host.Username, host.Addr, host.PortNum)
		err = r.CopyFrom(
			host,
			false,
			SCP_TEST_PATH,
			SCP_TEST_PATH)
		/* Clean up remote side, we don't care too much if it fails */
		r.Run(
			host,
			fmt.Sprintf("rm %s", SCP_TEST_PATH),
			0)
		if err != nil {
			fmt.Printf("FAIL (%s)\n", err.Error())
		} else {
			fmt.Printf("PASS\n")
		}
	}
	return nil
}

func testRemoteEcho() {
	var stdout, stderr string
	var command string = "echo -n test"
	var err error
	fmt.Println("Testing remote echo...")
	for _, host := range conf.Hosts {
		fmt.Printf("%s@%s:%d : ", host.Username, host.Addr, host.PortNum)
		stdout, stderr, err = r.Run(
			host,
			command,
			0)
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
	r = ssh.New()
	testConnection()
	fmt.Println("")
	testRemoteEcho()
	fmt.Println("")
	testCopyToRemote()
	fmt.Println("")
	testCopyFromRemote()
}
