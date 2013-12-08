/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package task

import (
	"fmt"
	"path/filepath"
	"testing"
)

const TESTDATADIR = "../../test/data"
const TEST_NAME = "test"

func TestNewScript(t *testing.T) {
	var s script_t = NewScript(TEST_NAME, nil)
	if s.name != TEST_NAME {
		t.Errorf(fmt.Sprintf(
			"Error constructing new script: Expected name \"%s\", got \"%s\"", TEST_NAME, s.name))
	}
	if s.maxConcurrent != nil {
		t.Errorf(fmt.Sprintf(
			"Error constructing new script: Expected \"nil\" maxConcurrent, got \"%d\"", s.maxConcurrent))
	}
}

func TestNewScriptFromPath(t *testing.T) {
	TEST_MAX_CONCURRENT := uint32(100)
	s, err := NewScriptFromPath(TEST_NAME, filepath.Join(TESTDATADIR, "script.sh"), &TEST_MAX_CONCURRENT)
	if err != nil {
		t.Fatalf("Error constructing new script: %s", err.Error())
	}
	if s.name != TEST_NAME {
		t.Fatalf(fmt.Sprintf(
			"Error constructing new script: Expected name \"%s\", got \"%s\"", TEST_NAME, s.name))
	}
	if s.maxConcurrent != nil && *s.maxConcurrent != TEST_MAX_CONCURRENT {
		t.Fatalf(fmt.Sprintf(
			"Error constructing new script: Expected \"%d\" maxConcurrent, got \"%d\"", TEST_MAX_CONCURRENT, s.maxConcurrent))
	}

	var expected []string = []string{
		"echo \"Hello World\"",
		"ls",
		"true",
		"false",
	}

	if fmt.Sprintf("%#v", expected) != fmt.Sprintf("%#v", s.commands) {
		t.Errorf("Unexpected script contents:\n"+
			"Actual: %#v\n"+
			"Expected: %#v",
			s.commands,
			expected)
	}
}
