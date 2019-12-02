// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package server

import (
	"net/http"
	"strings"

	"github.com/drone/autoscaler/config"
	"github.com/drone/autoscaler/logger"
	"github.com/drone/drone-go/drone"

	"golang.org/x/oauth2"
)

// CheckDrone returns a middleware function that authorizes
// the incoming http.Request using the Drone API.
func CheckDrone(conf config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)

			// the user can authenticate with a global authorization
			// token provied in the Authorization header.
			token := r.Header.Get("Authorization")
			token = strings.TrimPrefix(token, "Bearer ")
			token = strings.TrimSpace(token)
			if token == "" {
				log.Debugln("missing authorization header")
				writeUnauthorized(w, errInvalidToken)
				return
			}

			// creates a new drone client using the bearer token
			// in the incoming request to authenticate with drone.
			config := new(oauth2.Config)
			auther := config.Client(
				oauth2.NoContext,
				&oauth2.Token{
					AccessToken: token,
				},
			)
			server := conf.Server.Proto + "://" + conf.Server.Host
			client := drone.NewClient(server, auther)

			// fetch the user account associated with the currently
			// authenticated bearer token. This user must exist in
			// drone and must be an administrator.
			user, err := client.Self()
			if err != nil {
				log.WithError(err).
					Errorln("cannot authenticate user")
				writeUnauthorized(w, errUnauthorized)
				return
			}

			if !user.Admin {
				log.WithError(err).
					WithField("username", user.Login).
					Errorln("insufficient privileges")
				writeForbidden(w, errForbidden)
				return
			}

			log = log.WithField("username", user.Login)
			log.Debugln("user authorized")

			next.ServeHTTP(w, r.WithContext(
				logger.WithContext(ctx, log),
			))
		})
	}
}
