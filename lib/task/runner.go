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
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func acquireRemoteRunnerLock(conn remote.Remote, host host.Host) (stderr string, err error) {
	const RETRIES = 10
	const SLEEP_INTERVAL = 0.1
	// TODO: put LOCK_DIR in the config file
	const LOCK_DIR = "/var/tmp/geto_lock"
	for i := 0; i < RETRIES; i++ {
		_, stderr, err = conn.Run(
			host,
			fmt.Sprintf("mkdir %s", LOCK_DIR),
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
	wrapperScript := NewScriptWithCommands(
		"wrapper",
		[]string{
			"#!/bin/bash",
			fmt.Sprintf("%s 1>%s 2>%s &",
				innerTask.getRemoteScriptPath(),
				stdoutPath,
				stderrPath),
			fmt.Sprintf("rm -r %s", c.RemoteLockPath)},
		nil)
	return New([]string{}, wrapperScript, 0)
}

func RunOnHost(conn remote.Remote, task Task, host host.Host) (stdout string, stderr string, err error) {
	log.Printf("Running task %s on host %s (%s)...", task.Id, host.Name, host.Addr)
	c := config.GetParsedConfig()
	taskDirPath, err := task.CreateDir()
	if err != nil {
		return "", "", err
	}

	if stderr, err = acquireRemoteRunnerLock(conn, host); err != nil {
		return "", stderr, errors.New(
			"Failed to acquire remote lock: " + err.Error())
	}

	if stderr, err = createRemoteWorkPathDir(conn, host); err != nil {
		return "", stderr, err
	}

	conn.CopyTo(host, true, taskDirPath, c.RemoteWorkPath)

	return conn.Run(
		host,
		task.getRemoteScriptPath(),
		task.Timeout)
}

func RunOnHostBalancedByScriptName(conn remote.Remote, task Task) {
	fmt.Println("TODO: Try to find the host running the fewest scripts of the same name.")
}

func RunOnRandomHost(conn remote.Remote, task Task) (stdout string, stderr string, err error) {
	return RunOnHost(conn, task, getRandomHost())
}

func getRandomHost() host.Host {
	c := config.GetParsedConfig()
	return c.Hosts[rand.Intn(len(c.Hosts))]
}
