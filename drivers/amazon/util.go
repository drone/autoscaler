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

// static ami id list for Ubuntu Server 20.04 LTS
// source: https://cloud-images.ubuntu.com/locator/
// filters:
// - Cloud: Amazon AWS, Amazon GovCloud, Amazon AWS China
// - Version: 20.04
// - Instance Type: hvm-ssd
var images = map[string]string{
	// AWS Regions: Ubuntu Server 20.04 LTS
	// Upstream release version: 20220706
	"af-south-1":     "ami-0f5298ccab965edeb",
	"ap-east-1":      "ami-0dfad1f1f65cd083b",
	"ap-northeast-1": "ami-0986c991cc80c6ad9",
	"ap-northeast-2": "ami-0565d651769eb3de5",
	"ap-northeast-3": "ami-0e6078093a109801c",
	"ap-south-1":     "ami-0325e3016099f9112",
	"ap-southeast-1": "ami-0eaf04122a1ae7b3b",
	"ap-southeast-2": "ami-048a2d001938101dd",
	"ap-southeast-3": "ami-09915141a4f1dafdd",
	"ca-central-1":   "ami-04a579d2f00bb4001",
	"eu-central-1":   "ami-06cac34c3836ff90b",
	"eu-north-1":     "ami-0ede84a5f28ec932a",
	"eu-south-1":     "ami-0a39f417b8836bc59",
	"eu-west-1":      "ami-0141514361b6a3c1b",
	"eu-west-2":      "ami-014b642f603e350c3",
	"eu-west-3":      "ami-0d0b8d91779dec1e5",
	"me-south-1":     "ami-0c769d841005394ee",
	"sa-east-1":      "ami-088afbba294231fe0",
	"us-east-1":      "ami-0070c5311b7677678",
	"us-east-2":      "ami-07f84a50d2dec2fa4",
	"us-west-1":      "ami-040a251ee9d7d1a9b",
	"us-west-2":      "ami-0aab355e1bfa1e72e",

	// AWS GovCloud (US): Ubuntu Server 20.04 LTS
	// Upstream release version: 20220627.1
	"us-gov-east-1": "ami-0d8ee446ec886f5cf",
	"us-gov-west-1": "ami-0cbaf57cea1d72aec",

	// AWS China: Ubuntu Server 20.04 LTS
	// Upstream release version: 20210720
	"cn-north-1":     "ami-0741e7b8b4fb0001c",
	"cn-northwest-1": "ami-0883e8062ff31f727",
}
