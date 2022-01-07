// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package userdata

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/drone/autoscaler"
)

func TestUserdata(t *testing.T) {
	buf := new(bytes.Buffer)
	err := T.Execute(buf, &autoscaler.InstanceCreateOpts{
		Name:    "agent-123456",
		CACert:  []byte(dummyCA),
		TLSKey:  []byte(dummykey),
		TLSCert: []byte(dummyCert),
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestUserdataFuncmap(t *testing.T) {
	buf := new(bytes.Buffer)
	err := UD.Execute(buf, &map[string]interface{}{
		"Content": "foo",
	})
	fmt.Println(buf.String())
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != UDExpected {
		t.Errorf("expected '%s', got '%s'", UDExpected, buf.String())
	}
}

var dummyCA = `-----BEGIN CERTIFICATE-----
MIIGOTCCBCGgAwIBAgIJAOE/vJd8EB24MA0GCSqGSIb3DQEBBQUAMIGyMQswCQYD
VQQGEwJGUjEPMA0GA1UECAwGQWxzYWNlMRMwEQYDVQQHDApTdHJhc2JvdXJnMRgw
FgYDVQQKDA93d3cuZnJlZWxhbi5vcmcxEDAOBgNVBAsMB2ZyZWVsYW4xLTArBgNV
BAMMJEZyZWVsYW4gU2FtcGxlIENlcnRpZmljYXRlIEF1dGhvcml0eTEiMCAGCSqG
KvbxUcDaVvXB0EU0bg==
-----END CERTIFICATE-----`

var dummykey = `-----BEGIN RSA PRIVATE KEY-----
MIIJKQIBAAKCAgEA3W29+ID6194bH6ejLrIC4hb2Ugo8v6ZC+Mrck2dNYMNPjcOK
ABvxxEtBamnSaeU/IY7FC/giN622LEtV/3oDcrua0+yWuVafyxmZyTKUb4/GUgaf
RQPf/eiX9urWurtIK7XgNGFNUjYPq4dSJQPPhwCHE/LKAykWnZBXRrX0Dq4XyApN
ku0IpjIjEXH+8ixE12wH8wt7DEvdO7T3N3CfUbaITl1qBX+Nm2Z6q4Ag/u5rl8NJ
v3TGd3xXD9yQIjmugNgxNiwAZzhJs/ZJy++fPSJ1XQxbd9qPghgGoe/ff6G7
-----END RSA PRIVATE KEY-----`

var dummyCert = `-----BEGIN CERTIFICATE-----
MIIGJzCCBA+gAwIBAgIBATANBgkqhkiG9w0BAQUFADCBsjELMAkGA1UEBhMCRlIx
d3d3LmZyZWVsYW4ub3JnMRAwDgYDVQQLDAdmcmVlbGFuMS0wKwYDVQQDDCRGcmVl
bGFuIFNhbXBsZSBDZXJ0aWZpY2F0ZSBBdXRob3JpdHkxIjAgBgkqhkiG9w0BCQEW
E2NvbnRhY3RAZnJlZWxhbi5vcmcwHhcNMTIwNDI3MTAzMTE4WhcNMjIwNDI1MTAz
DiH5uEqBXExjrj0FslxcVKdVj5glVcSmkLwZKbEU1OKwleT/iXFhvooWhQ==
-----END CERTIFICATE-----`

var UD = Parse(`#cloud-config

apt_reboot_if_required: 
package_update: false
package_upgrade: false

write_files:
  - path: /etc/systemd/system/docker.service.d/override.conf
    content: | {{nindent .Content 6 }}
`)

var UDExpected = `#cloud-config

apt_reboot_if_required: 
package_update: false
package_upgrade: false

write_files:
  - path: /etc/systemd/system/docker.service.d/override.conf
    content: | 
      foo
`
