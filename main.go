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
	useProxy bool
	mtg      int
)

// getFlags
// 本当はこんなもの使いたくないが仕方無く全員ここで取得
func getFlags() (cmdOptionType, bool, meeting_type.MeetingType, error) {
	flag.IntVar(&cmd, "cmd", -1, "-1:\t(default, but do not use)\n0:\tminutes template generator\n1:\tminutes sender to wiki and mailer\n2:\tmail reminder")
	flag.BoolVar(&useProxy, "p", false, "false\t(default)")
	flag.IntVar(&mtg, "mtg", -1, "0: others \t(default, but do not use)\n1: TeamMEMS\n2: Executive")
	flag.Parse()

	return cmdOptionType(cmd), useProxy, meeting_type.MeetingType(mtg), nil
}

func main() {
	cmd, userProxy, mtg, err := getFlags()
	if err != nil {
		log.Println(err)
		return
	}

	if err != nil {
		log.Println(err)
		return
	}
	switch cmd {
	case MinutesTemplateGenerator:
		lab_cmd.GenerateMinutesTemplate(userProxy)
	case MinutesSender:
		lab_cmd.SendMinutes()
	case MailReminder:
		// この場合のみ変なパラメータならば処理を中止
		if mtg <= meeting_type.Unknown || mtg > meeting_type.Executive {
			err := fmt.Errorf("invalid value %q for flag -%s: %v", mtg, "mtg", "parse error")
			fmt.Fprintln(flag.CommandLine.Output(), err)
			fmt.Fprintf(flag.NewFlagSet(os.Args[0], flag.ExitOnError).Output(), "Usage of %s:\n", os.Args[0])
			flag.PrintDefaults()
			return
		}
		lab_cmd.SendReminderMail(mtg)
	default:
		return
	}
}
