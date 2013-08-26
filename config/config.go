/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package config

import (
	"fmt"
	"github.com/robfig/config"
	"log"
	"os"
)

var conf Config

func init() {
	conf = Config{}
}

type Config struct {
	IsParsed bool
	FilePath string
	Hosts    []string
}

func ParseConfig(configPath string) Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Print("No configuration file: ", configPath)
		return conf
	}
	log.Print("Parsing configuration file: ", configPath)
	/* TODO: look for system config first */
	/* TODO: How do I read a file with a relative path from this file? */
	c, err := config.ReadDefault(configPath)
	if err != nil {
		log.Print("Failed to parse config file: ", err.Error())
		return conf
	}
	/* TODO: replace test code with implementation */
	host, err := c.String("hosts", "localhost")
	if err != nil {
		log.Print("Failed to parse config file: ", err.Error())
		return conf
	}
	conf.IsParsed = true
	conf.FilePath = configPath
	conf.Hosts = append(conf.Hosts, host)
	fmt.Println("Host:", conf.Hosts[0])
	return conf
}
