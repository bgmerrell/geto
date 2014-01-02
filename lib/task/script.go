/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

/*
Provide the script structure and functions.
*/
package task

import (
	"bufio"
	"os"
)

// A script that runs on a target host
type Script struct {
	// Name is the name of a script.  It need not be unique.
	name string
	// The commands that make up a shell-style script.
	// Each index represents a line in the script.
	commands []string
	// The number of scripts of the same name that will run on a target host
	// concurrently.  A nil value means there is no limit.
	maxConcurrent *uint32
}

func NewScript(name string, maxConcurrent *uint32) Script {
	return Script{name, []string{}, maxConcurrent}
}

func NewScriptWithCommands(name string, commands []string, maxConcurrent *uint32) Script {
	return Script{name, commands, maxConcurrent}
}

// Takes a name and a path to a shell script and returns a Script object
func NewScriptFromPath(name string, path string, maxConcurrent *uint32) (Script, error) {
	var s Script = Script{name, []string{}, maxConcurrent}

	f, err := os.Open(path)
	if err != nil {
		return s, err
	}
	defer f.Close()

	/* TODO: Is there a better way to do this?  I need ot understand the bufio scanner
	(and Go scanners in general) better */
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s.commands = append(s.commands, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return s, err
	}
	return s, err
}
