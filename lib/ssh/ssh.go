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
	"crypto"
	"crypto/dsa"
	"crypto/rsa"
	_ "crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
)

const DEFAULT_SSH_PORT = 22

// keychain implements the ssh.ClientKeyring interface
type keychain struct {
	keys []interface{}
}

func (k *keychain) Key(i int) (ssh.PublicKey, error) {
	if i < 0 || i >= len(k.keys) {
		return nil, nil
	}
	switch key := k.keys[i].(type) {
	case *rsa.PrivateKey:
		return ssh.NewRSAPublicKey(&key.PublicKey), nil
	case *dsa.PrivateKey:
		return ssh.NewDSAPublicKey(&key.PublicKey), nil
	}
	panic("unknown key type")
}

func (k *keychain) Sign(i int, rand io.Reader, data []byte) (sig []byte, err error) {
	hashFunc := crypto.SHA1
	h := hashFunc.New()
	h.Write(data)
	digest := h.Sum(nil)
	switch key := k.keys[i].(type) {
	case *rsa.PrivateKey:
		return rsa.SignPKCS1v15(rand, key, hashFunc, digest)
	}
	return nil, errors.New("ssh: unknown key type")
}

func (k *keychain) loadPEM(file string) error {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(buf)
	if block == nil {
		return errors.New("ssh: no key found")
	}
	r, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	k.keys = append(k.keys, r)
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
// The portNum is the SSH port number of the remote host.
// The command parameter is the command to run on the remote host.
func Run(
	addr string,
	username string,
	password *string,
	privKeyPath string,
	portNum uint16,
	command string) (stdout string, stderr string, err error) {

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
	if err = session.Run(command); err != nil {
		return "", "", err
	}
	return stdout_buf.String(), stderr_buf.String(), err
}

// Run a separate scp process (for now) to secure copy files between hosts.
// Only copying between the master and a single slave is supported.
// The addr parameter is the address (IP, hostname, etc) of the remote host.
// The privKeyPath parameter is the path to the private key of the master.
// The portNum is the SSH port number of the remote host.
// The incoming parameter indicates which direction to perform the copy.
// The recursive parameter indicates whether to use the -r scp option.
// Password authentication not supported for this function.
func Scp(
	addr string,
	username string,
	portNum uint16,
	incoming bool,
	recursive bool) (stdout string, stderr string, err error) {

	/* Unfortunately, there doesn't appear to be an SFTP or SCP library, so
	we'll just have to run a separate scp process.  This means no password
	authentication for when calling this method. */
	var command string = "scp "
	if recursive {
		command += "-r "
	}
	/* TODO: Finish this */
	return "", "", nil
}
