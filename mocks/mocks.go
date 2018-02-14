// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package mocks

//go:generate mockgen -package=mocks -destination=mocks_gen.go github.com/drone/autoscaler ServerStore,Provider
//go:generate mockgen -package=mocks -destination=mocks_gen_drone.go github.com/drone/drone-go/drone Client
