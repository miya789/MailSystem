package lab_mail

import (
	"LabMeeting/pkg/meeting_type"
	"LabMeeting/pkg/schedule"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
)

func sendSMTPMail(host, message string) error {
	server := net.JoinHostPort(host, portSMTP)

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	c, err := smtp.Dial(server)
	if err != nil {
		return fmt.Errorf("Failed to tls.Dial(): %w", err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		return fmt.Errorf("Failed to c.Mail(): %w", err)
	}
	if err = c.Rcpt(to.Address); err != nil {
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

func sendSMTPMailSSL(host, message string) error {
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
	if err = c.Auth(smtp.PlainAuth("", userID, mailPassword, host)); err != nil {
		return fmt.Errorf("Failed to c.Auth(): %w", err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		return fmt.Errorf("Failed to c.Mail(): %w", err)
	}
	if err = c.Rcpt(to.Address); err != nil {
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
	message, err := buildMessage(mtg, ms, mzs)
	if err != nil {
		return err
	}

	// Connect to the SMTP Server
	mozartHost := "smtp.if.t.u-tokyo.ac.jp"
	return sendSMTPMail(mozartHost, message)
	// return sendSMTPMailSSL(mozartHost, message)
}
