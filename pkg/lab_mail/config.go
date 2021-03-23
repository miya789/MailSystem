package lab_mail

import (
	"LabMeeting/pkg/meeting_type"
	"LabMeeting/pkg/schedule"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	userID       string
	mailPassword string
	grade        string
	name_jp      string
	name_en      string

	me             *mail.Address
	teamMEMS       *mail.Address
	executive      *mail.Address
	executive_bccs []*mail.Address
)

const (
	PORT_SMTP  = "587"
	PORT_SMTPs = "465"
)

type MeetingPlace struct {
	jp string
	en string
}

func init() {
	if err := godotenv.Load("../config/.env"); err != nil {
		// GitLab-Runner で実行する為に失敗しても可とする
		// Local で実行する場合は .env によって簡単に設定できる
		log.Println(fmt.Errorf("Failed to lab_mail init(): failed to read \"../config/.env\""))
	}

	userID = os.Getenv("userID")
	mailPassword = os.Getenv("mailPassword")
	grade = os.Getenv("grade")
	name_jp = os.Getenv("name_jp")
	name_en = os.Getenv("name_en")

	// 予め用意する宛先
	me = &mail.Address{
		Name:    "Mail Reminder",
		Address: os.Getenv("me_email"),
	}
	teamMEMS = &mail.Address{
		Name:    "",
		Address: os.Getenv("teamMEMS_email"),
	}
	executive = &mail.Address{
		Name:    "",
		Address: os.Getenv("executive_email"),
	}
	// bcc のみ複数扱う為にこうしている
	bccs_address := os.Getenv("executive_email_bcc")
	if bccs_address != "" {
		for _, executive_email_bcc := range strings.Split(bccs_address, ",") {
			executive_bccs = append(executive_bccs, &mail.Address{
				Name:    "",
				Address: strings.TrimSpace(executive_email_bcc), // カンマ区切りで余計な空白があってもここで除去
			})
		}
	}

	return
}

// TODO: テストにしやすくするために切り出す
// 複数のbccになっているかとか?
func New(mtg meeting_type.MeetingType, ms *schedule.MailSchedule, mzs *schedule.MailZoomSchedule) *ReminderMail {
	return &ReminderMail{
		userID:       userID,
		mailPassword: mailPassword,
		grade:        grade,
		name_jp:      name_jp,
		name_en:      name_en,

		mtg:              mtg,
		mailSchedule:     ms,
		mailZoomSchedule: mzs,
	}
}
