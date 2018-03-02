// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package scripts

import (
	"bytes"
	"text/template"
)

// GenerateTeardown generates a teardown script.
func GenerateTeardown() (string, error) {
	buf := new(bytes.Buffer)
	err := teardownT.Execute(buf, nil)
	return buf.String(), err
}

var teardownT = template.Must(template.New("_").Funcs(funcs).Parse(`
set -x;

sudo docker ps
sudo docker stop -t 3600 agent
sudo docker ps -a
`))
