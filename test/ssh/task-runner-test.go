/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"flag"
	"fmt"
	"github.com/bgmerrell/geto/lib/config"
	"github.com/bgmerrell/geto/lib/host"
	"github.com/bgmerrell/geto/lib/remote/ssh"
	"github.com/bgmerrell/geto/lib/task"
	"os"
)

// Set by command line parsing
var configPath string

var conf config.Config
var testHost host.Host

func parseCommandLine() {
	// TODO: look for a system-wide config file in a portable manner
	flag.StringVar(&configPath, "config-path", "geto.ini", "Configuration file path")
	flag.Parse()
}

func testSleepTask() {
	fmt.Printf("Running a basic \"sleep\" task on %s...\n", testHost.Name)
	var script task.Script = task.NewScriptWithCommands(
		"sleep", []string{"#!/bin/bash", "sleep 5"}, nil)
	var depFiles []string
	t, err := task.New(depFiles, script, 0)
	if err != nil {
		fmt.Printf("FAIL (failed to create Task: %s)\n", err.Error())
		return
	}

	stdout, stderr, err := task.RunOnHost(ssh.New(), t, testHost)

	fmt.Println("stdout: ", stdout)
	fmt.Println("stderr: ", stderr)
	if err != nil {
		fmt.Println("err: ", err.Error())
	}
}

func main() {
	parseCommandLine()
	var err error
	if conf, err = config.ParseConfig(configPath); err != nil {
		os.Exit(1)
	}

	if len(conf.Hosts) < 1 {
		fmt.Println("No hosts found.")
		return
	}
	testHost = conf.Hosts[0]
	testSleepTask()
}
