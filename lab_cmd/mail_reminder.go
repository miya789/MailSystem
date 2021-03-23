package lab_cmd

import (
	"LabMeeting/pkg/lab_flag"
	"LabMeeting/pkg/lab_mail"
	"LabMeeting/pkg/schedule"
	"log"
	"time"
)

func SendReminderMail() {
	// Determine meeting type
	mtg, err := lab_flag.GetMeetingType()
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
	t := schedule.GetNextWeekday(now)
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
	if err := lab_mail.SendReminderMail(mtg, ms, mz); err != nil {
		log.Println(err)
		return
	}
}
