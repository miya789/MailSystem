package lab_mail

import (
	"LabMeeting/pkg/meeting_type"
	"LabMeeting/pkg/schedule"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
)

func (r *ReminderMail) sendSMTPMail(host, message string) error {
	server := net.JoinHostPort(host, portSMTP)

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	c, err := smtp.Dial(server)
	if err != nil {
		return fmt.Errorf("Failed to tls.Dial(): %w", err)
	}

	// To && From
	if err = c.Mail(r.from.Address); err != nil {
		return fmt.Errorf("Failed to c.Mail(): %w", err)
	}
	if err = c.Rcpt(r.to.Address); err != nil {
		return fmt.Errorf("Failed to c.Rcpt(): %w", err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("Failed to c.Data(): %w", err)
	}
	if _, err = w.Write([]byte(message)); err != nil {
		return fmt.Errorf("Failed to w.Write(): %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("Failed to w.Close(): %w", err)
	}
	c.Quit()

	return nil
}

func (r *ReminderMail) sendSMTPMailSSL(host, message string) error {
	server := net.JoinHostPort(host, portSMTPs)

	// TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", server, tlsConfig)
	if err != nil {
		return fmt.Errorf("Failed to tls.Dial(): %w", err)
	}
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("Failed to smtp.NewClient(): %w", err)
	}

	// Auth
	if err = c.Auth(smtp.PlainAuth("", r.userID, r.mailPassword, host)); err != nil {
		return fmt.Errorf("Failed to c.Auth(): %w", err)
	}

	// To && From
	if err = c.Mail(r.from.Address); err != nil {
		return fmt.Errorf("Failed to c.Mail(): %w", err)
	}
	if err = c.Rcpt(r.to.Address); err != nil {
		return fmt.Errorf("Failed to c.Rcpt(): %w", err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("Failed to c.Data(): %w", err)
	}
	if _, err = w.Write([]byte(message)); err != nil {
		return fmt.Errorf("Failed to w.Write(): %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("Failed to w.Close(): %w", err)
	}
	c.Quit()

	return nil
}

func SendMail(mtg meeting_type.MeetingType, ms *schedule.MailSchedule, mzs *schedule.MailZoomSchedule) error {
	r := New(mtg, ms, mzs)

	if err := r.buildMessage(); err != nil {
		return err
	}

	// Connect to the SMTP Server
	mozartHost := "smtp.if.t.u-tokyo.ac.jp"
	return r.sendSMTPMail(mozartHost, r.message)
}
