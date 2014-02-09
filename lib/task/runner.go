/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

/*
Run tasks on the hosts and get results
*/
package task

import (
	"errors"
	"fmt"
	"github.com/bgmerrell/geto/lib/config"
	"github.com/bgmerrell/geto/lib/host"
	"github.com/bgmerrell/geto/lib/remote"
	"log"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type RunOutput struct {
	Stdout string
	Stderr string
	Err    error
}

type NRunningScriptsOutput struct {
	n   uint32 // The number of running scripts on the target
	err error
}

// Acquire a remote lock by creating a remote directory that acts as a lock.
// IMPORTANT: This assumes mkdir is atomic on the target filesystem
func acquireRemoteRunnerLock(conn remote.Remote, host host.Host) (stderr string, err error) {
	const RETRIES = 10
	const SLEEP_INTERVAL = 0.1
	c := config.GetParsedConfig()
	for i := 0; i < RETRIES; i++ {
		_, stderr, err = conn.Run(
			host,
			fmt.Sprintf("mkdir %s", c.RemoteLockPath),
			0)
		if err == nil {
			break
		}
	}
	if err != nil {
		err = errors.New("Failed to acquire remote lock: " + err.Error())
	}

	return stderr, err
}

// Remove the remote runner lock from the master side.  This is only to be
// used when an error is encountered that prevents the task script from being
// executed on the target and the lock has already been acquired.
func removeRemoteRunnerLock(conn remote.Remote, host host.Host) {
	const RETRIES = 20
	const SLEEP_INTERVAL = 0.1
	c := config.GetParsedConfig()
	var err error
	for i := 0; i < RETRIES; i++ {
		_, _, err = conn.Run(
			host,
			fmt.Sprintf("rm -r %s", c.RemoteLockPath),
			0)
		if err == nil {
			log.Printf("Remote lock removed")
			break
		}
	}
	if err != nil {
		log.Printf("Failed to remove remote lock: %s", err.Error())
	}
}

func createRemoteWorkPathDir(conn remote.Remote, host host.Host) (stderr string, err error) {
	// Create the remote work path on the target host in case it hasn't
	// been corrected yet.
	// XXX: Is there a better way to do this?  It would be nice to not
	// have this extra ssh session and command for every run.
	c := config.GetParsedConfig()
	_, stderr, err = conn.Run(
		host,
		fmt.Sprintf("mkdir -p %s", c.RemoteWorkPath),
		0)
	if err != nil {
		err = errors.New(fmt.Sprintf(
			"Failed to create remote work directory, %s: (%s)",
			c.RemoteWorkPath, err.Error()))
	}
	return stderr, err
}

func getWrapperTask(innerTask Task) (wrapperTask Task, err error) {
	remoteInnerTaskDirPath := innerTask.getRemoteDirPath()
	stdoutPath := filepath.Join(remoteInnerTaskDirPath, "stdout")
	stderrPath := filepath.Join(remoteInnerTaskDirPath, "stderr")
	c := config.GetParsedConfig()
	var timeoutString string
	if innerTask.Timeout > 0 {
		timeoutString = fmt.Sprintf("%ss", innerTask.Timeout)
	} else {
		timeoutString = "3650d" // effectively no timeout
	}
	wrapperScript := NewScriptWithCommands(
		"wrapper",
		[]string{
			"#!/bin/bash",
			// It would be nice to not have the "timeout"
			// dependency.
			fmt.Sprintf("timeout --kill-after=10 %s %s 1>%s 2>%s &",
				timeoutString,
				innerTask.getRemoteScriptPath(),
				stdoutPath,
				stderrPath),
			fmt.Sprintf("rm -r %s", c.RemoteLockPath)},
		nil)
	return New([]string{}, wrapperScript, 0)
}

func getRemoteNRunningScripts(conn remote.Remote, task Task, host host.Host, ch chan<- NRunningScriptsOutput) {
	c := config.GetParsedConfig()
	// ^ is used to avoid matching the wrapper timeout process
	var pgrepPattern = fmt.Sprintf(
		"^/bin/bash %s/.*_%s", c.RemoteWorkPath, task.Script.name)
	stdout, _, err := conn.Run(
		host,
		// pgrep -c returns 1 if the count is zero (which
		// seems silly); we append the "; true" to work around
		// this.
		fmt.Sprintf("pgrep -c -f \"%s\"; true", pgrepPattern),
		0)
	n, err := strconv.ParseUint(strings.TrimSpace(stdout), 10, 32)
	ch <- NRunningScriptsOutput{uint32(n), err}
}

// Run a task on a target host
func RunOnHost(conn remote.Remote, task Task, host host.Host, resultChan chan<- RunOutput) {
	log.Printf("Running task %s on host %s (%s)...", task.Id, host.Name, host.Addr)

	// TODO: Better handle removeRemoteRunnerLock failures!  It might be
	// a problem that we just log a message when such a failure occurs,
	// since it means the user will have to go in and delete the lock
	// manually.  In generally, this code could probably benefit from more
	// specific types of errors.

	c := config.GetParsedConfig()
	taskDirPath, err := task.CreateDir()
	if err != nil {
		resultChan <- RunOutput{"", "", err}
		return
	}

	// Acquire the remote lock; if we fail after this, we need to make
	// sure the remote lock is removed.
	if stderr, err := acquireRemoteRunnerLock(conn, host); err != nil {
		resultChan <- RunOutput{"", stderr, err}
		return
	} else {
		log.Printf("%s acquired remote lock", task.Id)
	}

	if task.Script.maxConcurrent != nil {
		ch := make(chan NRunningScriptsOutput)
		go getRemoteNRunningScripts(conn, task, host, ch)
		nRunningScriptsOutput := <-ch
		if nRunningScriptsOutput.err != nil {
			removeRemoteRunnerLock(conn, host)
			resultChan <- RunOutput{"", "", errors.New("Failed to parse pgrep output: " + err.Error())}
		}
		if nRunningScriptsOutput.n >= *task.Script.maxConcurrent {
			removeRemoteRunnerLock(conn, host)
			resultChan <- RunOutput{"", "", errors.New(fmt.Sprintf(
				"Max concurrent (%d) \"%s\" scripts already running",
				nRunningScriptsOutput.n, task.Script.name))}
			return
		}
	}

	stderr, err := createRemoteWorkPathDir(conn, host)
	if err != nil {
		resultChan <- RunOutput{"", stderr, err}
		removeRemoteRunnerLock(conn, host)
		return
	}

	conn.CopyTo(host, true, taskDirPath, c.RemoteWorkPath)

	wrapperTask, err := getWrapperTask(task)
	if err != nil {
		removeRemoteRunnerLock(conn, host)
		resultChan <- RunOutput{"", stderr, err}
		return
	}

	log.Printf("Wrapper task: %s", wrapperTask.Id)
	wrapperTaskDirPath, err := wrapperTask.CreateDir()
	if err != nil {
		removeRemoteRunnerLock(conn, host)
		resultChan <- RunOutput{"", "", err}
		return
	}

	conn.CopyTo(host, true, wrapperTaskDirPath, c.RemoteWorkPath)

	stdout, stderr, err := conn.Run(
		host, wrapperTask.getRemoteScriptPath(), wrapperTask.Timeout)
	resultChan <- RunOutput{stdout, stderr, err}
}

func RunOnHostBalancedByScriptName(conn remote.Remote, task Task, ch chan<- RunOutput) {
	c := config.GetParsedConfig()
	hostToChan := map[host.Host](chan NRunningScriptsOutput){}
	var bestHost host.Host
	var minScriptsRunning uint32
	var failure error = nil
	for _, host := range c.Hosts {
		log.Printf("candidate: %s\n", host.Name)
	}
	for _, host := range c.Hosts {
		hostToChan[host] = make(chan NRunningScriptsOutput)
		host := host // new instance for go routine
		go func() {
			if _, err := acquireRemoteRunnerLock(conn, host); err != nil {
				log.Printf("%s acquired remote lock", task.Id)
			}
			go getRemoteNRunningScripts(conn, task, host, hostToChan[host])
		}()
	}
	for host, _ := range hostToChan {
		nRunningScriptsOutput := <-hostToChan[host]
		if nRunningScriptsOutput.err != nil {
			failure = nRunningScriptsOutput.err
		}
		log.Printf("%d scripts running on %s\n", nRunningScriptsOutput.n, host.Name)
		if nRunningScriptsOutput.n < minScriptsRunning || bestHost.Name == "" {
			bestHost = host
			minScriptsRunning = nRunningScriptsOutput.n
		}
		// TODO: parallelize lock removals
		removeRemoteRunnerLock(conn, host)
	}

	if failure != nil {
		ch <- RunOutput{"", "", errors.New("Failed : " + failure.Error())}
	} else {
		log.Printf("Selected host \"%s\" for load balancing", bestHost.Name)
		RunOnHost(conn, task, bestHost, ch)
	}
}

func RunOnRandomHost(conn remote.Remote, task Task, ch chan<- RunOutput) {
	RunOnHost(conn, task, getRandomHost(), ch)
}

func getRandomHost() host.Host {
	c := config.GetParsedConfig()
	return c.Hosts[rand.Intn(len(c.Hosts))]
}
