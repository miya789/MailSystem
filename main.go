package main

import (
	"LabMeeting/lab_cmd"
	"LabMeeting/pkg/meeting_type"
	"flag"
	"fmt"
	"log"
	"os"
)

type cmdOptionType int

const (
	MinutesTemplateGenerator = iota
	MinutesSender
	MailReminder
)

var (
	cmd      int
	mtg      int
	useProxy bool
	useSSL   bool
)

// getFlags
// 本当はこんなもの使いたくないが仕方無く全員ここで取得
func getFlags() (cmdOptionType, meeting_type.MeetingType, bool, bool, error) {
	flag.IntVar(&cmd, "cmd", -1, "-1:\t(default, but do not use)\n0:\tminutes template generator\n1:\tminutes sender to wiki and mailer\n2:\tmail reminder\n")
	flag.IntVar(&mtg, "mtg", -1, "0: others \t(default, but do not use)\n1: TeamMEMS\n2: Executive\n(When -cmd 2, you can use)\n")
	flag.BoolVar(&useProxy, "p", false, "false:\tnot use proxy (default)\ntrue:\tuse proxy\n(When -cmd 0 and 1, you can use)\n")
	flag.BoolVar(&useSSL, "s", false, "false:\tnot use SSL (default)\ntrue:\tuse SSL and need your mail password in .env\n(When -cmd 1 and 2, you can use)\n")
	flag.Parse()

	return cmdOptionType(cmd), meeting_type.MeetingType(mtg), useProxy, useSSL, nil
}

func main() {
	cmd, mtg, useProxy, useSSL, err := getFlags()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	switch cmd {
	case MinutesTemplateGenerator:
		lab_cmd.GenerateMinutesTemplate(useProxy)
	case MinutesSender:
		lab_cmd.SendMinutes(useProxy, useSSL)
	case MailReminder:
		// この場合のみ変なパラメータならば処理を中止
		if mtg <= meeting_type.Unknown || mtg > meeting_type.Executive {
			err := fmt.Errorf("invalid value %q for flag -%s: %v", mtg, "mtg", "parse error")
			fmt.Fprintln(flag.CommandLine.Output(), err)
			fmt.Fprintf(flag.NewFlagSet(os.Args[0], flag.ExitOnError).Output(), "Usage of %s:\n", os.Args[0])
			flag.PrintDefaults()
			os.Exit(1)
		}
		lab_cmd.SendReminderMail(mtg, useSSL)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}
