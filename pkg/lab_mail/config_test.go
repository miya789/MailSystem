package lab_mail

import (
	"reflect"
	"testing"
)

func Test_getMeetingPlace(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want *MeetingPlace
	}{
		{
			name: "should success",
			args: args{
				key: "113",
			},
			want: &MeetingPlace{
				jp: "工学部3号館 113号室 (電気系セミナー室3) ",
				en: "Bldg. 3 Room 113 (Seminar 3)",
			},
		},
		{
			name: "should success",
			args: args{
				key: "114",
			},
			want: &MeetingPlace{
				jp: "工学部3号館 114号室 (電気系セミナー室2) ",
				en: "Bldg. 3 Room 114 (Seminar 2)",
			},
		},
		{
			name: "should success",
			args: args{
				key: "128",
			},
			want: &MeetingPlace{
				jp: "工学部3号館128号室 (電気系セミナー室1) ",
				en: "Bldg. 3 Room 128 (Seminar 1)",
			},
		},
		{
			name: "should success",
			args: args{
				key: "VDEC306",
			},
			want: &MeetingPlace{
				jp: "VDEC 306",
				en: "VDEC 306",
			},
		},
		{
			name: "should success",
			args: args{
				key: "VDEC402",
			},
			want: &MeetingPlace{
				jp: "VDEC 402",
				en: "VDEC 402",
			},
		},
		{
			name: "should success",
			args: args{
				key: "Bldg13",
			},
			want: &MeetingPlace{
				jp: "13号館一般実験室",
				en: "Bldg. 13",
			},
		},
		{
			name: "should success",
			args: args{
				key: "zoom",
			},
			want: &MeetingPlace{
				jp: "zoom",
				en: "zoom",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMeetingPlace(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getMeetingPlace() = %v, want %v", got, tt.want)
			}
		})
	}
}
