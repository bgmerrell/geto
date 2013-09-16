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
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	_ "crypto/sha1"
)

var _privKeyPath string = ""

// keychain implements the ClientKeyring interface
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

func TestDial(privKeyPath string) (ok bool) {
	_privKeyPath = privKeyPath
	var err error

	// An SSH client is represented with a ClientConn. Currently only
	// the "password" authentication method is supported.
	//
	// To authenticate with the remote server you must pass at least one
	// implementation of ClientAuth via the Auth field in ClientConfig.
	var clientKeychain *keychain = new(keychain)
	clientKeychain.loadPEM(_privKeyPath)
	config := &ssh.ClientConfig{
		User: "bean",  /* FIXME: Do not hard-code user name */
		Auth: []ssh.ClientAuth{
			ssh.ClientAuthKeyring(clientKeychain),
		},
	}
	client, err := ssh.Dial("tcp", "10.50.26.147:22", config)  /* FIXME: Do not hard-code URI */
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil {
		panic("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())
	return true
}
