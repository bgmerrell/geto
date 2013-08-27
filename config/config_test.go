/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package config

import (
	"testing"
)

func TestParseEmptyConfig(t *testing.T) {
	if _, err := ParseConfig("../test-data/config-empty.ini"); err == nil {
		t.Errorf("Parsing an empty config should fail")
	}
}

func TestParseConfigWithoutHosts(t *testing.T) {
	if _, err := ParseConfig("../test-data/config-no-hosts.ini"); err == nil {
		t.Errorf("Parsing a config with no hosts should fail")
	}
}
