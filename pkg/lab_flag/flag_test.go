package lab_flag

// func TestGetMeetingType(t *testing.T) {
// 	type fileds struct {
// 		mtg int
// 	}
// 	tests := []struct {
// 		name    string
// 		fileds  fileds
// 		want    meeting_type.MeetingType
// 		wantErr bool
// 	}{
// 		{
// 			name: "should success",
// 			fileds: fileds{
// 				mtg: 1,
// 			},
// 			want:    meeting_type.Executive,
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// flag.CommandLine.Set("mtg", tt.fileds.mtg)

// 			got, err := GetMeetingType()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetMeetingType() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("GetMeetingType() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestGetUseProxy(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		want    bool
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := GetUseProxy()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetUseProxy() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("GetUseProxy() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
