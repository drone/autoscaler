// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package scripts

import (
	"testing"

	"github.com/pmezard/go-difflib/difflib"
)

func TestGenerateSetup(t *testing.T) {
	opts := SetupOpts{}
	opts.Server.Host = "localhost:9000"
	opts.Server.Secret = "a8842634682b789"
	opts.Agent.Image = "drone/agent:0.8"
	opts.Agent.Capacity = 2
	opts.Instance.Addr = "1.2.3.4"
	opts.Instance.Name = "server1"
	opts.Cadvisor.Disable = false
	opts.Cadvisor.Secret = "14bb43312eada8a"

	script, err := GenerateSetup(opts)
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := script, setupScript; got != want {
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

var setupScript = `
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
