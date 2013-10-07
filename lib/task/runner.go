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
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func RunOnHost(task Task, host host.Host) {
	fmt.Println("TODO: Run task on host: " + host.Name)
}

func RunOnHostBalancedByScriptName(task Task) {
	fmt.Println("TODO: Try to find the host running the fewest scripts of the same name.")
}

func RunOnRandomHost(task Task) {
	RunOnHost(task, getRandomHost())
}

func getRandomHost() host.Host {
	c := config.GetParsedConfig()
	return c.Hosts[rand.Intn(len(c.Hosts))]
}
