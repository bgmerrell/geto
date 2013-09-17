/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package config

import (
	"fmt"
	"strconv"
	"testing"
)

/* conf defined in config.go */
var err error

const TESTDATADIR = "../../test/data"

func TestParseEmptyConfig(t *testing.T) {
	if _, err = ParseConfig("../../test/data/config-empty.ini"); err == nil {
		t.Errorf("Parsing an empty config should fail")
	}
}

func TestParseConfigWithoutHosts(t *testing.T) {
	if _, err = ParseConfig("../../test/data/config-no-hosts.ini"); err == nil {
		t.Errorf("Parsing a config with no hosts should fail")
	}
}

func TestParseConfigWithMissingPassword(t *testing.T) {
	if _, err = ParseConfig("../../test/data/config-missing-password.ini"); err == nil {
		t.Errorf("Parsing a config with a missing password should fail")
		return
	}
	if err.Error() != "option not found: password" {
		t.Errorf("Expected to fail for missing password")
	}
}

func TestParseConfigWithBadPort(t *testing.T) {
	if _, err = ParseConfig("../../test/data/config-bad-port.ini"); err == nil {
		t.Errorf("Parsing a config with an invalid port should fail")
		return
	}
	if err.Error() != "Invalid port number: 123456789" {
		t.Errorf("Expected to fail for invalid port number")
	}
}

func TestParseConfigWithMissingHost(t *testing.T) {
	if _, err = ParseConfig("../../test/data/config-missing-host.ini"); err == nil {
		t.Errorf("Parsing a config with a missing host should fail")
		return
	}
	if err.Error() != "section not found: server3" {
		t.Errorf("Expected to fail for missing host \"server3\"")
	}
}

func TestParseMissingConfig(t *testing.T) {
	if _, err = ParseConfig("BOGUS-CONFIG.ini"); err == nil {
		t.Errorf("Attempting to parse a missing config should fail")
	}
}

func TestParseConfigPath(t *testing.T) {
	if conf, err = ParseConfig("../../test/data/geto.ini"); err != nil {
		t.Fatalf("Parse of good config should pass.")
	}

	expected := "../../test/data/geto.ini"
	actual := conf.FilePath
	if expected != actual {
		t.Errorf("Config file path (%s) does not match expected (%s)",
			actual, expected)
	}
}

func TestParsePrivKeyPath(t *testing.T) {
	expected := "/Users/bean/.ssh/y"
	actual := conf.PrivKeyPath
	if expected != actual {
		t.Errorf("Private key file path (%s) does not match expected (%s)",
			actual, expected)
	}
}

func TestParseConfigHosts(t *testing.T) {
	expected := map[string]string{
		"server1": "10.0.0.10",
		"server2": "server2.int.mydomain.com",
		"server3": "server3",
	}
	const N_HOSTS = 3
	if len(conf.Hosts) != N_HOSTS {
		t.Errorf("Expected %d hosts, got %d", N_HOSTS, len(conf.Hosts))
	}

	for _, host := range conf.Hosts {
		fmt.Printf("%s: %s\n", host.Name, host.Addr)
		var expectedHostAddr string
		var ok bool
		if expectedHostAddr, ok = expected[host.Name]; !ok {
			t.Errorf("Unexpected host name: %s", host.Name)
		}
		if expectedHostAddr != host.Addr {
			t.Errorf("Expected host addr \"%s\" for host name \"%s\", got \"%s\"",
				expectedHostAddr, host.Name, host.Addr)
		}
	}
}

func TestParseConfigPorts(t *testing.T) {
	expected := map[string]uint16{
		"server1": uint16(22),
		"server2": uint16(2222),
		"server3": uint16(22),
	}

	for _, host := range conf.Hosts {
		fmt.Printf("%s: %s\n", host.Name, host.Addr)
		var expectedPortNum uint16
		var ok bool
		if expectedPortNum, ok = expected[host.Name]; !ok {
			t.Errorf("Unexpected host name: %s", host.Name)
		}
		if expectedPortNum != host.PortNum {
			t.Errorf("Expected port number \"%s\" for host name \"%s\", got \"%s\"",
				strconv.FormatUint(uint64(expectedPortNum), 10),
				host.Name,
				strconv.FormatUint(uint64(host.PortNum), 10))
		}
	}
}
