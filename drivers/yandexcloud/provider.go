package yandexcloud

import (
	"context"
	"errors"
	"fmt"

	ycsdk "github.com/yandex-cloud/go-sdk"
	"github.com/yandex-cloud/go-sdk/iamkey"

	"github.com/drone/autoscaler"
)

type provider struct {
	token              string
	serviceAccountJSON string

	folderID string
	zone     []string
	subnetID string

	platformID string
	privateIP  bool

	diskSize             int64
	diskType             string
	resourceCores        int64
	resourceCoreFraction int64
	resourceMemory       int64
	preemptible          bool

	imageFolderID string
	imageFamily   string

	service *ycsdk.SDK
}

func New(opts ...Option) (autoscaler.Provider, error) {
	var (
		key         *iamkey.Key
		credentials ycsdk.Credentials
		err         error
	)

	p := new(provider)
	for _, opt := range opts {
		opt(p)
	}

	if p.token == "" && p.serviceAccountJSON == "" {
		return nil, errors.New("token or service account must be provided")
	}
	if p.folderID == "" {
		return nil, errors.New("folderID must be provided")
	}
	if p.subnetID == "" {
		return nil, errors.New("empty subnet id")
	}
	if len(p.zone) == 0 {
		p.zone = []string{"ru-central1-a"}
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
	if p.resourceMemory == 0 {
		p.resourceMemory = 2 * 1024 * 1024 * 1024
	}
	if p.platformID == "" {
		p.platformID = "standard-v3"
	}
	if p.imageFolderID == "" {
		p.imageFolderID = "standard-images"
	}
	if p.imageFamily == "" {
		p.imageFamily = "container-optimized-image"
	}
	if p.resourceCoreFraction == 0 {
		p.resourceCoreFraction = 100
	}

	if p.token != "" {
		credentials = ycsdk.OAuthToken(p.token)
	} else {
		key, err = iamkey.ReadFromJSONBytes([]byte(p.serviceAccountJSON))
		if err != nil {
			return nil, fmt.Errorf("read service account json: %w", err)
		}

		credentials, err = ycsdk.ServiceAccountKey(key)
		if err != nil {
			return nil, fmt.Errorf("make service account credentials: %w", err)
		}
	}

	p.service, err = ycsdk.Build(context.Background(), ycsdk.Config{Credentials: credentials})
	if err != nil {
		return nil, fmt.Errorf("init yandex cloud sdk: %w", err)
	}

	return p, nil
}
