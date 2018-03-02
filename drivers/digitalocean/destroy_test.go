// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package digitalocean

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/config"
	"github.com/drone/autoscaler/mocks"

	"github.com/golang/mock/gomock"
	"github.com/h2non/gock"
	"golang.org/x/crypto/ssh"
)

func TestDestroy(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	defer gock.Off()

	gock.New("https://api.digitalocean.com").
		Delete("/v2/droplets/3164494").
		Reply(204)

	mockContext := context.TODO()
	mockSigner, _ := ssh.ParsePrivateKey(testkey)
	mockConfig := config.Config{}
	mockInstance := &autoscaler.Instance{
		ID: "3164494",
	}

	// base provider to mock SSH calls.
	mockProvider := mocks.NewMockProvider(controller)
	mockProvider.EXPECT().Execute(mockContext, mockInstance, gomock.Any()).Return(nil, nil)

	p := Provider{
		Provider: mockProvider,
		config:   mockConfig,
		signer:   mockSigner,
	}

	err := p.Destroy(mockContext, mockInstance)
	if err != nil {
		t.Error(err)
	}
}

func TestDestroyShutdownError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	defer gock.Off()

	gock.New("https://api.digitalocean.com").
		Delete("/v2/droplets/3164494").
		Reply(204)

	mockError := errors.New("oh no")
	mockContext := context.TODO()
	mockSigner, _ := ssh.ParsePrivateKey(testkey)
	mockConfig := config.Config{}
	mockInstance := &autoscaler.Instance{
		ID: "3164494",
	}

	// base provider to mock SSH calls.
	mockProvider := mocks.NewMockProvider(controller)
	mockProvider.EXPECT().Execute(mockContext, mockInstance, gomock.Any()).Return(nil, mockError)

	p := Provider{
		Provider: mockProvider,
		config:   mockConfig,
		signer:   mockSigner,
	}

	err := p.Destroy(mockContext, mockInstance)
	if err != nil {
		t.Error(err)
	}
}

func TestDestroyDeleteError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	defer gock.Off()

	gock.New("https://api.digitalocean.com").
		Delete("/v2/droplets/3164494").
		Reply(500)

	mockContext := context.TODO()
	mockSigner, _ := ssh.ParsePrivateKey(testkey)
	mockConfig := config.Config{}
	mockInstance := &autoscaler.Instance{
		ID: "3164494",
	}

	// base provider to mock SSH calls.
	mockProvider := mocks.NewMockProvider(controller)
	mockProvider.EXPECT().Execute(mockContext, mockInstance, gomock.Any()).Return(nil, nil)

	p := Provider{
		Provider: mockProvider,
		config:   mockConfig,
		signer:   mockSigner,
	}

	err := p.Destroy(mockContext, mockInstance)
	if err == nil {
		t.Errorf("Expect error returned from digital ocean")
	} else if _, ok := err.(*godo.ErrorResponse); !ok {
		t.Errorf("Expect ErrorResponse digital ocean")
	}
}

func TestDestroyInvalidInput(t *testing.T) {
	i := &autoscaler.Instance{}
	p := Provider{}
	err := p.Destroy(context.TODO(), i)
	if _, ok := err.(*strconv.NumError); !ok {
		t.Errorf("Expected invalid or missing ID error")
	}
}

var testkey = []byte(`
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
`)
