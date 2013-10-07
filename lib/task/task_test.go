/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package task

import (
	"regexp"
	"testing"
)

func TestGenTaskId(t *testing.T) {
	const numToGen = 50
	var taskIds map[string]bool = make(map[string]bool, numToGen)
	var expectedPattern = "[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}"
	var validId = regexp.MustCompile(expectedPattern)
	for i := 0; i < numToGen; i++ {
		id := genTaskId()
		matched := validId.MatchString(id)
		if !matched {
			t.Errorf("Task ID %s did not match expected pattern of %s",
				id,
				expectedPattern)
		}
		taskIds[id] = true
	}
	if len(taskIds) != numToGen {
		t.Errorf("Expected %d unique IDs, got %d", numToGen, len(taskIds))
	}
}
