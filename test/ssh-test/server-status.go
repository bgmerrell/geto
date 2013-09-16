/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"github.com/bgmerrell/geto/lib/ssh"
	"flag"
	"github.com/bgmerrell/geto/lib/config"
	"os"
)

/* Variables set by command line parsing */
var configPath string

func parseCommandLine() {
	/* TODO: look for a system-wide config file in a portable manner */
	flag.StringVar(&configPath, "config-path", "geto.ini", "Configuration file path")
	flag.Parse()
}

func main() {
	parseCommandLine()
	var conf config.Config
	var err error
	if conf, err = config.ParseConfig(configPath); err != nil {
		os.Exit(1)
	}

	if ok := ssh.TestDial(conf.PrivKeyPath); !ok {
		os.Exit(1)
	}
}
