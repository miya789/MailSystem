package main

import (
	"LabMeeting/pkg/lab_mail"
	"LabMeeting/pkg/meeting_type"
	"LabMeeting/pkg/schedule"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func getMeetingType() (meeting_type.MeetingType, error) {
	var mtg int
	flag.IntVar(&mtg, "mtg", 0, "0: others \t(default, but do not use)\n1: TeamMEMS\n2: Executive")
	flag.Parse()
	if mtg <= int(meeting_type.Unknown) || mtg > int(meeting_type.Executive) {
		err := fmt.Errorf("invalid value %q for flag -%s: %v", strconv.Itoa(mtg), "mtg", "parse error")
		fmt.Fprintln(flag.CommandLine.Output(), err)
		fmt.Fprintf(flag.NewFlagSet(os.Args[0], flag.ExitOnError).Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return 0, err
	}
	return meeting_type.MeetingType(mtg), nil
}

func main() {
	// Determine meeting type
	mtg, err := getMeetingType()
	if err != nil {
		log.Println(err)
		return
	}

	// Check whether today is holiday
	now := time.Now()
	if schedule.IsHoliday(now) {
		log.Printf("Today is %s.\n", now.Format(schedule.TimeLayout))
		return
	}

	// Get next weekday
	t := now.AddDate(0, 0, 1)
	t = schedule.GetNextWeekday(t)
	log.Printf("Target date is %s.\n", t.Format(schedule.TimeLayout))

	// Get the next meeting specified with next weekday
	ms, mz, err := schedule.GetScheduleBy(t, mtg)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Next meeting is %s.\n", ms)

	// Send reminder mail
	if err := lab_mail.SendMail(mtg, ms, mz); err != nil {
		log.Println(err)
		return
	}
}
