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

type script_t struct {
	name string
	commands []string
}

func NewScript(name string) script_t {
	return script_t{name, []string{}}
}

func NewScriptWithCommands(name string, commands []string) script_t {
	return script_t{name, commands}
}

// Takes a name and a path to a shell script and returns a script_t object
func NewScriptFromPath(name string, path string) (script_t, error) {
	var s script_t = script_t{name, []string{}}

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
