package yandexcloud

// Option configures a Digital Ocean provider option.
type Option func(*provider)

func WithToken(token string) Option {
	return func(p *provider) {
		p.token = token
	}
}

func WithSubnetID(subnetID string) Option {
	return func(p *provider) {
		p.subnetID = subnetID
	}
}

func WithFolderID(folderID string) Option {
	return func(p *provider) {
		p.folderID = folderID
	}
}
