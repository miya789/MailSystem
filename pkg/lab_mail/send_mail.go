package lab_mail

import (
	"LabMeeting/pkg/meeting_type"
	"LabMeeting/pkg/schedule"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
)

func sendSMTPMail(host string, message *Message) error {
	server := net.JoinHostPort(host, PORT_SMTP)

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	c, err := smtp.Dial(server)
	if err != nil {
		return fmt.Errorf("Failed to tls.Dial(): %w", err)
	}

	// To && From
	if err = c.Mail(message.from.Address); err != nil {
		return fmt.Errorf("Failed to c.Mail(): %w", err)
	}
	if err = c.Rcpt(message.to.Address); err != nil {
		return fmt.Errorf("Failed to c.Rcpt(): %w", err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("Failed to c.Data(): %w", err)
	}
	if _, err = w.Write([]byte(message.body)); err != nil {
		return fmt.Errorf("Failed to w.Write(): %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("Failed to w.Close(): %w", err)
	}
	c.Quit()

	return nil
}

func sendSMTPMailSSL(host string, message *Message, userID, password string) error {
	server := net.JoinHostPort(host, PORT_SMTPs)

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
	if err = c.Auth(smtp.PlainAuth("", userID, password, host)); err != nil {
		return fmt.Errorf("Failed to c.Auth(): %w", err)
	}

	// To && From
	if err = c.Mail(message.from.Address); err != nil {
		return fmt.Errorf("Failed to c.Mail(): %w", err)
	}
	if err = c.Rcpt(message.to.Address); err != nil {
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

func SendReminderMail(mtg meeting_type.MeetingType, ms *schedule.MailSchedule, mzs *schedule.MailZoomSchedule) error {
	r := New(mtg, ms, mzs)

	if err := r.buildReminderMessage(); err != nil {
		return err
	}

	// Connect to the SMTP Server
	mozartHost := "smtp.if.t.u-tokyo.ac.jp"
	return r.sendSMTPMail(mozartHost, r.message)
}

// コメントアウトした行は削除してメール送信する
func SendMinutesMail(mtg meeting_type.MeetingType, date, msg string) error {
	// TODO: ミーティングスケジュールは不要なので仕様を変える
	r := New(mtg, nil, nil)

	if err := r.buildMessage(msg, date); err != nil {
		return err
	}

	// Connect to the SMTP Server
	mozartHost := "smtp.if.t.u-tokyo.ac.jp"
	return r.sendSMTPMailSSL(mozartHost, r.message)
}
