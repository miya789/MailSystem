package lab_mail

import (
	"LabMeeting/pkg/meeting_type"
	"LabMeeting/pkg/schedule"
	"fmt"
	"log"
	"net/mail"
	"strconv"
	"time"
)

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

func buildBody(from, to *mail.Address, subject, body string) string {
	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Bcc"] = bcc.String()
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

func buildMessage(mtg meeting_type.MeetingType, mailSchedule *schedule.MailSchedule, mailZoomSchedule *schedule.MailZoomSchedule) (string, error) {
	place := getMeetingPlace(mailSchedule.Location)
	var (
		wdays = [...]string{"日", "月", "火", "水", "木", "金", "土"}

		meetingTime         = mailSchedule.StartTime
		meetingPlaceJP      = place.jp
		meetingPlaceEN      = place.en
		meetingZoomURL      = mailZoomSchedule.URL
		meetingZoomPassword = mailZoomSchedule.Password
	)

	// Prepare mailSchedule.StartDate
	s, err := time.Parse("2006/01/02", mailSchedule.StartDate)
	if err != nil {
		return "", err
	}

	// Prepare DATE_FOR_CONTENTS
	dateMessageJP := s.Format("01/02") + "(" + wdays[s.Weekday()] + ")"
	dateMessageEN := s.Weekday().String() + ", " + s.Month().String() + " " + strconv.Itoa(s.Day())

	subj := "The next " + mtg.CaptitalString() + " Meeting【" + dateMessageEN + " " + meetingTime + " - @" + meetingPlaceJP + "}】"

	body := ""
	body += mtg.CaptitalString() + "の皆様\n"
	body += "\n"
	body += grade + "の" + name_en + "です．\n"

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

	body += "Dear " + mtg.CaptitalString() + " members,\n"
	body += "\n"
	body += "I'm " + grade + " " + name_en + ".\n"
	if mailSchedule.Location == "Zoom" {
		body += "The next " + mtg.CaptitalString() + " Meeting is going to be held at the Zoom from " + meetingTime + " on " + dateMessageEN + ".\n"
		body += "  Zoom URL: " + meetingZoomURL + "\n"
		body += "  Zoom Password: " + meetingZoomPassword + "\n"
	} else {
		body += "The next " + mtg.CaptitalString() + " Meeting is going to be held at the " + meetingPlaceEN + " from " + meetingTime + " on " + dateMessageEN + "."
	}
	body += "\n"
	body += "Please attend the meeting.\n"
	body += "Thank you.\n"
	body += "\n"

	// signature
	body += "--\n"
	body += "Mita Lab. Meeting Reminder\n"
	body += "\n"
	body += "MAIL: hosa@if.t.u-tokyo.ac.jp"

	// Setup message
	message := buildBody(from, to, subj, body)
	log.Printf("\n--------------------------------------------------\n%s\n--------------------------------------------------\n", message)

	return message, nil
}
