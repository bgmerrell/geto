/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package config

import (
	"fmt"
	"log"
	"github.com/robfig/config"
)

func Config() bool {
	log.Print("Parsing configuration file...")
	/* TODO: look for system config first */
	c, err := config.ReadDefault("./geto.ini")
	if err != nil {
		log.Print("Failed to parse config file: ", err.Error())
		return false;
	}
	fmt.Println(c.String("hosts", "localhost"))
	return true;
}
