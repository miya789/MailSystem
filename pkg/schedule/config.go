package schedule

// Schedule structs
type (
	ZoomSchedule struct {
		StartTime string `csv:"start_time"`
		Duration  string `csv:"duration"`
		Topic     string `csv:"topic"`
	}
	MailSchedule struct {
		Start
		Location string `csv:"Location"`
		Subject  string `csv:"Subject"`
	}
	MailZoomSchedule struct {
		Start
		URL      string `csv:"URL"`
		Password string `csv:"Password"`
	}
	CalendarSchedule struct {
		BaseSchedule
	}

	// Base componets
	BaseSchedule struct {
		Start
		End
		Location    string `csv:"Location"`
		Subject     string `csv:"Subject"`
		Description string `csv:"Description"`
	}
	Start struct {
		StartDate string `csv:"Start Date"`
		StartTime string `csv:"Start Time"`
	}
	End struct {
		EndDate string `csv:"End Date"`
		EndTime string `csv:"End Time"`
	}
)

// ScheduleType
type ScheduleType int

const (
	Seed ScheduleType = iota
	Zoom
	Mail
	MailZoom
	Calendar
)

var basename = map[ScheduleType]string{
	Seed:     "seed.csv",
	Zoom:     "zoom.csv",
	Mail:     "mail.csv",
	MailZoom: "mail_zoom.csv",
	Calendar: "calendar.csv",
}

const TimeLayout = "2006/01/02"
