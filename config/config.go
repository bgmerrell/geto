/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package config

import (
	"fmt"
	"github.com/robfig/config"
	"log"
)

func Config() bool {
	log.Print("Parsing configuration file...")
	/* TODO: look for system config first */
	/* TODO: How do I read a file with a relative path from this file? */
	c, err := config.ReadDefault("config/geto.ini")
	if err != nil {
		log.Print("Failed to parse config file:", err.Error())
		return false
	}
	/* TODO: replace test code with implementation */
	host, err := c.String("hosts", "localhost")
	if err != nil {
		log.Print("Failed to parse config file:", err.Error())
		return false
	}
	fmt.Println("Host:", host)
	return true
}
