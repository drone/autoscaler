// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package server

import (
	"net/http"
	"strings"

	"github.com/drone/autoscaler/config"
	"github.com/drone/drone-go/drone"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"golang.org/x/oauth2"
)

// CheckDrone returns a middleware function that authorizes
// the incoming http.Request using the Drone API.
func CheckDrone(conf config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := hlog.FromRequest(r)

			// the user can authenticate with a global authorization
			// token provied in the Authorization header.
			token := r.Header.Get("Authorization")
			token = strings.TrimPrefix(token, "Bearer ")
			token = strings.TrimSpace(token)
			if token == "" {
				logger.Debug().
					Msg("missing authorization header")
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
				logger.Error().
					Err(err).
					Msg("cannot authenticate user")
				writeUnauthorized(w, errUnauthorized)
				return
			}

			if !user.Admin {
				logger.Error().
					Err(err).
					Str("username", user.Login).
					Msg("insufficient privileges")
				writeForbidden(w, errForbidden)
				return
			}

			// add the authorized user to the logger context.
			logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Str("username", user.Login)
			})

			logger.Debug().
				Str("username", user.Login).
				Msg("user authorized")

			next.ServeHTTP(w, r)
		})
	}
}
