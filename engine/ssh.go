// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package engine

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"
	"golang.org/x/crypto/ssh"
)

type sshClient struct {
	*ssh.Client
}

// clientFunc defines a builder funciton used to build and return
// the SSH client to a Server.
type sshClientFunc func(*autoscaler.Server, *config.SSH) (*sshClient, io.Closer, error)

// newSSHClient returns a new SSH client configured for the
// Server host and certificate chain.
func newSSHClient(server *autoscaler.Server, sshConfig *config.SSH) (*sshClient, io.Closer, error) {
	// Read private key for agents
	key, err := ioutil.ReadFile(sshConfig.PrivateKeyFile)
	if err != nil {
		return nil, nil, err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, nil, err
	}

	// TODO: DO NOT USE INSECUREIGNOREHOSTKEY
	config := &ssh.ClientConfig{
		User: sshConfig.User,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the remote server and perform the SSH handshake.
	_client, err := ssh.Dial("tcp", server.Address+sshConfig.Port, config)
	if err != nil {
		return nil, nil, err
	}

	_sshClient := &sshClient{_client}
	return _sshClient, _sshClient, nil
}

// Wrapper session
func (client *sshClient) Run(cmd string) (string, error) {
	// create a session and call Run then close the session
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf
	err = session.Run(cmd)
	if err != nil {
		return stderrBuf.String(), err
	}

	return stdoutBuf.String(), err
}

// Ping
func (client *sshClient) Ping() (string, error) {
	out, err := client.Run("ps -C drone-runner-exec -ocomm=")
	if err != nil {
		return out, err
	}
	return out, nil
}
