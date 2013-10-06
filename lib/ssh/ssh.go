/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

/*
All of the calls to the external SSH library
(code.google.com/p/go.crypto/ssh) will go through this package.  This gives us
the opportunity to adjust the interface for our needs.  More importantly, it
will allow us to more easily swap out the backend if all of our external SSH
calls are in the same place.
*/
package ssh

import (
	"bytes"
	"code.google.com/p/go.crypto/ssh"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strconv"
)

const DEFAULT_SSH_PORT = 22

// keychain implements the ssh.ClientKeyring interface
type keychain struct {
	keys []ssh.Signer
}

func (k *keychain) Key(i int) (ssh.PublicKey, error) {
	if i < 0 || i >= len(k.keys) {
		return nil, nil
	}

	return k.keys[i].PublicKey(), nil
}

func (k *keychain) Sign(i int, rand io.Reader, data []byte) (sig []byte, err error) {
	return k.keys[i].Sign(rand, data)
}

func (k *keychain) add(key ssh.Signer) {
	k.keys = append(k.keys, key)
}

func (k *keychain) loadPEM(file string) error {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return err
	}
	k.add(key)
	return nil
}

// clientPassword implements the ssh.ClientPassword interface
type clientPassword string

func (p clientPassword) Password(user string) (string, error) {
	return string(p), nil
}

// Establish a code.google.com/p/go.crypto/ssh Session.
// The caller is responsible for closing the session.
func getSession(
	addr string,
	username string,
	password *string,
	privKeyPath string,
	portNum uint16) (session *ssh.Session, err error) {

	var authorizers []ssh.ClientAuth = []ssh.ClientAuth{}
	if privKeyPath != "" {
		var clientKeychain *keychain = new(keychain)
		if err := clientKeychain.loadPEM(privKeyPath); err != nil {
			return session, err
		}
		authorizers = append(
			authorizers, ssh.ClientAuthKeyring(clientKeychain))
	}

	if password != nil {
		authorizers = append(
			authorizers, ssh.ClientAuthPassword(clientPassword(*password)))
	}

	if len(authorizers) == 0 {
		return session, errors.New("No authorization methods provided")
	}

	/* Try to authenticate with a public SSH key first, try a password if that fails */
	config := &ssh.ClientConfig{
		User: username,
		Auth: authorizers,
	}
	client, err := ssh.Dial(
		"tcp",
		addr+":"+strconv.FormatUint(uint64(portNum), 10),
		config)
	if err != nil {
		return session, err
	}

	session, err = client.NewSession()
	if err != nil {
		return session, err
	}
	return session, err
}

func TestConnection(
	addr string,
	username string,
	password *string,
	privKeyPath string,
	portNum uint16) (err error) {

	var session *ssh.Session

	session, err = getSession(addr, username, password, privKeyPath, portNum)
	if err != nil {
		return err
	}
	defer session.Close()

	if err = session.Run("true"); err != nil {
		return err
	}
	return nil
}

// The addr parameter is the address (IP, hostname, etc) of the remote host.
// The username parameter is the username to use to SSH to the remote host.
// The password parameter is the password to use to SSH to the remote host.
// The privKeyPath parameter is the path to the private key of the master.
// The portNum parameter is the SSH port number of the remote host.
// The command parameter is the command to run on the remote host.
// The timeout parameter is the number of seconds before abandoning the command.
// A timeout of 0 means no timeout.
func Run(
	addr string,
	username string,
	password *string,
	privKeyPath string,
	portNum uint16,
	command string,
	timeout uint32) (stdout string, stderr string, err error) {

	var session *ssh.Session

	session, err = getSession(addr, username, password, privKeyPath, portNum)
	if err != nil {
		return "", "", err
	}
	defer session.Close()

	var stdout_buf bytes.Buffer
	var stderr_buf bytes.Buffer
	session.Stdout = &stdout_buf
	session.Stderr = &stderr_buf

	if timeout != 0 {
		// TODO: kill the remote process if there is a timeout
	}

	if err = session.Run(command); err != nil {
		return "", "", err
	}

	return stdout_buf.String(), stderr_buf.String(), err
}

// Secure copy (scp) from localhost to addr
// Run a separate scp process (for now) to secure copy files between hosts.
// The addr parameter is the address (IP, hostname, etc) of the remote host.
// The privKeyPath parameter is the path to the private key of the master.
// The portNum is the SSH port number of the remote host.
// The incoming parameter indicates which direction to perform the copy.
// The recursive parameter indicates whether to use the -r scp option.
// Password authentication not supported for this function.
// XXX: This function should probably go away in favor of a single Scp function when Issue #1 is fixed.
func ScpTo(
	addr string,
	username string,
	portNum uint16,
	recursive bool,
	localPath string,
	remotePath string) (err error) {

	/* Unfortunately, there doesn't appear to be an SFTP or SCP library, so
	we'll just have to run a separate scp process.  This means no password
	authentication for when calling this method. */
	var stdout bytes.Buffer
	var args []string = []string{}
	var cmd *exec.Cmd
	if recursive {
		args = append(args, "-r")
	}
	args = append(args, fmt.Sprintf("-P"))
	args = append(args, strconv.FormatUint(uint64(portNum), 10))
	args = append(args, localPath)
	args = append(args, fmt.Sprintf("%s@%s:%s", username, addr, remotePath))
	cmd = exec.Command("scp", args...)
	cmd.Stderr = &stdout
	err = cmd.Run()
	if err != nil {
		return errors.New("scp to " + addr + " failed: " + err.Error())
	}
	return nil

}

// Secure copy (scp) from addr to localhost
// Run a separate scp process (for now) to secure copy files between hosts.
// The addr parameter is the address (IP, hostname, etc) of the remote host.
// The privKeyPath parameter is the path to the private key of the master.
// The portNum is the SSH port number of the remote host.
// The incoming parameter indicates which direction to perform the copy.
// The recursive parameter indicates whether to use the -r scp option.
// Password authentication not supported for this function.
// XXX: This function should probably go away in favor of a single Scp function when Issue #1 is fixed.
func ScpFrom(
	addr string,
	username string,
	portNum uint16,
	recursive bool,
	remotePath string,
	localPath string) (err error) {

	/* Unfortunately, there doesn't appear to be an SFTP or SCP library, so
	we'll just have to run a separate scp process.  This means no password
	authentication for when calling this method. */
	var stdout bytes.Buffer
	var args []string = []string{}
	var cmd *exec.Cmd
	if recursive {
		args = append(args, "-r")
	}
	args = append(args, fmt.Sprintf("-P"))
	args = append(args, strconv.FormatUint(uint64(portNum), 10))
	args = append(args, fmt.Sprintf("%s@%s:%s", username, addr, remotePath))
	args = append(args, localPath)
	cmd = exec.Command("scp", args...)
	fmt.Printf("%v", cmd.Args)
	cmd.Stderr = &stdout
	err = cmd.Run()
	if err != nil {
		return errors.New("scp to " + addr + "failed: " + err.Error())
	}
	return nil
}
