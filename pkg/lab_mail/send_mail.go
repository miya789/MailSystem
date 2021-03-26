package lab_mail

import (
	"LabMeeting/pkg/meeting_type"
	"LabMeeting/pkg/schedule"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
)

var (
	ME_MAIL_ID         string // 使用者のSMTPsメール送信用のID
	ME_MAIL_PASS       string // 使用者のSMTPsメール送信用のパス
	REMINDER_MAIL_ID   string // リマインダーのSMTPsメール送信用のID
	REMINDER_MAIL_PASS string // リマインダーのSMTPsメール送信用のパス
)

const (
	PORT_SMTP  = "587" // 言わずもがな
	PORT_SMTPs = "465" // 言わずもがな
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
	if _, err = w.Write([]byte(message.body)); err != nil {
		return fmt.Errorf("Failed to w.Write(): %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("Failed to w.Close(): %w", err)
	}
	c.Quit()

	return nil
}

func SendReminderMail(mtg meeting_type.MeetingType, ms *schedule.MailSchedule, mzs *schedule.MailZoomSchedule, useSSL bool) error {
	to, bcc := getByMeetingType(mtg)
	message := &Message{
		from:    meetingReminder,
		to:      to,
		bcc:     bcc, // 複数であることに注意
		subject: buildReminderSubject(mtg, ms),
		rawBody: buildReminderRawBody(mtg, ms, mzs),
	}
	message.headers = buildHeader(message.from, message.to, message.bcc, message.subject)
	message.body = buildMessage(message.headers, message.rawBody)
	log.Printf("The send mail is as follows...\n--------------------------------------------------\n%s\n--------------------------------------------------\n", message.body)

	// Connect to the SMTP Server
	mozartHost := "smtp.if.t.u-tokyo.ac.jp"
	if useSSL {
		return sendSMTPMailSSL(mozartHost, message, ME_MAIL_ID, ME_MAIL_PASS)
	}
	return sendSMTPMail(mozartHost, message)
}

// コメントアウトした行は削除してメール送信する
func SendMinutesMail(mtg meeting_type.MeetingType, date, msg string, useSSL bool) error {
	to, _ := getByMeetingType(mtg)
	message := &Message{
		from:    me,
		to:      to,
		bcc:     nil,
		subject: buildMinutesSubject(date),
		rawBody: buildMinutesRawBody(msg),
	}
	message.headers = buildHeader(message.from, message.to, message.bcc, message.subject)
	message.body = buildMessage(message.headers, message.rawBody)
	log.Printf("The send mail is as follows...\n--------------------------------------------------\n%s\n--------------------------------------------------\n", message.body)

	// Connect to the SMTP Server
	mozartHost := "smtp.if.t.u-tokyo.ac.jp"
	if useSSL {
		return sendSMTPMailSSL(mozartHost, message, ME_MAIL_ID, ME_MAIL_PASS)
	}
	return sendSMTPMail(mozartHost, message)
}
