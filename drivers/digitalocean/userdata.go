// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package digitalocean

import "github.com/drone/autoscaler/drivers/internal/userdata"

var userdataT = userdata.Parse(`#cloud-config
write_files:
  - path: /etc/systemd/system/docker.service.d/override.conf
    content: |
      [Service]
      ExecStart=
      ExecStart=/usr/bin/dockerd
  - path: /etc/default/docker
    content: |
      DOCKER_OPTS=""
  - path: /etc/docker/daemon.json
    content: |
      {
        "dns": [ "8.8.8.8", "8.8.4.4" ],
        "hosts": [ "0.0.0.0:2376", "unix:///var/run/docker.sock" ],
        "tls": true,
        "tlsverify": true,
        "tlscacert": "/etc/docker/ca.pem",
        "tlscert": "/etc/docker/server-cert.pem",
        "tlskey": "/etc/docker/server-key.pem"
      }
  - path: /etc/docker/ca.pem
    encoding: b64
    content: {{ .CACert | base64 }}
  - path: /etc/docker/server-cert.pem
    encoding: b64
    content: {{ .TLSCert | base64 }}
  - path: /etc/docker/server-key.pem
    encoding: b64
    content: {{ .TLSKey | base64 }}

runcmd:
  - [ systemctl, daemon-reload ]
  - [ systemctl, restart, docker ]
`)
