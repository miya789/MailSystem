package lab_mail

import (
	"LabMeeting/pkg/meeting_type"
	"LabMeeting/pkg/schedule"
	"fmt"
	"mime"
	"net/mail"
	"regexp"
	"strconv"
	"time"
)

var (
	ME_NAME_JP       string // 使用者の名字の日本語表記
	ME_NAME_EN       string // 使用者の名字の英語表記
	ME_GRADE         string // 使用者の学年
	REMINDER_NAME_JP string // リマインダーの名前の日本語表記
	REMINDER_NAME_EN string // リマインダーの名前の英語表記
	REMINDER_GRADE   string // リマインダーの学年

)

// getByMeetingType returns 2 mail address objects (to and bcc) based on "meeting_type.MeetingType".
func getByMeetingType(mt meeting_type.MeetingType) (*mail.Address, []*mail.Address) {
	switch mt {
	case meeting_type.TeamMEMS:
		return teamMEMS, nil
	case meeting_type.Executive:
		return executive, executive_bcc
	default:
		return nil, nil
	}
}

// メールで送信するメッセージの塊
type Message struct {
	from    *mail.Address     // 送信者
	to      *mail.Address     // to に複数は不可
	bcc     []*mail.Address   // BCCは基本的に複数
	subject string            // subjetc は headers に含まれるので先に設定する
	headers map[string]string // headers は subjectが決まってから設定する
	rawBody string            // メッセージの headers を取り除いた本文の部分
	body    string            // メールで送信する headers も含めた全文
}

// common functions

// buildHeader returns headers based on *mail.Address (from and to), []*mail.Address (bcc), subject.
func buildHeader(from, to *mail.Address, bcc []*mail.Address, subject string) map[string]string {
	headers := make(map[string]string)

	headers["From"] = from.String()
	headers["To"] = to.String()
	// bcc のみ複数扱う為にこうしている
	if bcc != nil {
		var bccText string
		for _, v := range bcc {
			bccText += v.String() + ","
		}
		headers["Bcc"] = bccText
	}
	headers["Subject"] = subject
	headers["Content-Type"] = "text/plain; charset=UTF-8"
	headers["Content-Transfer-Encoding"] = "8bit"
	headers["MIME-Version"] = "1.0"
	headers["Date"] = time.Now().Local().Format(time.RFC1123Z)

	return headers
}

// buildMessage returns message body based on headers map[string]string, body string.
func buildMessage(headers map[string]string, body string) string {
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	return message + "\r\n" + body
}

func buildReminderSubject(mtg meeting_type.MeetingType, mailSchedule *schedule.MailSchedule) string {
	wdays := [...]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	s, _ := time.Parse("2006/01/02", mailSchedule.StartDate)
	dateMessageEN := mailSchedule.StartDate + "(" + wdays[s.Weekday()] + ")"
	place := getMeetingPlace(mailSchedule.Location)

	subjStr := "The next " + mtg.CaptitalString() + " Meeting【" + dateMessageEN + " " + mailSchedule.StartTime + " - @" + place.jp + "】"
	return mime.QEncoding.Encode("utf-8", subjStr)
}

func buildReminderRawBody(mtg meeting_type.MeetingType, mailSchedule *schedule.MailSchedule, mailZoomSchedule *schedule.MailZoomSchedule) string {
	place := getMeetingPlace(mailSchedule.Location)
	var (
		// TODO: 共通化しても良いのでは?
		wdays = [...]string{"日", "月", "火", "水", "木", "金", "土"}

		meetingTime         = mailSchedule.StartTime
		meetingPlaceJP      = place.jp
		meetingPlaceEN      = place.en
		meetingZoomURL      = mailZoomSchedule.URL
		meetingZoomPassword = mailZoomSchedule.Password
	)

	// Prepare mailSchedule.StartDate
	s, _ := time.Parse("2006/01/02", mailSchedule.StartDate)

	// Prepare DATE_FOR_CONTENTS
	dateMessageJP := s.Format("01/02") + "(" + wdays[s.Weekday()] + ")"
	dateMessageEN := s.Weekday().String() + ", " + s.Month().String() + " " + strconv.Itoa(s.Day())

	body := ""
	body += mtg.CaptitalString() + "の皆様\n"
	body += "\n"
	body += REMINDER_GRADE + "の" + REMINDER_NAME_JP + "です．\n"
	if mailSchedule.Location == "Zoom" {
		body += "次回の" + mtg.CaptitalString() + " Meetingは" + dateMessageJP + meetingTime + " - @Zoomで行われます．\n"
		body += "  Zoom URL: " + meetingZoomURL + "\n"
		body += "  Zoom Password: " + meetingZoomPassword + "\n"
	} else {
		body += "次回の" + mtg.CaptitalString() + " Meetingは" + dateMessageJP + " " + meetingTime + " - @" + meetingPlaceJP + "で行われます．\n"
	}
	body += "尚，ミーティングに関する連絡はこちらのメーリングリストのメッセージ宛で返信をお願いいたします．\n"
	body += "よろしくお願いいたします．\n"
	body += "\n"
	body += "\n"
	// en
	body += "Dear " + mtg.CaptitalString() + " members,\n"
	body += "\n"
	body += "I'm " + REMINDER_GRADE + " " + REMINDER_NAME_EN + ".\n"
	if mailSchedule.Location == "Zoom" {
		body += "The next " + mtg.CaptitalString() + " Meeting is going to be held at the Zoom from " + meetingTime + " on " + dateMessageEN + ".\n"
		body += "  Zoom URL: " + meetingZoomURL + "\n"
		body += "  Zoom Password: " + meetingZoomPassword + "\n"
	} else {
		body += "The next " + mtg.CaptitalString() + " Meeting is going to be held at the " + meetingPlaceEN + " from " + meetingTime + " on " + dateMessageEN + "."
	}
	body += "Please attend the meeting.\n"
	body += "Thank you.\n"
	body += "\n"
	// signature
	body += "--\n"
	body += "Mita Lab. Meeting Reminder\n"
	body += "\n"
	body += "MAIL: hosa@if.t.u-tokyo.ac.jp" // TODO

	return body
}

// buildMinutesSubject は入力される筈の日付を元に件名を組み立てる
func buildMinutesSubject(date string) string {
	t, _ := time.Parse("20060102", date)
	subjStr := "Executive meeting 議事録 " + t.Format("2006/01/02")
	return mime.QEncoding.Encode("utf-8", subjStr)
}

// buildMinutesRawBody は議事録をメール送信する為に本文を組み立てる
// 尚，コメントアウトは削除する
func buildMinutesRawBody(msg string) string {
	rep := regexp.MustCompilePOSIX(`^// .*?$\n`)
	msgMail := rep.ReplaceAllString(msg, "")

	body := ""
	body += "Executiveの皆様\n"
	body += "\n"
	body += "議事録をお送りします．\n"
	body += "不備がございましたらご指摘ください．\n"
	body += "\n"
	body += "よろしくお願いいたします．\n"
	body += "\n"
	body += ME_NAME_JP + "\n"
	body += "\n"
	body += msgMail

	return body
}
