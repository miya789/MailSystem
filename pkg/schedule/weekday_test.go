package schedule

import (
	"LabMeeting/pkg/meeting_type"
	"reflect"
	"testing"
	"time"
)

func Test_IsHoliday(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should return true with 2021/5/1 as holiday",
			args: args{
				t: time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
			},
			want: true,
		},
		{
			name: "should return true with 2021/5/2 as holiday",
			args: args{
				t: time.Date(2021, 5, 2, 0, 0, 0, 0, time.Local),
			},
			want: true,
		},
		{
			name: "should return true with 2021/5/3 as holiday",
			args: args{
				t: time.Date(2021, 5, 3, 0, 0, 0, 0, time.Local),
			},
			want: true,
		},
		{
			name: "should return false with 2021/5/6 as weekday",
			args: args{
				t: time.Date(2021, 5, 6, 0, 0, 0, 0, time.Local),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsHoliday(tt.args.t); got != tt.want {
				t.Errorf("isHoliday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetNextWeekday(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "should return true with 2021/5/1 as holiday",
			args: args{
				t: time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local),
			},
			want: time.Date(2021, 5, 6, 0, 0, 0, 0, time.Local),
		},
		{
			name: "should return true with 2021/5/2 as holiday",
			args: args{
				t: time.Date(2021, 5, 2, 0, 0, 0, 0, time.Local),
			},
			want: time.Date(2021, 5, 6, 0, 0, 0, 0, time.Local),
		},
		{
			name: "should return true with 2021/5/3 as holiday",
			args: args{
				t: time.Date(2021, 5, 3, 0, 0, 0, 0, time.Local),
			},
			want: time.Date(2021, 5, 6, 0, 0, 0, 0, time.Local),
		},
		{
			name: "should return false with 2021/5/6 as weekday",
			args: args{
				t: time.Date(2021, 5, 6, 0, 0, 0, 0, time.Local),
			},
			want: time.Date(2021, 5, 7, 0, 0, 0, 0, time.Local),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNextWeekday(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNextWeekday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetScheduleBy(t *testing.T) {
	type args struct {
		targetDate time.Time
		mtg        meeting_type.MeetingType
	}
	tests := []struct {
		name    string
		args    args
		want    *MailSchedule
		want1   *MailZoomSchedule
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetScheduleBy(tt.args.targetDate, tt.args.mtg)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetScheduleBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetScheduleBy() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetScheduleBy() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
