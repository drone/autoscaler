// +build nolimit
// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package limiter

import "github.com/drone/autoscaler"

// Limit is an alternate implementaiton of the limiter that
// does not enforce a limit or license key.
func Limit(server autoscaler.ServerStore, token string) autoscaler.ServerStore {
	print(message)
	return server
}

var message = `
ATTENTION

Use of this software is governed by the Polyform License
that can be found in the LICENSE file.

If your use of the Licensed Work does not comply with the requirements
currently in effect as described in this License, you must purchase a
commercial license from the Licensor, its affiliated entities, or authorized
resellers, or you must refrain from using the Licensed Work.
`
