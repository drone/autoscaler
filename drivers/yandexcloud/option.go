package yandexcloud

// Option configures a Digital Ocean provider option.
type Option func(*provider)

// WithToken returns an option to set the token.
func WithToken(token string) Option {
	return func(p *provider) {
		p.token = token
	}
}

// WithServiceAccountJSON returns an option to set the token.
func WithServiceAccountJSON(serviceAccount string) Option {
	return func(p *provider) {
		p.serviceAccountJSON = serviceAccount
	}
}

// WithFolderID returns an option to set the folder id.
func WithFolderID(folderID string) Option {
	return func(p *provider) {
		p.folderID = folderID
	}
}

// WithSubnetID returns an option to set the subnet id.
func WithSubnetID(subnetID string) Option {
	return func(p *provider) {
		p.subnetID = subnetID
	}
}

// WithZone returns an option to set the zone.
func WithZone(zone []string) Option {
	return func(p *provider) {
		p.zone = zone
	}
}

// WithDiskSize returns an option to set the disk
// size in bytes.
func WithDiskSize(diskSize int64) Option {
	return func(p *provider) {
		p.diskSize = diskSize
	}
}

// WithDiskType returns an option to set the disk type.
func WithDiskType(diskType string) Option {
	return func(p *provider) {
		p.diskType = diskType
	}
}

// WithResourceCores returns an option to set the resource cores.
func WithResourceCores(resourceCores int64) Option {
	return func(p *provider) {
		p.resourceCores = resourceCores
	}
}

// WithResourceCoreFraction returns an option to set the resource core fraction.
func WithResourceCoreFraction(resourceCoreFraction int64) Option {
	return func(p *provider) {
		p.resourceCoreFraction = resourceCoreFraction
	}
}

// WithResourceMemory returns an option to set the resource
// memory in bytes.
func WithResourceMemory(resourceMemory int64) Option {
	return func(p *provider) {
		p.resourceMemory = resourceMemory
	}
}

// WithPlatformID returns an option to set the platform id.
func WithPlatformID(platformID string) Option {
	return func(p *provider) {
		p.platformID = platformID
	}
}

// WithImageFolderID returns an option to set the image folder id.
func WithImageFolderID(imageFolderID string) Option {
	return func(p *provider) {
		p.imageFolderID = imageFolderID
	}
}

// WithImageFamily returns an option to set the image family.
func WithImageFamily(imageFamily string) Option {
	return func(p *provider) {
		p.imageFamily = imageFamily
	}
}

// WithPreemptible returns an options to set the preemptible.
func WithPreemptible(preemptible bool) Option {
	return func(p *provider) {
		p.preemptible = preemptible
	}
}

// WithPrivateIP returns an options to set the privateIP.
func WithPrivateIP(privateIP bool) Option {
	return func(p *provider) {
		p.privateIP = privateIP
	}
}

// WithSecurityGroups returns an option to set security groups.
func WithSecurityGroups(groupIDs []string) Option {
	return func(p *provider) {
		p.securityGroupIDs = groupIDs
	}
}

// WithSSHUserKeyPair returns an option to set ssh user key pair.
func WithSSHUserKeyPair(pair string) Option {
	return func(p *provider) {
		p.sshUserPublicKeyPair = pair
	}
}

// WithDockerComposeConfig returns an option to set docker-compose config.
func WithDockerComposeConfig(conf string) Option {
	return func(p *provider) {
		p.dockerComposeMetadata = conf
	}
}
