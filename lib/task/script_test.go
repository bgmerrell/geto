/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package task

import (
	"path/filepath"
	"testing"
	"fmt"
)

const TESTDATADIR = "../../test/data"

func TestNewScript(t *testing.T) {
	var s script_t = NewScript("test")
	if s.name != "test" {
		t.Errorf("Error constructing new script")
	}
}

func TestNewScriptFromPath(t *testing.T) {
	script, err := NewScriptFromPath("test", filepath.Join(TESTDATADIR, "script.sh"))
	if err != nil {
		t.Fatalf("Error constructing new script: %s", err.Error())
	}
	if script.name != "test" {
		t.Fatalf("Error constructing new script")
	}

	var expected []string = []string{
		"echo \"Hello World\"",
		"ls",
		"true",
		"false",
	}

	if fmt.Sprintf("%#v", expected) != fmt.Sprintf("%#v", script.commands) {
		t.Fatalf("Unexpected script contents:\n" +
			 "Got: %#v\n" + 
			 "Expected: %#v",
			 script.commands,
		  	 expected)
	}
}
