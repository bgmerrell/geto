/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

/*
Provide the task structure and functions.
*/
package task

import (
	"os"
	"fmt"
)

type Task struct {
	Id       string
	DepFiles []string
	Script   script_t
}

// Generate a new task ID
func genTaskId() string {
	// 8 bytes should be good enough
	const numBytes = 8
        f, _ := os.Open("/dev/urandom")
        b := make([]byte, numBytes)
        f.Read(b)
        f.Close()
        uuid := fmt.Sprintf("%x-%x-%x-%x", b[0:2], b[2:4], b[4:6], b[6:8])
        return uuid
}
