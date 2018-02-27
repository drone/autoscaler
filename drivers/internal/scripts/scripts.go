// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package scripts

import (
	"bytes"
	"text/template"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"
)

// GenerateInstall generates an installation script.
func GenerateInstall(config config.Config, server *autoscaler.Server) (string, error) {
	buf := new(bytes.Buffer)
	err := installT.Execute(buf, map[string]interface{}{
		"Server": server,
		"Config": config,
	})
	return buf.String(), err
}

// GenerateTeardown generates a teardown script.
func GenerateTeardown(config config.Config) (string, error) {
	buf := new(bytes.Buffer)
	err := teardownT.Execute(buf, map[string]interface{}{
		"Config": config,
	})
	return buf.String(), err
}

var teardownT = template.Must(template.New("_").Funcs(funcs).Parse(`
set -x;

sudo docker ps
sudo docker stop -t 3600 agent
sudo docker ps -a
`))

var installT = template.Must(template.New("_").Funcs(funcs).Parse(`
set -x;
set -e;

if ! [ -x "$(command -v docker)" ]; then
  curl -fsSL get.docker.com -o get-docker.sh
  sh get-docker.sh
  sudo usermod -aG docker $(whoami)
fi

echo -n 'admin:{SHA}{{sha .Server.Secret}}' > $HOME/.htpasswd;

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
-e DRONE_SECRET={{.Config.Agent.Token}} \
-e DRONE_SERVER={{.Config.Agent.Host}} \
-e DRONE_MAX_PROCS={{.Server.Capacity}} \
-e DRONE_HOSTNAME={{.Server.Name}} \
--name=agent \
{{.Config.Agent.Image}};

sudo docker ps;
`))
