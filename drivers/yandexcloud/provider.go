package yandexcloud

import (
	"context"
	"errors"
	"fmt"

	ycsdk "github.com/yandex-cloud/go-sdk"

	"github.com/drone/autoscaler"
)

type provider struct {
	token    string
	folderID string
	zones    []string
	subnetID string

	platformId string
	privateIP  bool

	diskSize            int64
	diskType            string
	resourceCores       int64
	resourceMemoryBytes int64

	imageFolderID string
	imageFamily   string

	service *ycsdk.SDK
}

func New(opts ...Option) (autoscaler.Provider, error) {
	var err error

	p := new(provider)
	for _, opt := range opts {
		opt(p)
	}

	if p.token == "" {
		return nil, errors.New("token must be provided")
	}
	if p.folderID == "" {
		return nil, errors.New("folderID must be provided")
	}
	if len(p.zones) == 0 {
		p.zones = []string{"ru-central1-a"}
	}
	if p.subnetID == "" {
		return nil, errors.New("empty subnet id")
	}
	if p.diskSize == 0 {
		p.diskSize = 10 * 1024 * 1024 * 1024
	}
	if p.diskType == "" {
		p.diskType = "network-hdd"
	}
	if p.resourceCores == 0 {
		p.resourceCores = 2
	}
	if p.resourceMemoryBytes == 0 {
		p.resourceMemoryBytes = 2 * 1024 * 1024 * 1024
	}
	if p.platformId == "" {
		p.platformId = "standard-v3"
	}
	if p.imageFolderID == "" {
		p.imageFolderID = "standard-images"
	}
	if p.imageFamily == "" {
		p.imageFamily = "debian-9"
	}

	p.service, err = ycsdk.Build(context.Background(), ycsdk.Config{
		Credentials: ycsdk.OAuthToken(p.token),
	})
	if err != nil {
		return nil, fmt.Errorf("init yandex cloud sdk: %w", err)
	}

	return p, nil
}
