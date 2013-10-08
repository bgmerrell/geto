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
	"time"
)

const SCP_TEST_PATH = "/tmp/geto-scp-test.txt"

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

func testScpToRemote() (err error) {
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
		err = ssh.ScpTo(
			host.Addr,
			host.Username,
			host.PortNum,
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

func testScpFromRemote() (err error) {
	fmt.Println("Testing SCP of remote file to localhost...")
	for _, host := range conf.Hosts {
		fmt.Printf("%s@%s:%d : ", host.Username, host.Addr, host.PortNum)
		err = ssh.ScpFrom(
			host.Addr,
			host.Username,
			host.PortNum,
			false,
			SCP_TEST_PATH,
			SCP_TEST_PATH)
		/* Clean up remote side, we don't care too much if it fails */
		ssh.Run(
			host.Addr,
			host.Username,
			host.Password,
			conf.PrivKeyPath,
			host.PortNum,
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
		stdout, stderr, err = ssh.Run(
			host.Addr,
			host.Username,
			host.Password,
			conf.PrivKeyPath,
			host.PortNum,
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

func testTimeout() {
	/* All durations in seconds */
	const sleepDuration = 8
	/* padding for things over than the execution of the actual remote sleep */
	const padDuration = 2
	const timeoutDuration = 3

	var stdout, stderr string
	var err error
	fmt.Println("Testing timeout for remote run...")
	for _, host := range conf.Hosts {
		fmt.Printf("%s@%s:%d : ", host.Username, host.Addr, host.PortNum)
		start := time.Now()
		stdout, stderr, err = ssh.Run(
			host.Addr,
			host.Username,
			host.Password,
			conf.PrivKeyPath,
			host.PortNum,
			fmt.Sprintf("%s %d", "sleep", sleepDuration),
			timeoutDuration)
		elapsed := time.Since(start)
		if err != nil {
			fmt.Printf("FAIL (%s)\n", err.Error())
		} else if elapsed > timeoutDuration+padDuration {
			fmt.Printf("FAIL (took %.1f seconds, expected < %d)\n",
				elapsed.Seconds(),
				timeoutDuration+padDuration)
		} else if stdout == "" && stderr == "" {
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
	testTimeout()
	fmt.Println("")
	testConnection()
	fmt.Println("")
	testRemoteEcho()
	fmt.Println("")
	testScpToRemote()
	fmt.Println("")
	testScpFromRemote()
}
