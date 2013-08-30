/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package config

import (
	"errors"
	"github.com/robfig/config"
	"log"
	"os"
)

var conf Config

func init() {
	conf = Config{}
}

type Host struct {
	name	string
	addr	string
}

type Config struct {
	FilePath string
	Hosts    []Host
}

func ParseConfig(configPath string) (Config, error) {
	var err error

	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		log.Print("No configuration file: ", configPath)
		return conf, err
	}
	log.Print("Parsing configuration file: ", configPath)

	var c *config.Config
	c, err = config.ReadDefault(configPath)
	if err != nil {
		log.Print("Failed to parse config file: ", err.Error())
		return conf, err
	}

	var opts []string
	if opts, err = c.Options("hosts"); err != nil {
		log.Print("Failed to parse hosts section: ", err.Error())
		return conf, err
	}

	const N_MIN_REQUIRED_HOSTS = 1
	if len(opts) < N_MIN_REQUIRED_HOSTS {
		err = errors.New("Config must have at least one host")
		log.Print("Failed to parse hosts section: ", err.Error())
		return conf, err
	}

	for _, hostname := range opts {
		var addr string
		if addr, err = c.String("hosts", hostname); err != nil {
			log.Print("Failed to parse hosts section: ", err.Error())
			return conf, err
		}
		conf.Hosts = append(conf.Hosts, Host{hostname, addr})
	}

	conf.FilePath = configPath
	return conf, nil
}
