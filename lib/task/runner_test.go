/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package task

import (
	"fmt"
	"github.com/bgmerrell/geto/lib/config"
	"github.com/bgmerrell/geto/lib/remote/dummy"
	"sort"
	"testing"
)

func init() {
	if _, err := config.ParseConfig("../../test/data/geto.ini"); err != nil {
		panic("Failed to parse test config file.")
	}
}

func TestRunOnRandomHost(t *testing.T) {
	dummyConn := dummy.New()
	task := Task{"test-task", []string{}, NewScript("test-script", nil), 0}
	ch := make(chan RunOutput)
	go RunOnRandomHost(dummyConn, task, ch)
	_ = <-ch
}

func TestGetRandomHost(t *testing.T) {
	m := map[string]struct{}{}
	c := config.GetParsedConfig()
	const TRIES = 100
	for i := 0; i < TRIES; i++ {
		m[getRandomHost().Name] = struct{}{}
	}
	/* statistically, this should be true  */
	if len(m) != len(c.Hosts) {
		actual := make([]string, len(m))
		var expected []string
		for i, host := range c.Hosts {
			actual[i] = host.Name
		}
		fmt.Println(actual)
		for hostname := range m {
			expected = append(expected, hostname)
		}
		fmt.Println(expected)
		sort.Strings(actual)
		sort.Strings(expected)
		/* Just print out the expect and actual as sorted lists to be compared by the user */
		t.Errorf("Failed to randomly get all hosts in %d tries:\n"+
			"Actual: %#v\n"+
			"Expected: %#v",
			TRIES,
			actual,
			expected)
	}
}

func TestRunOnHostBalancedByScript(t *testing.T) {
	dummyConn := dummy.New()
	task := Task{"test-task", []string{}, NewScript("test-script", nil), 0}
	ch := make(chan RunOutput)
	go RunOnHostBalancedByScriptName(dummyConn, task, ch)
	<-ch
}
