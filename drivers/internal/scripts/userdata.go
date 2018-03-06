// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package scripts

import (
	"encoding/base64"
	"io/ioutil"
	"os"

	"text/template"
)

var (
	UserdataFuncmap = map[string]interface{}{
		"base64": func(src []byte) string {
			return base64.StdEncoding.EncodeToString(src)
		},
	}
)

func UserdataPrepare(val string) *template.Template {
	tmpl := template.New("_").Funcs(UserdataFuncmap)

	if _, err := os.Stat(val); err == nil {
		content, err := ioutil.ReadFile(val)

		if err != nil {
			return nil
		}

		return template.Must(tmpl.Parse(string(content)))
	}

	decoded, err := base64.StdEncoding.DecodeString(val)

	if err != nil {
		return template.Must(tmpl.Parse(val))
	}

	return template.Must(tmpl.Parse(string(decoded)))
}
