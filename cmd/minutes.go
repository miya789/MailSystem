package main

import (
	"LabMeeting/pkg/meeting_type"
	"LabMeeting/pkg/memswiki"
	"LabMeeting/pkg/redmine"
	"LabMeeting/pkg/schedule"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../config/.env"); err != nil {
		log.Println(fmt.Errorf("Failed to read \"../config/.env\""))
	}
	RECEPTION_URL := os.Getenv("RECEPTION_URL")
	NANOTECH_HELP_URL := os.Getenv("NANOTECH_HELP_URL")
	// TEST_URL := os.Getenv("TEST_URL")

	var proxy bool
	flag.BoolVar(&proxy, "p", false, "false\t(default)")
	flag.Parse()
	log.Printf("Setting proxy: \"%+v\"", proxy)

	log.Println("Getting issues...")
	r := &redmine.Redmine{
		UseProxy: proxy,
	}
	receptionIssues, err := r.GetIssues(RECEPTION_URL)
	if err != nil {
		return
	}
	nanotechHelpIssues, err := r.GetIssues(NANOTECH_HELP_URL)
	if err != nil {
		return
	}

	log.Println("Loading schedules...")
	cs, err := schedule.Read(meeting_type.Executive, schedule.Calendar)
	if err != nil {
		log.Println(err)
		return
	}
	calendarSchdules := cs.([]*schedule.CalendarSchedule)

	// now := time.Now()

	template, err := memswiki.WriteTemplate(receptionIssues, nanotechHelpIssues, calendarSchdules)
	if err != nil {
		fmt.Errorf("Failed to GetScheduleBy(): %w", err)
		return
	}

	outPth := "../out/executive_minutes.txt"
	if err := ioutil.WriteFile(outPth, []byte(template), 0666); err != nil {
		log.Println(fmt.Errorf("Failed to Write(): %w", err))
		return
	}

	return
}
