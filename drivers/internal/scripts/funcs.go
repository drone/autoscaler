// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package scripts

import (
	"crypto/sha1"
	"encoding/base64"
)

// template func map.
var funcs = map[string]interface{}{
	"sha": generateSha,
}

// helper function generates a base64 encoded sha for use
// with .htaccess files.
func generateSha(password string) string {
	s := sha1.New()
	s.Write([]byte(password))
	hash := []byte(s.Sum(nil))
	return base64.StdEncoding.EncodeToString(hash)
}
