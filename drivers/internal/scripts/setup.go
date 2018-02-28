// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package scripts

import (
	"bytes"
	"text/template"
)

// SetupOpts provides options for generating the setup
// shell script.
type SetupOpts struct {
	Instance struct {
		Name string
		Addr string
	}
	Server struct {
		Host   string
		Secret string
	}
	Agent struct {
		Image    string
		Capacity int
	}
	Cadvisor struct {
		Secret  string
		Disable bool
	}
}

// GenerateSetup generates the agent setup script.
func GenerateSetup(opts SetupOpts) (string, error) {
	buf := new(bytes.Buffer)
	err := setupT.Execute(buf, &opts)
	return buf.String(), err
}

var setupT = template.Must(template.New("_").Funcs(funcs).Parse(`
set -x;
set -e;

if ! [ -x "$(command -v docker)" ]; then
  curl -fsSL get.docker.com -o get-docker.sh
  sh get-docker.sh
  sudo usermod -aG docker $(whoami)
fi

echo -n 'admin:{SHA}{{sha .Cadvisor.Secret}}' > $HOME/.htpasswd;

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
-e DRONE_SECRET={{.Server.Secret}} \
-e DRONE_SERVER={{.Server.Host}} \
-e DRONE_MAX_PROCS={{.Agent.Capacity}} \
-e DRONE_HOSTNAME={{.Instance.Name}} \
--name=agent \
{{.Agent.Image}};

sudo docker ps;
`))
