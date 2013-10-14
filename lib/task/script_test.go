/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package task

import (
	"fmt"
	"io/ioutil"
	"os"
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

func TestScriptFileFrom(t *testing.T) {
	var existing_script_path string = filepath.Join(TESTDATADIR, "script.sh")
	script, err := NewScriptFromPath(TEST_NAME, existing_script_path, nil)
	if err != nil {
		t.Fatalf("Error constructing new script: %s", err.Error())
	}
	temp_f, err := ioutil.TempFile("", "geto-test")
	if err != nil {
		t.Fatal("Failed to open temporary file: %s", err.Error())
	}
	temp_f.Close()
	defer os.Remove(temp_f.Name())

	script.ToFile(temp_f.Name())

	expected, err := ioutil.ReadFile(temp_f.Name())
	if err != nil {
		t.Fatal("Failed to read temporary script file: %s", err.Error())
	}
	actual, err := ioutil.ReadFile(existing_script_path)
	if err != nil {
		t.Fatal("Failed to read existing script file: %s", err.Error())
	}
	if fmt.Sprintf("%#v", string(expected)) != fmt.Sprintf("%#v", string(actual)) {
		t.Errorf("Unexpected script file contents:\n"+
			"Actual:\n"+
			"---------\n"+
			"%s"+
			"\nExpected:\n"+
			"---------\n"+
			"%s",
			string(actual),
			string(expected))
	}
}
