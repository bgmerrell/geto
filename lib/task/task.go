/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

/*
Provide the task structure and functions.
*/
package task

import (
	"errors"
	"os/exec"
	"fmt"
	"github.com/bgmerrell/geto/lib/config"
	"os"
	"path/filepath"
)

// A task that runs on a target host
type Task struct {
	// A unique ID for the task
	Id       string
	// A list of files and/or directories that the task requires
	DepFiles []string
	// A script for the task to run
	Script   script_t
	// The number of seconds before giving up on a task after it has been
	// started
	Timeout  uint32
}

func NewTask(depFiles []string, script script_t, timeout uint32) (Task, error) {
	taskId, err := genTaskId()
	return Task{taskId, depFiles, script, timeout}, err
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

// Creates a directory from a task object.
// The directory contains everything needed to run the task (e.g., required
// files and the script to run).
// The path to the created file is returned
// If there is a problem a non-nil error is returned
func (t *Task) CreateDir() (path string, err error) {
	c := config.GetParsedConfig()
	taskDirPath := filepath.Join(c.LocalWorkPath, t.Id)
	taskDepsDirPath := filepath.Join(taskDirPath, "DEPS")
	// The directory name is the task ID; we also create a special
	// dependency subdirectory where any file dependencies will live.
	os.MkdirAll(taskDepsDirPath, 0755)
	if err != nil {
		return "", errors.New(fmt.Sprintf(
			"Failed to create task directory: %s", err.Error()))
	}

	// The script name is the task ID combined with the script name.
	scriptFilePath := filepath.Join(
		taskDirPath,
		fmt.Sprintf("%s_%s", t.Id, t.Script.name))
	f, err := os.Create(scriptFilePath)
	defer f.Close()
	if err != nil {
		return "", errors.New(fmt.Sprintf(
			"Failed to create script file: %s", err.Error()))
	}

	for _, c := range t.Script.commands {
		_, err := f.Write([]byte(c + "\n"))
		if err != nil {
			return "", errors.New(fmt.Sprintf(
				"Failed to write script file: %s", err.Error()))
		}
	}

	// Now copy over the file dependencies to this task directory
	for _, depFilePath := range t.DepFiles {
		cmd := exec.Command("cp", "-r", depFilePath, taskDepsDirPath)
		err := cmd.Run()
		if err != nil {
			return "", errors.New(fmt.Sprintf(
				"Failed to copy file dependencies: %s", err.Error()))
		}
	}

	return taskDirPath, nil
}
