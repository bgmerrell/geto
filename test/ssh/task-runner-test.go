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

	c := make(chan task.RunOutput)
	stdout, stderr, err := task.RunOnHost(ssh.New(), t, testHost, c)
	output := <- c1

	fmt.Println("stdout: ", output.Stdout)
	fmt.Println("stderr: ", output.Stderr)
	if err != nil {
		fmt.Println("err: ", output.Err.Error())
	}
}

func testMaxConcurrentSleepTasks() {
	fmt.Printf("Running two maxConcurrent=1 \"sleep\" tasks on %s...\n", testHost.Name)
	maxConcurrent := uint32(1)
	var script task.Script = task.NewScriptWithCommands(
		"sleep", []string{"#!/bin/bash", "sleep 60"}, &maxConcurrent)
	var depFiles []string
	t1, err := task.New(depFiles, script, 0)
	if err != nil {
		fmt.Printf("FAIL (failed to create task 1: %s)\n", err.Error())
		return
	}
	t2, err := task.New(depFiles, script, 0)
	if err != nil {
		fmt.Printf("FAIL (failed to create task 2: %s)\n", err.Error())
		return
	}

	c1 := make(chan task.RunOutput)
	c2 := make(chan task.RunOutput)
	go task.RunOnHost(ssh.New(), t1, testHost, c1)
	go task.RunOnHost(ssh.New(), t2, testHost, c2)

	output1 := <- c1
	output2 := <- c2

	fmt.Printf("Task 1 (%s)\n", t1.Id)
	fmt.Println("------------------------")
	fmt.Println("stdout: ", output1.Stdout)
	fmt.Println("stderr: ", output1.Stderr)
	if output1.Err != nil {
		fmt.Println("err: ", output1.Err.Error())
	}
	fmt.Println("")
	fmt.Printf("Task 2 (%s)\n", t2.Id)
	fmt.Println("------------------------")
	fmt.Println("stdout: ", output2.Stdout)
	fmt.Println("stderr: ", output2.Stderr)
	if output2.Err != nil {
		fmt.Println("err: ", output2.Err.Error())
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
	// testSleepTask()
	testMaxConcurrentSleepTasks()
}
