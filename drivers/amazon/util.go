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
	"af-south-1":     "ami-038cbae6ba4c280fa", // Release: 2022-01-31, Reagion Name: Africa (Cape Town)
	"ap-east-1":      "ami-064b888ccdf1708e2", // Release: 2022-01-31, Reagion Name: Asia Pacific (Hong Kong)
	"ap-northeast-1": "ami-0ec4d40472158dbd2", // Release: 2022-01-31, Reagion Name: Asia Pacific (Tokyo)
	"ap-northeast-2": "ami-077c1171fcf5593ec", // Release: 2022-01-31, Reagion Name: Asia Pacific (Seoul)
	"ap-northeast-3": "ami-035913f232768c8c8", // Release: 2022-01-31, Reagion Name: Asia Pacific (Osaka)
	"ap-south-1":     "ami-0b8959ac764ad4343", // Release: 2022-01-31, Reagion Name: Asia Pacific (Mumbai)
	"ap-southeast-1": "ami-042f884c037e74d76", // Release: 2022-01-31, Reagion Name: Asia Pacific (Singapore)
	"ap-southeast-2": "ami-0154c902f0267d0ce", // Release: 2022-01-31, Reagion Name: Asia Pacific (Sydney)
	"ap-southeast-3": "ami-0d5ff02b9a1041622", // Release: 2022-01-31, Reagion Name: Asia Pacific (Jakarta)
	"ca-central-1":   "ami-03953974e61b3bd41", // Release: 2022-01-31, Reagion Name: Canada (Central)
	"eu-central-1":   "ami-05b308c240ae70bb6", // Release: 2022-01-31, Reagion Name: Europe (Frankfurt)
	"eu-north-1":     "ami-0820d427f94ae9361", // Release: 2022-01-31, Reagion Name: Europe (Stockholm)
	"eu-south-1":     "ami-0da1c7edc984ae450", // Release: 2022-01-31, Reagion Name: Europe (Milan)
	"eu-west-1":      "ami-081ff4b9aa4e81a08", // Release: 2022-01-31, Reagion Name: Europe (Ireland)
	"eu-west-2":      "ami-0d19fa6f37a659a28", // Release: 2022-01-31, Reagion Name: Europe (London)
	"eu-west-3":      "ami-0c0f763628afa7f8b", // Release: 2022-01-31, Reagion Name: Europe (Paris)
	"me-south-1":     "ami-0538f1780d51a45d3", // Release: 2022-01-31, Reagion Name: Middle East (Bahrain)
	"sa-east-1":      "ami-0b919b06e4a6a7040", // Release: 2022-01-31, Reagion Name: South America (São Paulo)
	"us-east-1":      "ami-01b996646377b6619", // Release: 2022-01-31, Reagion Name: US East (N. Virginia)
	"us-east-2":      "ami-039af3bfc52681cd5", // Release: 2022-01-31, Reagion Name: US East (Ohio)
	"us-west-1":      "ami-08fa7c8891945eae4", // Release: 2022-01-31, Reagion Name: US West (N. California)
	"us-west-2":      "ami-0637e7dc7fcc9a2d9", // Release: 2022-01-31, Reagion Name: US West (Oregon)

	// AWS China: Ubuntu Server 20.04 LTS
	"cn-north-1":     "ami-0741e7b8b4fb0001c", // Release: 2021-07-20, Region Name: China (Beijing)
	"cn-northwest-1": "ami-0883e8062ff31f727", // Release: 2021-07-20, Region Name: China (Ningxia)

	// AWS GovCloud (US): Ubuntu Server 20.04 LTS
	"us-gov-east-1": "ami-0e1e4f0f5c274fb48", // Release: 2022-01-18, Region Name: AWS GovCloud (US-East)
	"us-gov-west-1": "ami-047d3c53ce8db1f52", // Release: 2022-01-18, Region Name: AWS GovCloud (US-West)
}
