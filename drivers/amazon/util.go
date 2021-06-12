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
	// AWS Regions: Ubuntu Server 20.04 LTS
	"af-south-1":     "ami-0b24f8168513355d8", // Release: 20210610
	"ap-east-1":      "ami-02df46111756ad5da", // Release: 20210610
	"ap-northeast-1": "ami-0339d948b9577fc0b", // Release: 20210610
	"ap-northeast-2": "ami-0c06610b37cf86e57", // Release: 20210610
	"ap-northeast-3": "ami-06b93c6f10372b44d", // Release: 20210610
	"ap-south-1":     "ami-0a180a50150c94d58", // Release: 20210610
	"ap-southeast-1": "ami-000892e0bc49a951d", // Release: 20210610
	"ap-southeast-2": "ami-093be32c43c9f6bba", // Release: 20210610
	"ca-central-1":   "ami-034e0f5c73669ebe4", // Release: 20210610
	"eu-central-1":   "ami-0eb471e022a0d8fc6", // Release: 20210610
	"eu-north-1":     "ami-0dd744cb0659e0cc2", // Release: 20210610
	"eu-south-1":     "ami-0c67192cfd77d2216", // Release: 20210610
	"eu-west-1":      "ami-04d9ce6b8cf0f422c", // Release: 20210610
	"eu-west-2":      "ami-086a05edfed3d225e", // Release: 20210610
	"eu-west-3":      "ami-01fea7461e7f2adba", // Release: 20210610
	"me-south-1":     "ami-06abb23e5f7fdea42", // Release: 20210610
	"sa-east-1":      "ami-07de2a52834192a7b", // Release: 20210610
	"us-east-1":      "ami-0c536cd6abac1a385", // Release: 20210610
	"us-east-2":      "ami-06382629a9eb569e3", // Release: 20210610
	"us-west-1":      "ami-0a1a90c77c33d81f9", // Release: 20210610
	"us-west-2":      "ami-0a1b477074e2f1708", // Release: 20210610

	// AWS China: Ubuntu Server 20.04 LTS
	"cn-north-1":     "ami-0b6277ff0310832fb", // Release: 20210518
	"cn-northwest-1": "ami-06ada3264ec0d22b4", // Release: 20210518

	// AWS GovCloud (US): Ubuntu Server 20.04 LTS
	"us-gov-east-1": "ami-0cb7240ebe7a02450", // Release: 20210503
	"us-gov-west-1": "ami-0777d7afb149b4542", // Release: 20210503
}
