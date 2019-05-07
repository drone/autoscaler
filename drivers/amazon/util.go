// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
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
	"ap-south-1":     "ami-0189d76e",
	"us-east-1":      "ami-43a15f3e",
	"ap-northeast-1": "ami-0d74386b",
	"eu-west-1":      "ami-f90a4880",
	"ap-southeast-1": "ami-52d4802e",
	"ca-central-1":   "ami-ae55d2ca",
	"us-west-1":      "ami-925144f2",
	"eu-central-1":   "ami-7c412f13",
	"sa-east-1":      "ami-423d772e",
	"cn-north-1":     "ami-cc4499a1",
	"cn-northwest-1": "ami-fd0e1a9f",
	"us-gov-west-1":  "ami-893fb4e8",
	"ap-southeast-2": "ami-d38a4ab1",
	"eu-west-2":      "ami-f4f21593",
	"ap-northeast-2": "ami-a414b9ca",
	"us-west-2":      "ami-4e79ed36",
	"us-east-2":      "ami-916f59f4",
	"eu-west-3":      "ami-0e55e373",
}
