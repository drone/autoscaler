// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package sshutil

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/crypto/ssh"
)

// Fingerprint returns the md5 fingerprint.
func Fingerprint(signer ssh.Signer) (out string) {
	hash := md5.Sum(
		signer.PublicKey().Marshal(),
	)
	for i := 0; i < 16; i++ {
		if i > 0 {
			out += ":"
		}
		out += fmt.Sprintf("%02x", hash[i])
	}
	return out
}

// ParsePrivateKey parses the private key and returns an ssh.Signer.
func ParsePrivateKey(s string) (ssh.Signer, error) {
	if strings.Contains(s, "-----BEGIN RSA PRIVATE KEY-----") {
		return parsePrivateKeyString(s)
	} else if strings.HasSuffix(s, "=") {
		return parsePrivateKeyBase64(s)
	}
	return parsePrivateKeyFile(s)
}

// helper function parses the private key string.
func parsePrivateKeyString(s string) (ssh.Signer, error) {
	return ssh.ParsePrivateKey([]byte(s))
}

// helper function parses the base64 private key string.
func parsePrivateKeyBase64(s string) (ssh.Signer, error) {
	out, _ := base64.StdEncoding.DecodeString(s)
	return ssh.ParsePrivateKey(out)
}

// helper function parses the private key file.
func parsePrivateKeyFile(file string) (ssh.Signer, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return ssh.ParsePrivateKey(data)
}
