/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

/*
Geto's main package

Parse command line arguments and let the fun begin!
*/
package main

import (
	"flag"
	"github.com/bgmerrell/geto/lib/config"
	"github.com/bgmerrell/geto/server"
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
	if _, err := config.ParseConfig(configPath); err != nil {
		os.Exit(1)
	}
	if server.Serve() {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
