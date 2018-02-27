// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package scripts

import (
	"testing"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"

	"github.com/pmezard/go-difflib/difflib"
)

func TestGenerateInstall(t *testing.T) {
	conf := config.Config{}
	conf.Agent.Host = "localhost:9000"
	conf.Agent.Image = "drone/agent:0.8"
	conf.Agent.Token = "a8842634682b789"

	server := autoscaler.Server{}
	server.Name = "server1"
	server.Secret = "14bb43312eada8a"
	server.Capacity = 2

	script, err := GenerateInstall(conf, &server)
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := script, installScript; got != want {
		t.Errorf("Invalid script")
		diff := difflib.UnifiedDiff{
			A:        difflib.SplitLines(got),
			B:        difflib.SplitLines(want),
			FromFile: "Got",
			ToFile:   "Want",
			Context:  5,
		}
		text, _ := difflib.GetUnifiedDiffString(diff)
		t.Log(text)
	}
}

func TestGenerateTeardown(t *testing.T) {
	conf := config.Config{}

	script, err := GenerateTeardown(conf)
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := script, teardownScript; got != want {
		t.Errorf("Invalid script")
		diff := difflib.UnifiedDiff{
			A:        difflib.SplitLines(got),
			B:        difflib.SplitLines(want),
			FromFile: "Got",
			ToFile:   "Want",
			Context:  5,
		}
		text, _ := difflib.GetUnifiedDiffString(diff)
		t.Log(text)
	}
}

var teardownScript = `
set -x;

sudo docker ps
sudo docker stop -t 3600 agent
sudo docker ps -a
`

var installScript = `
set -x;
set -e;

if ! [ -x "$(command -v docker)" ]; then
  curl -fsSL get.docker.com -o get-docker.sh
  sh get-docker.sh
  sudo usermod -aG docker $(whoami)
fi

echo -n 'admin:{SHA}0RIWnjGvcw2wHMferV9MJVSo0Uw=' > $HOME/.htpasswd;

sudo docker run \
--volume=/:/rootfs:ro \
--volume=/var/run:/var/run:rw \
--volume=/sys:/sys:ro \
--volume=/var/lib/docker/:/var/lib/docker:ro \
--volume=/dev/disk/:/dev/disk:ro \
--volume=$HOME/.htpasswd:/root/.htpasswd \
--publish=8080:8080 \
--detach=true \
--name=cadvisor \
google/cadvisor:latest \
--http_auth_realm localhost \
--http_auth_file /root/.htpasswd;

sudo docker run \
--detach=true \
--restart=always \
--volume /var/run/docker.sock:/var/run/docker.sock \
-e DRONE_SECRET=a8842634682b789 \
-e DRONE_SERVER=localhost:9000 \
-e DRONE_MAX_PROCS=2 \
-e DRONE_HOSTNAME=server1 \
--name=agent \
drone/agent:0.8;

sudo docker ps;
`
