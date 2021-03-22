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
	log.Printf("Checking whether today is a holiday or not... (Today: %s)\n", now.Format(schedule.TimeLayout))
	if schedule.IsHoliday(now) {
		log.Printf("Today (%s) is a holiday, so finished.\n", now.Format(schedule.TimeLayout))
		return
	}
	log.Printf("Today (%s) is not a holiday, so continuing...\n", now.Format(schedule.TimeLayout))

	// Get next weekday
	t := now.AddDate(0, 0, 1)
	t = now.AddDate(0, 0, 3)
	t = schedule.GetNextWeekday(t)
	log.Printf("The next weekday is %s.\n", t.Format(schedule.TimeLayout))

	// Get the next meeting specified with next weekday
	ms, mz, err := schedule.GetScheduleBy(t, mtg)
	if err != nil {
		log.Println(err)
		log.Printf("The announced schedule (%s) do not exist, so finished.\n", t.Format(schedule.TimeLayout))
		return
	}
	log.Printf("The announced schedule (%s) is %s, so continuing...\n", t.Format(schedule.TimeLayout), ms)

	// Send reminder mail
	log.Printf("Sending reminder mail...\n")
	if err := lab_mail.SendMail(mtg, ms, mz); err != nil {
		log.Println(err)
		return
	}
}
