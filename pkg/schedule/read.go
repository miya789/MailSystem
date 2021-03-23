package schedule

import (
	"LabMeeting/pkg/meeting_type"
	"fmt"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

// Read returns array of some Schedule.
func Read(mt meeting_type.MeetingType, st ScheduleType) (interface{}, error) {
	outPth := "config/" + mt.String() + "_" + basename[st]

	seedFile, err := os.OpenFile(outPth, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("Failed to Read(): %w", err)
	}
	defer seedFile.Close()

	switch st {
	case Mail:
		mailSchedules := []*MailSchedule{}
		if err = gocsv.UnmarshalFile(seedFile, &mailSchedules); err != nil {
			return nil, fmt.Errorf("Failed to Read(): %w", err)
		}
		return mailSchedules, nil
	case MailZoom:
		mailZoomSchedules := []*MailZoomSchedule{}
		if err = gocsv.UnmarshalFile(seedFile, &mailZoomSchedules); err != nil {
			return nil, fmt.Errorf("Failed to Read(): %w", err)
		}
		return mailZoomSchedules, nil
	case Calendar:
		calendarSchedules := []*CalendarSchedule{}
		if err = gocsv.UnmarshalFile(seedFile, &calendarSchedules); err != nil {
			return nil, fmt.Errorf("Failed to Read(): %w", err)
		}
		return calendarSchedules, nil
	}

	return nil, nil
}

func GetSchedulesAfter(t time.Time, mtg meeting_type.MeetingType, st ScheduleType) (interface{}, error) {
	cs, err := Read(mtg, st)
	if err != nil {
		return nil, err
	}
	allCalendarSchdules := cs.([]*CalendarSchedule)
	var returnCalendarSchdules []*CalendarSchedule
	var layout = "2006/01/02"
	for _, v := range allCalendarSchdules {
		s, err := time.Parse(layout, v.StartDate)
		if err != nil {
			return nil, err
		}
		if s.After(t) {
			returnCalendarSchdules = append(returnCalendarSchdules, v)
		}
	}

	return returnCalendarSchdules, nil
}
