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
	"github.com/bgmerrell/geto/lib/config"
	"os"
	"path/filepath"
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

// Creates a file from a task object.
// The path to the created file is returned
// If there is a problem a non-nil error is returned
func (t *Task) ToFile() (path string, err error) {
	c := config.GetParsedConfig()
	path = filepath.Join(
		c.LocalWorkPath,
		fmt.Sprintf("%s_%s", t.Id, t.Script.name))

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	for _, c := range t.Script.commands {
		_, err := f.Write([]byte(c + "\n"))
		if err != nil {
			return "", err
		}
	}
	return path, nil
}
