/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

/*
Provide the host structure
*/
package host

// A target host (on which tasks run)
type Host struct {
	// A name for the host.  This name is only for logging purposes, and
	// does not need to be able to resolve to an IP.
	Name     string
	// The address (which can be a hostname or an IP)
	Addr     string
	// The username to use to login to the host
	Username string
	// The password for the username, nil means no password, as opposed to
	// an empty password
	Password *string
	// The port on which to connect to the host
	PortNum  uint16
}
