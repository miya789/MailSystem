package schedule

import (
	"LabMeeting/pkg/meeting_type"
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
)

// Read returns array of some Schedule.
func Read(mt meeting_type.MeetingType, st ScheduleType) (interface{}, error) {
	outPth := "../config/" + mt.String() + "_" + basename[st]

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
