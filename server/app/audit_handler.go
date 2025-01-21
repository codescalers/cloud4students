package app

import (
	"errors"
	"net/http"

	"github.com/codescalers/cloud4students/middlewares"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Example endpoint: List user's logs
// @Summary List user's logs
// @Description List user's logs
// @Tags Audit
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} []models.AuditLog
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /user/log [get]
func (a *App) ListLogsHandler(req *http.Request) (interface{}, Response) {
	userID := req.Context().Value(middlewares.UserIDKey("UserID")).(string)

	logs, err := a.db.GetUserLogs(userID)
	if err == gorm.ErrRecordNotFound || len(logs) == 0 {
		return ResponseMsg{
			Message: "no logs found",
			Data:    logs,
		}, Ok()
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errors.New(internalServerErrorMsg))
	}

	return ResponseMsg{
		Message: "Logs are found",
		Data:    logs,
	}, Ok()
}
