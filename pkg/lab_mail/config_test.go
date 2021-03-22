package lab_mail

import (
	"LabMeeting/pkg/meeting_type"
	"LabMeeting/pkg/schedule"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		mtg meeting_type.MeetingType
		ms  *schedule.MailSchedule
		mzs *schedule.MailZoomSchedule
	}
	tests := []struct {
		name string
		args args
		want *ReminderMail
	}{
		{
			name: "should success",
			args: args{
				mtg: meeting_type.TeamMEMS,
				ms:  &schedule.MailSchedule{},
				mzs: &schedule.MailZoomSchedule{},
			},
			want: &ReminderMail{
				mtg:              meeting_type.TeamMEMS,
				mailSchedule:     &schedule.MailSchedule{},
				mailZoomSchedule: &schedule.MailZoomSchedule{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.mtg, tt.args.ms, tt.args.mzs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
