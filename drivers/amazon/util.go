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
	"af-south-1":     "ami-0852a941175b30c13", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"ap-east-1":      "ami-6d03411c",          // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"ap-northeast-1": "ami-09b86f9709b3c33d4", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"ap-northeast-2": "ami-044057cb1bc4ce527", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"ap-northeast-3": "ami-0c733c715f0e4ee50", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"ap-south-1":     "ami-0cda377a1b884a1bc", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"ap-southeast-1": "ami-093da183b859d5a4b", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"ap-southeast-2": "ami-0f158b0f26f18e619", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"ca-central-1":   "ami-0edab43b6fa892279", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"cn-north-1":     "ami-04ca1d006fc3c7320", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"cn-northwest-1": "ami-0b3b3edd594a6d6bd", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"eu-central-1":   "ami-0c960b947cbb2dd16", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"eu-north-1":     "ami-008dea09a148cea39", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"eu-south-1":     "ami-01eec6bdfa20f008e", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"eu-west-1":      "ami-06fd8a495a537da8b", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"eu-west-2":      "ami-05c424d59413a2876", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"eu-west-3":      "ami-078db6d55a16afc82", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"me-south-1":     "ami-053a63c3a3f73ca70", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"sa-east-1":      "ami-02dc8ad50da58fffd", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"us-east-1":      "ami-0dba2cb6798deb6d8", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"us-east-2":      "ami-07efac79022b86107", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"us-gov-east-1":  "ami-7a09e00b",          // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"us-gov-west-1":  "ami-581c2539",          // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"us-west-1":      "ami-021809d9177640a20", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
	"us-west-2":      "ami-06e54d05255faf8f6", // Ubuntu Server 20.04 LTS (Released Date: 20200907)
}
