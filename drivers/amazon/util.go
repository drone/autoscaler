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
	"af-south-1":     "ami-08502d6674c2156ed", // Release: 20210108
	"ap-east-1":      "ami-2c5f135d",          // Release: 20210108
	"ap-northeast-1": "ami-08c925c2f4fe5aee0", // Release: 20210108
	"ap-northeast-2": "ami-0de358223fd8dc276", // Release: 20210108
	"ap-northeast-3": "ami-0598298c1d3a099c2", // Release: 20210108
	"ap-south-1":     "ami-097b711d946240d58", // Release: 20210108
	"ap-southeast-1": "ami-03425fe81e8a82dfb", // Release: 20210108
	"ap-southeast-2": "ami-0b31ea67df6098216", // Release: 20210108
	"ca-central-1":   "ami-0a20346326d3d1853", // Release: 20210108
	"eu-central-1":   "ami-0bb75d95f668ff5a7", // Release: 20210108
	"eu-north-1":     "ami-0000da4d489188a4b", // Release: 20210108
	"eu-south-1":     "ami-092549ba527b0e05a", // Release: 20210108
	"eu-west-1":      "ami-04ffbabc7935ec0e9", // Release: 20210108
	"eu-west-2":      "ami-0d738342f2d2fd5fd", // Release: 20210108
	"eu-west-3":      "ami-0b209583a4a1146dd", // Release: 20210108
	"me-south-1":     "ami-0161166adfc549dbf", // Release: 20210108
	"sa-east-1":      "ami-04ef59610885412b7", // Release: 20210108
	"us-east-1":      "ami-011899242bb902164", // Release: 20210108
	"us-east-2":      "ami-07d5003620a5450ee", // Release: 20210108
	"us-west-1":      "ami-034bf895b736be04a", // Release: 20210108
	"us-west-2":      "ami-089668cd321f3cf82", // Release: 20210108

	// AWS China: Ubuntu Server 20.04 LTS
	"cn-north-1":     "ami-0592ccadb56e65f8d", // Release: 20201112.1
	"cn-northwest-1": "ami-007d0f254ea0f8588", // Release: 20201112.1

	// AWS GovCloud (US): Ubuntu Server 20.04 LTS
	"us-gov-east-1": "ami-5ce8032d", // Release: 20201112.1
	"us-gov-west-1": "ami-da2a11bb", // Release: 20201112.1
}
