package yandexcloud

import (
	"bytes"
	"testing"

	"github.com/drone/autoscaler"
)

func TestTemplate(t *testing.T) {
	type args struct {
		data extendedOpts
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				extendedOpts{
					autoscaler.InstanceCreateOpts{
						Name:    "test",
						CAKey:   []byte("cakey"),
						TLSKey:  []byte("tls"),
						TLSCert: []byte(`-----BEGIN CERTIFICATE-----)`),
					},
					"ubuntu",
					[]string{"ssh-rsa auth-key"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			if err := Template(buf, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Template() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
