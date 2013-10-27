/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

/*
Provide the host structure
*/
package host

type Host struct {
	Name     string
	Addr     string
	Username string
	/* nil means no password, as opposed to an empty password */
	Password *string
	PortNum  uint16
}
