// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package autoscaler

import (
	"context"
)

// A Scaler implements an algorithm to automatically scale up
// or scale down the available pool of servers.
type Scaler interface {
	Scale(context.Context) error
}

// // Stats stores details about the synchronization process.
// type Stats struct {
// 	Start    int
// 	Finish   int
// 	Servers  int
// 	Capacity int
// 	Running  int
// 	Pending  int
// }
