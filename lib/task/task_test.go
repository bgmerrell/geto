/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package task

import (
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
