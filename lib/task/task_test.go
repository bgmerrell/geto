/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package task

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"testing"
)

const expectedPattern = "[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}"

var validId *regexp.Regexp

func init() {
	validId = regexp.MustCompile(expectedPattern)
}

func TestGenTaskId(t *testing.T) {
	const numToGen = 50
	var taskIds map[string]struct{} = make(map[string]struct{}, numToGen)
	var id string
	var err error
	for i := 0; i < numToGen; i++ {
		if id, err = genTaskId(); err != nil {
			t.Fatalf(err.Error())
		}
		matched := validId.MatchString(id)
		if !matched {
			t.Errorf("Task ID %s did not match expected pattern of %s",
				id,
				expectedPattern)
		}
		taskIds[id] = struct{}{}
	}
	if len(taskIds) != numToGen {
		t.Errorf("Expected %d unique IDs, got %d", numToGen, len(taskIds))
	}
}

func TestNewTask(t *testing.T) {
	var script script_t = NewScript("test", nil)
	var depFiles []string
	var task Task
	var err error
	if task, err = NewTask(depFiles, script); err != nil {
		t.Fatalf("Failed to create new Task: " + err.Error())
	}
	if !validId.MatchString(task.Id) {
		t.Errorf("Task ID %s did not match expected pattern of %s",
			task.Id,
			expectedPattern)
	}
}

func TestTaskToFile(t *testing.T) {
	var depFiles []string
	var existing_script_path string = filepath.Join(TESTDATADIR, "script.sh")
	var task Task
	script, err := NewScriptFromPath(TEST_NAME, existing_script_path, nil)
	if err != nil {
		t.Fatalf("Error constructing new script: %s", err.Error())
	}

	if task, err = NewTask(depFiles, script); err != nil {
		t.Fatalf("Failed to create new Task: " + err.Error())
	}

	taskFilePath, err := task.ToFile()
	if err != nil {
		t.Fatalf("Failed to create file from Task: " + err.Error())
	}
	fmt.Println(taskFilePath)

	expected, err := ioutil.ReadFile(taskFilePath)
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
