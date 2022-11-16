package yandexcloud

import "github.com/drone/autoscaler/drivers/internal/userdata"

var userdataT = userdata.Parse(`#cloud-config

write_files:
  - path: /etc/docker/ca.pem
    permissions: '0444'
    encoding: b64
    content: {{ .CACert | base64 }}
  - path: /etc/docker/cert.pem
    permissions: '0444'
    encoding: b64
    content: {{ .TLSCert | base64 }}
  - path: /etc/docker/key.pem
    permissions: '0400'
    encoding: b64
    content: {{ .TLSKey | base64 }}

runcmd:
  - sudo sed -i "s/ExecStart=\\/usr\\/bin\\/dockerd/ExecStart=\\/usr\\/bin\\/dockerd -H tcp:\\/\\/0.0.0.0:2376 --tlsverify --tlscacert=\\/etc\\/docker\\/ca.pem --tlskey=\\/etc\\/docker\\/key.pem --tlscert=\\/etc\\/docker\\/cert.pem/g" /lib/systemd/system/docker.service
  - sudo sed -i "s/DOCKER_OPTS=.*/DOCKER_OPTS=\"-H tcp:\\/\\/0.0.0.0:2376\"/g" /etc/init.d/docker
  - sudo systemctl daemon-reload
  - sudo systemctl restart docker.service
`)
