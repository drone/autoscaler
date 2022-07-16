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
	"af-south-1":     "ami-05f63785d68762258", // Release: 20220110, Region Name: Africa (Cape Town)
	"ap-east-1":      "ami-0504f786c867c393b", // Release: 20220110, Region Name: Asia Pacific (Hong Kong)
	"ap-northeast-1": "ami-09c48a5d777342713", // Release: 20220110, Region Name: Asia Pacific (Tokyo)
	"ap-northeast-2": "ami-02f4931c64ffab121", // Release: 20220110, Region Name: Asia Pacific (Seoul)
	"ap-northeast-3": "ami-0268e471de7bae4b4", // Release: 20220110, Region Name: Asia Pacific (Osaka)
	"ap-south-1":     "ami-0a5b602444b05877e", // Release: 20220110, Region Name: Asia Pacific (Mumbai)
	"ap-southeast-1": "ami-044c31103a8df7b51", // Release: 20220110, Region Name: Asia Pacific (Singapore)
	"ap-southeast-2": "ami-00fb768e8271d7560", // Release: 20220110, Region Name: Asia Pacific (Sydney)
	"ap-southeast-3": "ami-0a2d445c9f3e6f7fc", // Release: 20220110, Region Name: Asia Pacific (Jakarta)
	"ca-central-1":   "ami-0ada0b5f702f77242", // Release: 20220110, Region Name: Canada (Central)
	"eu-central-1":   "ami-0d267e97f16681cd8", // Release: 20220110, Region Name: Europe (Frankfurt)
	"eu-north-1":     "ami-056bbd85327482a72", // Release: 20220110, Region Name: Europe (Stockholm)
	"eu-south-1":     "ami-04fc3b43df99bed69", // Release: 20220110, Region Name: Europe (Milan)
	"eu-west-1":      "ami-08307ebe62cde256a", // Release: 20220110, Region Name: Europe (Ireland)
	"eu-west-2":      "ami-065536dd9c2967f5c", // Release: 20220110, Region Name: Europe (London)
	"eu-west-3":      "ami-06bc48ca7e148506d", // Release: 20220110, Region Name: Europe (Paris)
	"me-south-1":     "ami-05a3d520f37b54e25", // Release: 20220110, Region Name: Middle East (Bahrain)
	"sa-east-1":      "ami-05572d2849beb8330", // Release: 20220110, Region Name: South America (São Paulo)
	"us-east-1":      "ami-0d4c664d2c7345cf1", // Release: 20220110, Region Name: US East (N. Virginia)
	"us-east-2":      "ami-08be70d36872187b9", // Release: 20220110, Region Name: US East (Ohio)
	"us-west-1":      "ami-0889468400f78c6cd", // Release: 20220110, Region Name: US West (N. California)
	"us-west-2":      "ami-02da538d84c7792a9", // Release: 20220110, Region Name: US West (Oregon)

	// AWS China: Ubuntu Server 20.04 LTS
	"cn-north-1":     "ami-0741e7b8b4fb0001c", // Release: 20210720, Region Name: China (Beijing)
	"cn-northwest-1": "ami-0883e8062ff31f727", // Release: 20210720, Region Name: China (Ningxia)

	// AWS GovCloud (US): Ubuntu Server 20.04 LTS
	"us-gov-east-1": "ami-03958388786db1513", // Release: 20211129, Region Name: AWS GovCloud (US-East)
	"us-gov-west-1": "ami-066189aeb91baa0ab", // Release: 20211129, Region Name: AWS GovCloud (US-West)
}
