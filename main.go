package main

import (
	"LabMeeting/lab_cmd"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

const usgaeCmd = `-1: (default, but do not use)
0:  minutes template generator
1:  minutes sender to wiki and mailer
2:  mail reminder
`

type cmdOptionType int

const (
	MinutesTemplateGenerator = iota
	MinutesSender
	MailReminder
)

func getCmdOption() (cmdOptionType, error) {
	var cmd int
	flag.IntVar(&cmd, "cmd", -1, usgaeCmd)
	flag.Parse()
	if cmd < int(MinutesTemplateGenerator) || cmd > int(MailReminder) {
		err := fmt.Errorf("invalid value %q for flag -%s: %v", strconv.Itoa(cmd), "cmd", "parse error")
		fmt.Fprintln(flag.CommandLine.Output(), err)
		fmt.Fprintf(flag.NewFlagSet(os.Args[0], flag.ExitOnError).Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return 0, err
	}
	return cmdOptionType(cmd), nil
}

func main() {
	cmd, err := getCmdOption()
	if err != nil {
		log.Println(err)
		return
	}
	switch cmd {
	case MinutesTemplateGenerator:
		lab_cmd.GenerateMinutesTemplate()
	case MinutesSender:
		lab_cmd.SendMinutes()
	case MailReminder:
		lab_cmd.SendReminderMail()
	default:
		return
	}
}
