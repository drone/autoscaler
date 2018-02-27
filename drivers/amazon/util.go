// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package amazon

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// helper function converts an array of tags in string format to
// an array of ec2 tags.
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
