package schedule

import (
	"LabMeeting/pkg/meeting_type"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/PuloV/ics-golang"
)

func IsHoliday(t time.Time) bool {
	parser := ics.New()
	pc := parser.GetInputChan()
	// このURLは Google カレンダーの「カレンダー設定」→「日本の祝日」→「ICAL」から取得可能 (2017/01/12現在)
	pc <- "https://calendar.google.com/calendar/ical/ja.japanese%23holiday%40group.v.calendar.google.com/public/basic.ics"
	parser.Wait()
	cal, err := parser.GetCalendars()
	if err != nil {
		// TODO: log level
		log.Println(fmt.Errorf("Failed to isNotWeekday(): %w", err))
		return false
	}

	isSat := (t.Weekday() == time.Saturday)
	isSun := (t.Weekday() == time.Sunday)
	var isPublicHoliday bool
	for _, e := range cal[0].GetEvents() {
		if t.Format(TimeLayout) == e.GetStart().Format(TimeLayout) {
			isPublicHoliday = true
			break
		}
	}

	return isSat || isSun || isPublicHoliday
}

func GetNextWeekday(t time.Time) time.Time {
	t = t.AddDate(0, 0, 1)
	for IsHoliday(t) {
		t = t.AddDate(0, 0, 1)
	}
	return t
}

func GetScheduleBy(targetDate time.Time, mtg meeting_type.MeetingType) (*MailSchedule, *MailZoomSchedule, error) {
	ms, err := Read(mtg, Mail)
	if err != nil {
		return nil, nil, err
	}
	mailSchedules := ms.([]*MailSchedule)
	mzs, err := Read(mtg, MailZoom)
	if err != nil {
		return nil, nil, err
	}
	mailZoomSchedules := mzs.([]*MailZoomSchedule)

	for i, mailSchedule := range mailSchedules {
		if mailSchedule.StartDate == targetDate.Format("2006/01/02") {
			if &mailSchedule.Start == &mailZoomSchedules[i].Start {
				return nil, nil, fmt.Errorf("Failed to GetScheduleBy(): %w", errors.New("Meeting file conflict"))
			}
			return mailSchedule, mailZoomSchedules[i], nil
		}
	}

	return nil, nil, fmt.Errorf("There is no meetings")
}
