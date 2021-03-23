package main

import (
	"LabMeeting/cmd"
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
	var cmdOpt int
	flag.IntVar(&cmdOpt, "cmdOpt", -1, usgaeCmd)
	flag.Parse()
	if cmdOpt < int(MinutesTemplateGenerator) || cmdOpt > int(MailReminder) {
		err := fmt.Errorf("invalid value %q for flag -%s: %v", strconv.Itoa(cmdOpt), "cmdOpt", "parse error")
		fmt.Fprintln(flag.CommandLine.Output(), err)
		fmt.Fprintf(flag.NewFlagSet(os.Args[0], flag.ExitOnError).Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return 0, err
	}
	return cmdOptionType(cmdOpt), nil
}

func main() {
	cmdOpt, err := getCmdOption()
	if err != nil {
		log.Println(err)
		return
	}
	switch cmdOpt {
	case MinutesTemplateGenerator:
		cmd.GenerateMinutesTemplate()
	case MinutesSender:
		cmd.SendMinutes()
	case MailReminder:
		cmd.SendReminderMail()
	default:
		return
	}
}
