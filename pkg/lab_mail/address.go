package lab_mail

import (
	"fmt"
	"log"
	"net/mail"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	ME_EMAIL            string
	REMINDER_EMAIL      string
	TEAMMEMS_EMAIL      string
	EXECUTIVE_EMAIL     string
	EXECUTIVE_EMAIL_BCC string

	me              *mail.Address
	meetingReminder *mail.Address
	teamMEMS        *mail.Address
	executive       *mail.Address
	executive_bcc   []*mail.Address
)

func init() {
	if err := godotenv.Load("config/.env"); err != nil {
		// GitLab-Runner で実行する為に失敗しても可とする
		// Local で実行する場合は .env によって簡単に設定できる
		// 恐らくは， GitLab-Runner で読み込めても環境変数が優先されるが，一応は 変なものを混入させないように注意
		log.Println(fmt.Errorf("Failed to lab_mail init(): failed to read \"config/.env\""))
	}

	// ここで環境変数の読み込みを一通り終える
	ME_NAME_JP = os.Getenv("ME_NAME_JP")
	ME_NAME_EN = os.Getenv("ME_NAME_EN")
	ME_GRADE = os.Getenv("ME_GRADE")
	ME_MAIL_ID = os.Getenv("ME_MAIL_ID")
	ME_MAIL_PASS = os.Getenv("ME_MAIL_PASS")
	REMINDER_NAME_JP = os.Getenv("REMINDER_NAME_JP")
	REMINDER_NAME_EN = os.Getenv("REMINDER_NAME_EN")
	REMINDER_GRADE = os.Getenv("REMINDER_GRADE")
	REMINDER_MAIL_ID = os.Getenv("REMINDER_MAIL_ID")
	REMINDER_MAIL_PASS = os.Getenv("REMINDER_MAIL_PASS")

	// 予め用意する宛先
	me = &mail.Address{
		Name:    "",
		Address: os.Getenv("ME_EMAIL"),
	}
	meetingReminder = &mail.Address{
		Name:    "Mita Lab. Meeting Reminder",
		Address: os.Getenv("REMINDER_EMAIL"),
	}
	teamMEMS = &mail.Address{
		Name:    "",
		Address: os.Getenv("TEAMMEMS_EMAIL"),
	}
	executive = &mail.Address{
		Name:    "",
		Address: os.Getenv("EXECUTIVE_EMAIL"),
	}
	// bcc のみ複数扱う為にこうしている
	bcc_address := os.Getenv("EXECUTIVE_EMAIL_BCC")
	if bcc_address != "" {
		for _, executive_email_bcc := range strings.Split(bcc_address, ",") {
			executive_bcc = append(executive_bcc, &mail.Address{
				Name:    "",
				Address: strings.TrimSpace(executive_email_bcc), // カンマ区切りで余計な空白があってもここで除去
			})
		}
	}

	return
}
