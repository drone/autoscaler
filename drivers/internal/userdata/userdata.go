// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package userdata

import (
	"encoding/base64"
	"fmt"
	"strings"
	"text/template"
)

var funcs = map[string]interface{}{
	"base64": func(src []byte) string {
		return base64.StdEncoding.EncodeToString(src)
	},
}

// Parse parses the userdata template.
func Parse(text string) *template.Template {
	if decoded, err := base64.StdEncoding.DecodeString(text); err == nil {
		return template.Must(
			template.New("_").Funcs(funcs).Parse(string(decoded)),
		)
	}

	return template.Must(
		template.New("_").Funcs(funcs).Parse(text),
	)
}

// DetectUserdata returns the docker provisioning userdata script based on the OS string
func DetectUserdata(os string) (*template.Template, error) {
	switch {
	case strings.Contains(os, "ubuntu"):
		return Ubuntu, nil
	default:
		return nil, fmt.Errorf("docker provisioning script unsuported for the given OS image : %s", os)
	}

}

// Ubuntu is the docker installation script for the Ubuntu OS
var Ubuntu = Parse(`#cloud-config

apt_reboot_if_required: false
package_update: false
package_upgrade: false

apt:
  sources:
    docker.list:
      source: deb [arch=amd64] https://download.docker.com/linux/ubuntu $RELEASE stable
      keyid: 0EBFCD88

packages:
  - docker-ce

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
