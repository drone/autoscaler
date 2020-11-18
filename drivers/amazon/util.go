// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package amazon

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// helper function converts an array of tags in string
// format to an array of ec2 tags.
func convertTags(in map[string]string) []*ec2.Tag {
	var out []*ec2.Tag
	for k, v := range in {
		out = append(out, &ec2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return out
}

// helper function creates a copy of map[string]string
func createCopy(in map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range in {
		out[k] = v
	}
	return out
}

// helper function returns the default image based on the
// selected region.
func defaultImage(region string) string {
	return images[region]
}

var images = map[string]string{
	"af-south-1":     "ami-0196a23f828d6e619", // Ubuntu Server 20.04 LTS (Release: 20201112.1)
	"ap-east-1":      "ami-f6511c87",          // Ubuntu Server 20.04 LTS (Release: 20201112.1)
	"ap-northeast-1": "ami-0e40a27db137d33cb", // Ubuntu Server 20.04 LTS (Release: 20201112.1)
	"ap-northeast-2": "ami-01ff1255cee8004b8", // Ubuntu Server 20.04 LTS (Release: 20201112.1)
	"ap-northeast-3": "ami-0b58a665b8d0d720c", // Ubuntu Server 20.04 LTS (Release: 20201112.1)
	"ap-south-1":     "ami-0cecfffd8cae9481c", // Ubuntu Server 20.04 LTS (Release: 20201112.1)
	"ap-southeast-1": "ami-0dbb8181cd0ce9cff", // Ubuntu Server 20.04 LTS (Release: 20201112.1)
	"ap-southeast-2": "ami-0f150e4544fb95045", // Ubuntu Server 20.04 LTS (Release: 20201112.1)
	"ca-central-1":   "ami-03060448f5c8f2199", // Ubuntu Server 20.04 LTS (Release: 20201112.1)
	"cn-north-1":     "ami-00cc446c1f9e0b72a", // Ubuntu Server 20.04 LTS (Release: 20201014)
	"cn-northwest-1": "ami-09f14afb2e15caab5", // Ubuntu Server 20.04 LTS (Release: 20201014)
	"eu-central-1":   "ami-09f14afb2e15caab5", // Ubuntu Server 20.04 LTS (Release: 20201112.1)
	"eu-north-1":     "ami-01450210d4ebb3bab", // Ubuntu Server 20.04 LTS (Release: 20201112.1)
	"eu-south-1":     "ami-0e3c0649c89ccddc9", // Ubuntu Server 20.04 LTS (Release: 20201112.1)
	"eu-west-1":      "ami-048309a44dad514df", // Ubuntu Server 20.04 LTS (Release: 20201112.1)
	"eu-west-2":      "ami-099ae17a6a688b1cc", // Ubuntu Server 20.04 LTS (Release: 20201112.1)
	"eu-west-3":      "ami-098efdd0afb686fd5", // Ubuntu Server 20.04 LTS (Release: 20201112.1)
	"me-south-1":     "ami-098b94183f8e74ecc", // Ubuntu Server 20.04 LTS (Release: 20201112.1)
	"sa-east-1":      "ami-0cc03bf224d6eb2fc", // Ubuntu Server 20.04 LTS (Release: 20201112.1)
	"us-east-1":      "ami-08306577a6694f5e7", // Ubuntu Server 20.04 LTS (Release: 20201112.1)
	"us-east-2":      "ami-0be9fcdb56a1f1226", // Ubuntu Server 20.04 LTS (Release: 20201112.1)
	"us-gov-east-1":  "ami-2d658d5c",          // Ubuntu Server 20.04 LTS (Release: 20201026)
	"us-gov-west-1":  "ami-f482ba95",          // Ubuntu Server 20.04 LTS (Release: 20201026)
	"us-west-1":      "ami-04d12df4da18327bd", // Ubuntu Server 20.04 LTS (Release: 20201112.1)
	"us-west-2":      "ami-082e4f383a98efbe9", // Ubuntu Server 20.04 LTS (Release: 20201112.1)
}
