package lab_mail

import (
	"LabMeeting/pkg/meeting_type"
	"LabMeeting/pkg/schedule"
	"fmt"
	"log"
	"mime"
	"net/mail"
	"strconv"
	"time"
)

type ReminderMail struct {
	userID       string
	mailPassword string
	grade        string
	name_jp      string
	name_en      string

	mtg              meeting_type.MeetingType
	mailSchedule     *schedule.MailSchedule
	mailZoomSchedule *schedule.MailZoomSchedule

	from *mail.Address
	to   *mail.Address
	bccs []*mail.Address

	header  string
	subject string
	body    string
	message string
}

func getMeetingPlace(key string) *MeetingPlace {
	meetingPlaces := map[string]*MeetingPlace{
		"113": {
			jp: "工学部3号館 113号室 (電気系セミナー室3) ",
			en: "Bldg. 3 Room 113 (Seminar 3)",
		},
		"114": {
			jp: "工学部3号館 114号室 (電気系セミナー室2) ",
			en: "Bldg. 3 Room 114 (Seminar 2)",
		},
		"128": {
			jp: "工学部3号館128号室 (電気系セミナー室1) ",
			en: "Bldg. 3 Room 128 (Seminar 1)",
		},
		"VDEC306": {
			jp: "VDEC 306",
			en: "VDEC 306",
		},
		"VDEC402": {
			jp: "VDEC 402",
			en: "VDEC 402",
		},
		"Bldg13": {
			jp: "13号館一般実験室",
			en: "Bldg. 13",
		},
	}

	v, ok := meetingPlaces[key]
	if !ok {
		return &MeetingPlace{
			jp: key,
			en: key,
		}
	}

	return v
}

func (r *ReminderMail) setAddress() {
	r.from = me
	switch r.mtg {
	case meeting_type.TeamMEMS:
		r.to = teamMEMS
	case meeting_type.Executive:
		r.to = executive
		r.bccs = executive_bccs
	default:
	}

	log.Printf("Setting from: %s", r.from)
	log.Printf("Setting to: %s", r.to)
	log.Printf("Setting bcc: %s", r.bccs)
}

func buildHaeder(from, to *mail.Address, bccs []*mail.Address, subject, body string) string {
	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	// bcc のみ複数扱う為にこうしている
	if bccs != nil {
		var bccText string
		for _, v := range bccs {
			bccText += v.String() + ","
		}
		headers["Bcc"] = bccText
	}
	headers["Subject"] = subject
	headers["Content-Type"] = "text/plain; charset=UTF-8"
	headers["Content-Transfer-Encoding"] = "8bit"
	headers["MIME-Version"] = "1.0"

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	return message + "\r\n" + body
}

func (r *ReminderMail) buildSubject() error {
	wdays := [...]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

	s, err := time.Parse("2006/01/02", r.mailSchedule.StartDate)
	if err != nil {
		return err
	}
	dateMessageEN := r.mailSchedule.StartDate + "(" + wdays[s.Weekday()] + ")"

	subjStr := "The next " + r.mtg.CaptitalString() + " Meeting【" + dateMessageEN + " " + r.mailSchedule.StartTime + " - @" + getMeetingPlace(r.mailSchedule.Location).jp + "】"
	log.Println(subjStr)
	r.subject = mime.QEncoding.Encode("utf-8", subjStr)
	return nil
}

func (r *ReminderMail) buildBody() error {
	place := getMeetingPlace(r.mailSchedule.Location)
	var (
		wdays = [...]string{"日", "月", "火", "水", "木", "金", "土"}

		meetingTime         = r.mailSchedule.StartTime
		meetingPlaceJP      = place.jp
		meetingPlaceEN      = place.en
		meetingZoomURL      = r.mailZoomSchedule.URL
		meetingZoomPassword = r.mailZoomSchedule.Password
	)

	// Prepare r.mailSchedule.StartDate
	s, err := time.Parse("2006/01/02", r.mailSchedule.StartDate)
	if err != nil {
		return err
	}

	// Prepare DATE_FOR_CONTENTS
	dateMessageJP := s.Format("01/02") + "(" + wdays[s.Weekday()] + ")"
	dateMessageEN := s.Weekday().String() + ", " + s.Month().String() + " " + strconv.Itoa(s.Day())

	body := ""
	body += r.mtg.CaptitalString() + "の皆様\n"
	body += "\n"
	body += r.grade + "の" + r.name_en + "です．\n"
	if r.mailSchedule.Location == "Zoom" {
		body += "次回の" + r.mtg.CaptitalString() + " Meetingは" + dateMessageJP + meetingTime + " - @Zoomで行われます．\n"
		body += "  Zoom URL: " + meetingZoomURL + "\n"
		body += "  Zoom Password: " + meetingZoomPassword + "\n"
	} else {
		body += "次回の" + r.mtg.CaptitalString() + " Meetingは" + dateMessageJP + " " + meetingTime + " - @" + meetingPlaceJP + "で行われます．\n"
	}
	body += "尚，ミーティングに関する連絡はこちらのメーリングリストのメッセージ宛で返信をお願いいたします．\n"
	body += "よろしくお願いいたします．\n"
	body += "\n"
	body += "\n"
	// en
	body += "Dear " + r.mtg.CaptitalString() + " members,\n"
	body += "\n"
	body += "I'm " + r.grade + " " + r.name_en + ".\n"
	if r.mailSchedule.Location == "Zoom" {
		body += "The next " + r.mtg.CaptitalString() + " Meeting is going to be held at the Zoom from " + meetingTime + " on " + dateMessageEN + ".\n"
		body += "  Zoom URL: " + meetingZoomURL + "\n"
		body += "  Zoom Password: " + meetingZoomPassword + "\n"
	} else {
		body += "The next " + r.mtg.CaptitalString() + " Meeting is going to be held at the " + meetingPlaceEN + " from " + meetingTime + " on " + dateMessageEN + "."
	}
	body += "Please attend the meeting.\n"
	body += "Thank you.\n"
	body += "\n"
	// signature
	body += "--\n"
	body += "Mita Lab. Meeting Reminder\n"
	body += "\n"
	body += "MAIL: hosa@if.t.u-tokyo.ac.jp"

	r.body = body

	return nil
}

func (r *ReminderMail) buildMessage() error {
	r.setAddress()

	// Setup message
	if err := r.buildBody(); err != nil {
		return err
	}
	if err := r.buildSubject(); err != nil {
		return err
	}
	r.message = buildHaeder(r.from, r.to, r.bccs, r.subject, r.body)
	log.Printf("The send mail is as follows...\n--------------------------------------------------\n%s\n--------------------------------------------------\n", r.message)

	return nil
}
