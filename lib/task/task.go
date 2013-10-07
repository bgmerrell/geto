/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

/*
Provide the task structure and functions.
*/
package task

import (
	"errors"
	"fmt"
	"os"
)

type Task struct {
	Id       string
	DepFiles []string
	Script   script_t
}

func NewTask(depFiles []string, script script_t) (Task, error) {
	taskId, err := genTaskId()
	return Task{taskId, depFiles, script}, err
}

// Generate a new task ID
func genTaskId() (string, error) {
	// 8 bytes should be good enough
	const numBytes = 8
	const failPrefix = "Failed to generate task ID: "
	f, err := os.Open("/dev/urandom")
	if err != nil {
		return "", errors.New(failPrefix + err.Error())
	}
	b := make([]byte, numBytes)
	count, err := f.Read(b)
	if err != nil {
		return "", errors.New(failPrefix + err.Error())
	} else if count != numBytes {
		return "", errors.New(fmt.Sprintf(
			"%sRead %d bytes, expected %d", failPrefix, count, numBytes))
	}
	f.Close()
	uuid := fmt.Sprintf("%x-%x-%x-%x", b[0:2], b[2:4], b[4:6], b[6:8])
	return uuid, nil
}
