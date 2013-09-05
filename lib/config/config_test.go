/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package config

import (
	"fmt"
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

func TestParseConfigHosts(t *testing.T) {
	fmt.Println("%v", conf.Hosts)

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
		fmt.Printf("%s: %s\n", host.name, host.addr)
		var expectedHostAddr string
		var ok bool
		if expectedHostAddr, ok = expected[host.name]; !ok {
			t.Errorf("Unexpected host name: %s", host.name)
		}
		if expectedHostAddr != host.addr {
			t.Errorf("Expected host addr \"%s\" for host name \"%s\", got \"%s\"",
				expectedHostAddr, host.name, host.addr)
		}
	}
}
