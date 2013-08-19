/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"flag"
	"github.com/bgmerrell/geto/config"
	"github.com/bgmerrell/geto/server"
	"os"
)

/* Variables set by command line parsing */
var configPath string

func parseCommandLine() {
	/* XXX: How to make this portable? */
	flag.StringVar(&configPath, "config-path", "/etc/geto.ini", "Configuration file path")
	flag.Parse()
}

func main() {
	parseCommandLine()
	if c := config.ParseConfig(configPath); !c.IsParsed {
		os.Exit(1)
	}
	if server.Serve() {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
