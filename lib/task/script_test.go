/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package task

import (
	"io/ioutil"
	"os"
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
		t.Errorf("Unexpected script contents:\n" +
			 "Actual: %#v\n" + 
			 "Expected: %#v",
			 script.commands,
		  	 expected)
	}
}

func TestScriptFileFrom(t *testing.T) {
	var existing_script_path string = filepath.Join(TESTDATADIR, "script.sh")
	script, err := NewScriptFromPath("test", existing_script_path)
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
	if fmt.Sprintf("%#v", string(expected)) == fmt.Sprintf("%#v", string(actual)) {
		t.Errorf("Unexpected script file contents:\n" +
			 "Actual:\n" + 
			 "---------\n" +
			 "%s" +
			 "\nExpected:\n" +
			 "---------\n" +
			 "%s",
			 string(actual),
		  	 string(expected))
	}
}
