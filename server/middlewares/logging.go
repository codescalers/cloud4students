// Package middlewares for middleware between api and backend
package middlewares

import (
	"net/http"
	"time"

	"github.com/codescalers/cloud4students/models"
	"github.com/rs/zerolog/log"
	"github.com/urfave/negroni/v3"
)

func AuditLogMiddleware(db models.DB) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(UserIDKey("UserID")).(string)

			lrw := negroni.NewResponseWriter(w)
			h.ServeHTTP(lrw, r)

			statusCode := lrw.Status()

			l := models.AuditLog{
				UserID:     userID,
				Method:     r.Method,
				URL:        r.URL.Path,
				Timestamp:  time.Now(),
				Success:    statusCode >= 200 && statusCode < 300,
				StatusCode: statusCode,
			}

			if err := db.CreateAuditLog(&l); err != nil {
				log.Error().Err(err).Msg("logging audit failed")
			}
		})
	}
}
