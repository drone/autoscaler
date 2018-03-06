// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package scripts

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"os"
	"testing"
)

func TestUserdataPrepare_File(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "test")

	if err != nil {
		t.Error(err)
		return
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(dummyUserdata)); err != nil {
		t.Error(err)
		return
	}

	if err := tmpfile.Close(); err != nil {
		t.Error(err)
		return
	}

	tmpl := UserdataPrepare(tmpfile.Name())
	buf := new(bytes.Buffer)

	if err := tmpl.Execute(buf, nil); err != nil {
		t.Error(err)
		return
	}

	if got, want := buf.String(), dummyUserdata; got != want {
		t.Errorf("Want parsed template of %v, got %v", want, got)
	}
}

func TestUserdataPrepare_Base64(t *testing.T) {
	input := base64.StdEncoding.EncodeToString([]byte(dummyUserdata))

	tmpl := UserdataPrepare(input)
	buf := new(bytes.Buffer)

	if err := tmpl.Execute(buf, nil); err != nil {
		t.Error(err)
		return
	}

	if got, want := buf.String(), dummyUserdata; got != want {
		t.Errorf("Want parsed template of %v, got %v", want, got)
	}

}

func TestUserdataPrepare_String(t *testing.T) {
	input := dummyUserdata

	tmpl := UserdataPrepare(input)
	buf := new(bytes.Buffer)

	if err := tmpl.Execute(buf, nil); err != nil {
		t.Error(err)
		return
	}

	if got, want := buf.String(), dummyUserdata; got != want {
		t.Errorf("Want parsed template of %v, got %v", want, got)
	}

}

var dummyUserdata = `#cloud-config
apt_reboot_if_required: false
package_update: false
package_upgrade: false
`
