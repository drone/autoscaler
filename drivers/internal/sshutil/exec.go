// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package sshutil

import (
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// default timeout for establishing the connection.
const timeout = 15 * time.Minute

// Execute executes an ssh script and returns a snapshot of
// the term output and result in the form of an error.
func Execute(address, port, username, script string, signer ssh.Signer) ([]byte, error) {
	config := &ssh.ClientConfig{
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
		User:            username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	addr := net.JoinHostPort(address, port)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	sess, err := conn.NewSession()
	if err != nil {
		return nil, err
	}

	return sess.CombinedOutput(script)
}
