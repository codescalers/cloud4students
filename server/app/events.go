package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/codescalers/cloud4students/internal"
	"github.com/codescalers/cloud4students/middlewares"
	"github.com/codescalers/cloud4students/models"
	"github.com/pkg/errors"
)

type role string

var (
	userRole   role = "User"
	adminRole  role = "Admin"
	systemRole role = "System"
)

func (a *App) logUserDelete(userID string) error {
	event := models.AuditEvent{
		UserID:    userID,
		Action:    "delete_user",
		Role:      string(userRole),
		Timestamp: time.Now(),
		Metadata:  fmt.Sprintf("User %v is deleted", userID),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logBalanceCharge(userID, currency string, balance float64) error {
	event := models.AuditEvent{
		UserID:    userID,
		Action:    "update_balance",
		Role:      string(userRole),
		Timestamp: time.Now(),
		Metadata: fmt.Sprintf(
			"Balance is charged with %v %v",
			balance, currency,
		),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logUserVoucherActivate(userID, currency, voucher string, balance uint64) error {
	event := models.AuditEvent{
		UserID:    userID,
		Action:    "apply_voucher",
		Role:      string(userRole),
		Timestamp: time.Now(),
		Metadata:  fmt.Sprintf("User activated a voucher %v with %v %v", voucher, balance, currency),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logUserVoucherApply(userID, currency string, balance uint64) error {
	event := models.AuditEvent{
		UserID:    userID,
		Action:    "apply_voucher",
		Role:      string(userRole),
		Timestamp: time.Now(),
		Metadata:  fmt.Sprintf("User applied for a voucher with %v %v", balance, currency),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logUserUpdate(userID string) error {
	event := models.AuditEvent{
		UserID:    userID,
		Action:    "update_user",
		Role:      string(userRole),
		Timestamp: time.Now(),
		Metadata:  "User data is updated successfully",
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logUserPasswordUpdate(userID string) error {
	event := models.AuditEvent{
		UserID:    userID,
		Action:    "update_user_password",
		Role:      string(userRole),
		Timestamp: time.Now(),
		Metadata:  "Password is updated successfully",
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logUserSignedIn(userID string) error {
	event := models.AuditEvent{
		UserID:    userID,
		Action:    "signin_user",
		Role:      string(userRole),
		Timestamp: time.Now(),
		Metadata:  fmt.Sprintf("User %v is signed in successfully", userID),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logUserCreated(userID string) error {
	event := models.AuditEvent{
		UserID:    userID,
		Action:    "create_user",
		Role:      string(userRole),
		Timestamp: time.Now(),
		Metadata:  fmt.Sprintf("User %v is created", userID),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logVoucherReset(userID string, voucherBalance float64) error {
	event := models.AuditEvent{
		UserID:    userID,
		Action:    "reset_voucher",
		Role:      string(adminRole),
		Timestamp: time.Now(),
		Metadata:  fmt.Sprintf("Voucher balance %v is reset", voucherBalance),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logVoucherUpdate(userID, voucher string, balance uint64, approved bool) error {
	state := "Approved"
	if approved {
		state = "Rejected"
	}

	event := models.AuditEvent{
		UserID:    userID,
		Action:    "update_voucher",
		Role:      string(adminRole),
		Timestamp: time.Now(),
		Metadata:  fmt.Sprintf("Voucher `%v` with balance %v, is %v", voucher, balance, state),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logVoucherCreate(userID, voucher string, balance uint64) error {
	event := models.AuditEvent{
		UserID:    userID,
		Action:    "create_voucher",
		Role:      string(adminRole),
		Timestamp: time.Now(),
		Metadata:  fmt.Sprintf("Voucher `%v` with balance %v, is created successfully", voucher, balance),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logCardDelete(userID, last4digits string) error {
	event := models.AuditEvent{
		UserID:    userID,
		Action:    "delete_card",
		Role:      string(userRole),
		Timestamp: time.Now(),
		Metadata:  fmt.Sprintf("Card ending in %v is deleted", last4digits),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logCardDefaultSet(userID, last4digits string) error {
	event := models.AuditEvent{
		UserID:    userID,
		Action:    "set_default_card",
		Role:      string(userRole),
		Timestamp: time.Now(),
		Metadata:  fmt.Sprintf("Card ending in %v is set as default", last4digits),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logCardAdded(userID, last4digits string) error {
	event := models.AuditEvent{
		UserID:    userID,
		Action:    "add_card",
		Role:      string(userRole),
		Timestamp: time.Now(),
		Metadata:  fmt.Sprintf("Card ending in %v is added", last4digits),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logVoucherBalanceUpdate(userID, currency string, role role, balance float64) error {
	event := models.AuditEvent{
		UserID:    userID,
		Action:    "update_voucher_balance",
		Role:      string(role),
		Timestamp: time.Now(),
		Metadata: fmt.Sprintf(
			"Voucher balance is updated to %v %v",
			balance, currency,
		),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logBalanceUpdate(userID, currency string, role role, balance float64) error {
	event := models.AuditEvent{
		UserID:    userID,
		Action:    "update_balance",
		Role:      string(role),
		Timestamp: time.Now(),
		Metadata: fmt.Sprintf(
			"Balance is updated to %v %v",
			balance, currency,
		),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logK8sDelete(userID string, role role, k8sID int, createdAt time.Time) error {
	event := models.AuditEvent{
		UserID:    userID,
		Action:    "delete_k8s",
		Role:      string(role),
		Timestamp: time.Now(),
		Metadata: fmt.Sprintf(
			"Kubernetes %v which created at %v, is deleted", k8sID, createdAt.Format("January 2, 2006"),
		),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logVMDelete(userID string, role role, vmID int, createdAt time.Time) error {
	event := models.AuditEvent{
		UserID:    userID,
		Action:    "delete_vm",
		Role:      string(role),
		Timestamp: time.Now(),
		Metadata: fmt.Sprintf(
			"Virtual machine %v which created at %v, is deleted",
			vmID, createdAt.Format("January 2, 2006"),
		),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logInvoiceCreate(userID, currency string, invoiceID int, invoiceTotal float64, createdAt time.Time) error {
	event := models.AuditEvent{
		UserID:    userID,
		Action:    "create_invoice",
		Role:      string(systemRole),
		Timestamp: time.Now(),
		Metadata: fmt.Sprintf(
			"Invoice %v with value: %v %v is created at %v",
			invoiceID, invoiceTotal, currency, createdAt.Format("January 2, 2006"),
		),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logInvoicePayment(userID, currency string, invoiceTotal float64, paymentDetails models.PaymentDetails) error {
	event := models.AuditEvent{
		UserID:    userID,
		Action:    "pay_invoice",
		Role:      string(userRole),
		Timestamp: time.Now(),
		Metadata: fmt.Sprintf(
			"Invoice %v with value: %v %v is paid using: {balance: %v, vouchers: %v, card:%v}",
			paymentDetails.InvoiceID, invoiceTotal, currency,
			paymentDetails.Balance, paymentDetails.VoucherBalance, paymentDetails.Card,
		),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logInvoicePDFUpdate(req *http.Request, invoiceID int) error {
	userID := req.Context().Value(middlewares.UserIDKey("UserID")).(string)

	event := models.AuditEvent{
		UserID:    userID,
		Action:    "update_invoice",
		Role:      string(userRole),
		Timestamp: time.Now(),
		Metadata:  fmt.Sprintf("Invoice %v pdf data is updated", invoiceID),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logInvoiceDownload(req *http.Request, invoiceID int) error {
	userID := req.Context().Value(middlewares.UserIDKey("UserID")).(string)

	event := models.AuditEvent{
		UserID:    userID,
		Action:    "download_invoice",
		Role:      string(userRole),
		Timestamp: time.Now(),
		Metadata:  fmt.Sprintf("Invoice %v is downloaded", invoiceID),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logEmailSent(req *http.Request, targetUserID, subject string) error {
	userID := req.Context().Value(middlewares.UserIDKey("UserID")).(string)

	event := models.AuditEvent{
		UserID:    userID,
		Action:    "send_email",
		Role:      string(adminRole),
		Timestamp: time.Now(),
		Metadata:  fmt.Sprintf("An email is sent to: %v, with subject: %v", targetUserID, subject),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logAnnouncementCreate(req *http.Request, subject string) error {
	userID := req.Context().Value(middlewares.UserIDKey("UserID")).(string)

	event := models.AuditEvent{
		UserID:    userID,
		Action:    "create_announcement",
		Role:      string(adminRole),
		Timestamp: time.Now(),
		Metadata:  fmt.Sprintf("An announcement is created with subject: %v", subject),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logAdminSet(req *http.Request, adminID string, admin bool) error {
	userID := req.Context().Value(middlewares.UserIDKey("UserID")).(string)

	metaData := fmt.Sprintf("A new admin %v is added", adminID)
	if !admin {
		metaData = fmt.Sprintf("An admin %v is removed", adminID)
	}

	event := models.AuditEvent{
		UserID:    userID,
		Action:    "set_admin",
		Role:      string(adminRole),
		Timestamp: time.Now(),
		Metadata:  metaData,
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logNextLaunchUpdate(req *http.Request, on bool) error {
	userID := req.Context().Value(middlewares.UserIDKey("UserID")).(string)

	event := models.AuditEvent{
		UserID:    userID,
		Action:    "update_next_launch",
		Role:      string(adminRole),
		Timestamp: time.Now(),
		Metadata:  fmt.Sprintf("Next launch value is updated to: %v", on),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logMaintenanceUpdate(req *http.Request, on bool) error {
	userID := req.Context().Value(middlewares.UserIDKey("UserID")).(string)

	event := models.AuditEvent{
		UserID:    userID,
		Action:    "update_maintenance",
		Role:      string(adminRole),
		Timestamp: time.Now(),
		Metadata:  fmt.Sprintf("Maintenance value is updated to: %v", on),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logAllDeploymentsDelete(req *http.Request) error {
	userID := req.Context().Value(middlewares.UserIDKey("UserID")).(string)

	event := models.AuditEvent{
		UserID:    userID,
		Action:    "delete_all_deployments",
		Role:      string(adminRole),
		Timestamp: time.Now(),
		Metadata:  "All virtual machines are deleted",
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}

func (a *App) logVMsPriceUpdate(req *http.Request, prices internal.Prices) error {
	userID := req.Context().Value(middlewares.UserIDKey("UserID")).(string)

	event := models.AuditEvent{
		UserID:    userID,
		Action:    "update_vms_prices",
		Role:      string(adminRole),
		Timestamp: time.Now(),
		Metadata: fmt.Sprintf(
			"Virtual machines prices are updated {small: %v, medium: %v, large: %v, public IPs: %v}",
			prices.SmallVM, prices.MediumVM, prices.LargeVM, prices.PublicIP,
		),
	}

	if err := a.db.CreateAuditEvent(&event); err != nil {
		return errors.Wrapf(err, "Failed to log audit event: %s", event.Action)
	}

	return nil
}
