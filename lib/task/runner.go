/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

/*
Run tasks on the hosts and get results
*/
package task

import (
	"fmt"
	"github.com/bgmerrell/geto/lib/config"
	"github.com/bgmerrell/geto/lib/host"
	"github.com/bgmerrell/geto/lib/remote"
	"log"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func RunOnHost(conn remote.Remote, task Task, host host.Host) (stdout string, stderr string, err error) {
	log.Printf("Running task %s on host %s (%s)...", task.Id, host.Name, host.Addr)
	c := config.GetParsedConfig()
	taskDirPath, err := task.CreateDir()
	if err != nil {
		return "", "", err
	}

	// Create the remote work path on the target host in case it hasn't
	// been corrected yet.  Is there a better way to do this?  It would be
	// nice to not have this extra ssh session and command for every run.
	_, stderr, err = conn.Run(
		host,
		fmt.Sprintf("mkdir -p %s", c.RemoteWorkPath),
		0)
	if err != nil {
		log.Printf("Failed to create remote work directory, %s: (%s)",
			c.RemoteWorkPath,
			err.Error())
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

func RunOnRandomHost(conn remote.Remote, task Task) {
	RunOnHost(conn, task, getRandomHost())
}

func getRandomHost() host.Host {
	c := config.GetParsedConfig()
	return c.Hosts[rand.Intn(len(c.Hosts))]
}
