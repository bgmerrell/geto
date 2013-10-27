/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package dummy

import (
	"fmt"
	"github.com/bgmerrell/geto/lib/config"
	"github.com/bgmerrell/geto/lib/remote"
	"testing"
)

var r remote.Remote
var conf config.Config

func init() {
	var err error
	if conf, err = config.ParseConfig("../../../test/data/geto.ini"); err != nil {
		panic("Failed to parse config file: " + err.Error())
	}
}

func TestNewRemote(t *testing.T) {
	r = New()
}

func TestTestConnection(t *testing.T) {
	for _, host := range conf.Hosts {
		if err := r.TestConnection(host); err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestRun(t *testing.T) {
	const EXPECTED_STDOUT = "test"
	const EXPECTED_STDERR = ""
	for _, host := range conf.Hosts {
		fmt.Println(host)
		stdout, stderr, err := r.Run(host, "test", 0)
		fmt.Println("stdout: " + stdout)
		fmt.Println("stderr: " + stderr)
		if err != nil {
			t.Errorf(err.Error())
		} else if stdout != EXPECTED_STDOUT || stderr != EXPECTED_STDERR {
			t.Errorf("\nActual:\n"+
				"    stdout: %s\n"+
				"    stderr: %s\n"+
				"Expected:\n"+
				"    stdout: %s\n"+
				"    stderr: %s\n",
				stdout,
				stderr,
				EXPECTED_STDOUT,
				EXPECTED_STDERR)
		}
	}
}

func TestCopyTo(t *testing.T) {
	for _, host := range conf.Hosts {
		if err := r.CopyTo(host, true, "", ""); err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestCopyFrom(t *testing.T) {
	for _, host := range conf.Hosts {
		if err := r.CopyFrom(host, false, "", ""); err != nil {
			t.Errorf(err.Error())
		}
	}
}
