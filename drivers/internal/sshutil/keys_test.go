// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package sshutil

import (
	"encoding/base64"
	"testing"
)

func TestFingerprint(t *testing.T) {
	signer, err := ParsePrivateKey(snakeoil)
	if err != nil {
		t.Error(err)
	}
	got, want := Fingerprint(signer), "58:8e:30:66:fc:e2:ff:ad:4f:6f:02:4b:af:28:0d:c7"
	if got != want {
		t.Errorf("Want fingerprint %s, got %s", want, got)
	}
}

func TestParsePrivateKey(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString([]byte(snakeoil))
	signer, err := ParsePrivateKey(encoded)
	if err != nil {
		t.Error(err)
	}
	got, want := Fingerprint(signer), "58:8e:30:66:fc:e2:ff:ad:4f:6f:02:4b:af:28:0d:c7"
	if got != want {
		t.Errorf("Want fingerprint %s, got %s", want, got)
	}
}

const snakeoil = `
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAwz+uIrhrf/C+Ku0EaofJQvPrmkQpeV2Bx3zVWjJ7GGn1k8GQ
bGXhNEsvzV3XAoFEkSD3uc6W5fU0/HxdbLtxFxUy2dg8tO/AmTrVsHAtZBMjS6ha
2ada/Xg/iH2ZvWj2t5XCAcOpotl5SBafZYytfWal2IEQ/qBEm8olZqyo3GGPxCBn
vUspbJt/i4Nzp/0mQmUrwv5EvaJa8G0ILt/AmYrKsFY2VCTt8v4XRJaPYejSyTDt
DddH6HWxIzc5fn0WMD+nCLb+a7Q33RPfZ7Y2CmuG088INEwYNsURYVLGzHwNkT7a
lHmDVH2+0YgU08m9bnEFTkcL1aH/HVIN/tXe7wIDAQABAoIBAF1FL19gr+HHVGDX
JrPpN8inExZ3l0Rl2dg9FwJmeQ05mNnDrsVJieJcRHKbcFm+/M1DbXOyb71cfLpc
gpitliGLu+X6+U0J9vx78Za+j8Btr/+1ZejxnHLXHaqLLYUg/jLG9I25NXEY6Gn6
fJybLkloXrNlPIQWdY/iail5M5VKo/CtOuyoSNJzN2HShe+uU4CR6js4URK7QiAA
y6dUW4VtkYlZCOATqHIMUAIx5fi/734v5b8ABUHOFNpaLBbUgKyvPS0H4MatjAJf
n4MiEj0A+Fyvv1UiXQJV5uQ5/Z3mv4Zf9dxZ6qHbGEzdGZWB19AClYTwKYfu/odu
IK5ubjECgYEA6uzRiDGxnB/Xfbb19O4WDgGut/qwP9qm+2/z08HFRv8n1VGSXsWY
AYl0VNAYbovOCGHGzhubYWbw0RenYaL/9YPiWa+ealMbDf6YsmU+0pvXuzyFtYI9
RHIlP81ViiDXcLzu/H7BVvEv4DfHSD9jFkWMlc83TCjQX/kEuSnulUcCgYEA1MOs
bB5V/Sw4dVj/a46zdo9ZK3JkTrNnqdI8nYqlqqSO4L11i2tJWhoh5ueDJ+NLiEIu
Wujdba7I3LzyetGITTREsLShPmM0EIuX0jJpeeTO3ylaUJppiummnsGmqMfVqdik
UYrwlD0ekm5vko7tUZqmEIgVXNA2kGDyNpR93RkCgYAUP997vtTRYUlA09F1kEQk
Zu65ewlQJ7e2+ppoyU4I5Zt4XrSgKKYGk+OMH/fLJ4/V1x+8ylJlXesqCsDpwJQR
hJGxK1sbTRiK50QgNGvq2XYJ9JiN4bEIQlKFolxaMKSBWje7We2uYdG/oO8zggs3
cz0/+IGKtgXoD93hXATtpwKBgFyPb9Btdh05AqrSd/Pz1dErVbCYCFlQpTV099fV
vHK7Okk9QwjPOM8Q9VS9vQo6UN7LY9061zHjSxD0xkx2IWTs60Ewo8E/aSQVhov0
UHyt9O2S0O6l7mp3cXw5ZOaiYSqNzBaJalYjLMypbLKGqWnJ7JreiOSi1EoFUvo5
qXPpAoGAVZCcIes93+9MHbcbhVyHv3X0XulrPQMRcve0M4MarZ7UWAIPrcIKhQC9
2M+S1L0BJaHZxTQyxCb2XNRKqzQnmuhrT/A+tH8SUpl2ZLlF1KiNn4UsznBrvsaJ
OoKlArIydogC0Ugu2LoKMw8oIkpf2ANTii3dYJN6ulBKE6EYaRk=
-----END RSA PRIVATE KEY-----
`
