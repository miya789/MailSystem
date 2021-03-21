package lab_mail

import (
	"fmt"
	"log"
	"net/mail"
	"os"

	"github.com/joho/godotenv"
)

var (
	me_email        string
	userID          string
	mailPassword    string
	teamMEMS_email  string
	executive_email string
	grade           string
	name_jp         string
	name_en         string

	me        *mail.Address
	teamMEMS  *mail.Address
	executive *mail.Address
	from      *mail.Address
	to        *mail.Address
	bcc       *mail.Address
)

const (
	portSMTP  = "587"
	portSMTPs = "465"
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
	me_email = os.Getenv("me_email")
	userID = os.Getenv("userID")
	teamMEMS_email = os.Getenv("teamMEMS_email")
	executive_email = os.Getenv("executive_email")
	mailPassword = os.Getenv("mailPassword")
	grade = os.Getenv("grade")
	name_jp = os.Getenv("name_jp")
	name_en = os.Getenv("name_en")

	me = &mail.Address{
		Name:    "Mail Reminder",
		Address: me_email,
	}
	teamMEMS = &mail.Address{
		Name:    "",
		Address: teamMEMS_email,
	}
	executive = &mail.Address{
		Name:    "",
		Address: executive_email,
	}

	from = me
	to = me
	bcc = me

	return
}
