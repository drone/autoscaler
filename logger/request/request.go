// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package request

import (
	"net/http"
	"time"

	"github.com/drone/autoscaler/logger"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

// Logger provides logrus middleware.
func Logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		fields := logrus.Fields{
			"method":     r.Method,
			"request":    r.RequestURI,
			"remote":     r.RemoteAddr,
			"referer":    r.Referer(),
			"user-agent": r.UserAgent(),
		}
		log := logrus.WithFields(fields)
		ctx := r.Context()
		ctx = logger.WithContext(ctx, logger.Logrus(log))
		next.ServeHTTP(rw, r)

		fields["status"] = rw.Status()
		fields["duration"] = time.Since(start)
		if id := r.Context().Value(middleware.RequestIDKey); id != nil {
			fields["request-id"] = id
		}
		log.WithFields(fields).Debugln("request completed")
	}
	return http.HandlerFunc(fn)
}
