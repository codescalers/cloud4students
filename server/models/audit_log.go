package models

import (
	"time"
)

// Struct to represent audit log details
type AuditLog struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     string    `json:"user_id" gorm:"not null"`
	Method     string    `json:"method"`
	URL        string    `json:"url"`
	Timestamp  time.Time `json:"timestamp"`
	Success    bool      `json:"success"`
	StatusCode int       `json:"status_code"`
}

func (d *DB) CreateAuditLog(l *AuditLog) error {
	return d.db.Create(&l).Error
}

// GetUserLogs gets user logs
func (d *DB) GetUserLogs(userID string) ([]AuditLog, error) {
	var res []AuditLog
	return res, d.db.Find(&res, "user_id = ?", userID).Error
}
