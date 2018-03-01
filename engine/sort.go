// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package engine

import "github.com/drone/autoscaler"

// byCreated sorts the server list by created date.
type byCreated []*autoscaler.Server

func (a byCreated) Len() int           { return len(a) }
func (a byCreated) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byCreated) Less(i, j int) bool { return a[i].Created < a[j].Created }
