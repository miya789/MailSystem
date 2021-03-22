package lab_flag

import (
	"LabMeeting/pkg/meeting_type"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func GetMeetingType() (meeting_type.MeetingType, error) {
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
