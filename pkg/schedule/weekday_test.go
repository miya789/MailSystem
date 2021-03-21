package schedule

import (
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

func Test_getNextWeekday(t *testing.T) {
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
