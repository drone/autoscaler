// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package base

import (
	"context"
	"net"
	"time"

	"github.com/drone/autoscaler"
	"golang.org/x/crypto/ssh"
)

// Provider returns a base Provider.
func Provider(username, port string, signer ssh.Signer) autoscaler.Provider {
	return &provider{
		username: username,
		port:     port,
		timeout:  time.Minute,
		signer:   signer,
	}
}

type provider struct {
	username string
	port     string
	timeout  time.Duration
	signer   ssh.Signer
}

func (p *provider) Create(ctx context.Context, opts autoscaler.InstanceCreateOpts) (*autoscaler.Instance, error) {
	panic("not implemented")
}

func (p *provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	panic("not implemented")
}

func (p *provider) Execute(ctx context.Context, instance *autoscaler.Instance, command string) ([]byte, error) {
	return p.execute(ctx, instance.Address, command)
}

func (p *provider) Ping(ctx context.Context, instance *autoscaler.Instance) error {
	// Hosting providers may block ping requests and some ping libraries
	// require special linux capabilities. Using SSH to verify connectivity.
	out, err := p.execute(ctx, instance.Address, "whoami")
	if err != nil {
		err = &autoscaler.InstanceError{
			Err:  err,
			Logs: out,
		}
	}
	return err
}

func (p *provider) execute(ctx context.Context, address, command string) ([]byte, error) {
	config := &ssh.ClientConfig{
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         p.timeout,
		User:            p.username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(p.signer),
		},
	}

	addr := net.JoinHostPort(address, p.port)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	sess, err := conn.NewSession()
	if err != nil {
		return nil, err
	}
	return sess.CombinedOutput(command)
}
