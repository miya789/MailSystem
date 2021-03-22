package meeting_type

type MeetingType int

const (
	Unknown MeetingType = iota
	TeamMEMS
	Executive
)

func (m MeetingType) String() string {
	switch m {
	case TeamMEMS:
		return "teamMEMS"
	case Executive:
		return "executive"
	default:
		return "unknown"
	}
}

func (m MeetingType) CaptitalString() string {
	switch m {
	case TeamMEMS:
		return "TeamMEMS"
	case Executive:
		return "Executive"
	default:
		return "unknown"
	}
}
