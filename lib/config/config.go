/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

/*
Configuration file management
*/
package config

import (
	"errors"
	"github.com/bgmerrell/geto/lib/host"
	"github.com/bgmerrell/geto/lib/ssh"
	"github.com/bgmerrell/config"
	"log"
	"os"
	"strconv"
)

var conf Config
var isParsed bool = false

func init() {
	conf = Config{}
}

type Config struct {
	FilePath    string
	PrivKeyPath string
	RemoteWorkPath string
	LocalWorkPath string
	Hosts       []host.Host
}

// Parse the config file
// The configPath parameter is the path to the config file on the filesystem
func ParseConfig(configPath string) (Config, error) {
	conf = Config{} /* Zero out the package-scope conf */
	var err error
	var portNum int

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

	var privKeyPath string
	if privKeyPath, err = c.String("geto", "privkey_path"); err == nil {
		conf.PrivKeyPath = privKeyPath
	}

	if conf.RemoteWorkPath, err = c.String("geto", "remote_work_path"); err != nil {
		log.Print("Failed to parse remote work path: ", err.Error())
		return conf, err
	}

	if conf.LocalWorkPath, err = c.String("geto", "local_work_path"); err != nil {
		log.Print("Failed to parse local work path: ", err.Error())
		return conf, err
	}

	var opts []string
	if opts, err = c.Options("hosts"); err != nil {
		log.Print("Could not find \"hosts\" section: ", err.Error())
		return conf, err
	}

	const N_MIN_REQUIRED_HOSTS = 1
	if len(opts) < N_MIN_REQUIRED_HOSTS {
		err = errors.New("Config must have at least one host")
		log.Print("Failed to parse \"hosts\" section: ", err.Error())
		return conf, err
	}

	for _, hostname := range opts {
		var addr, username string
		var password *string = new(string)
		if addr, err = c.String("hosts", hostname); err != nil {
			log.Print("Failed to parse \"hosts\" section: ", err.Error())
			return conf, err
		}
		if opts, err = c.Options(hostname); err != nil {
			log.Print("Could not find \"", hostname, "\" section: ", err.Error())
			return conf, err
		}
		if username, err = c.String(hostname, "username"); err != nil {
			log.Print("Failed to parse \"username\" option for \"",
				hostname, "\" section: ", err.Error())
			return conf, err
		}
		if *password, err = c.String(hostname, "password"); err != nil {
			password = nil
			/* If there's no private key path, a username is required */
			if conf.PrivKeyPath == "" {
				log.Print("Failed to parse \"username\" option for \"",
					hostname, "\" section: ", err.Error())
				return conf, err
			}
		}
		if portNum, err = c.Int(hostname, "port"); err == nil {
			/* Max of 2^16 port numbers */
			if portNum>>16 != 0 {
				err = errors.New("Invalid port number: " + strconv.FormatUint(uint64(portNum), 10))
				return conf, err
			}
		} else {
			portNum = ssh.DEFAULT_SSH_PORT
		}
		conf.Hosts = append(
			conf.Hosts,
			host.Host{hostname, addr, username, password, uint16(portNum)})
	}

	conf.FilePath = configPath
	isParsed = true
	return conf, nil
}

// Return the parsed Config object.
// Panic if the config has not been parsed
// ParseConfig should be called before this function
func GetParsedConfig() Config {
	if !isParsed {
		panic("unparsed configuration")
	}
	return conf
}
