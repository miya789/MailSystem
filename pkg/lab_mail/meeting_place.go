package lab_mail

type meetingPlace struct {
	jp string
	en string
}

func getMeetingPlace(key string) *meetingPlace {
	meetingPlaces := map[string]*meetingPlace{
		"113": {
			jp: "工学部3号館 113号室 (電気系セミナー室3) ",
			en: "Bldg. 3 Room 113 (Seminar 3)",
		},
		"114": {
			jp: "工学部3号館 114号室 (電気系セミナー室2) ",
			en: "Bldg. 3 Room 114 (Seminar 2)",
		},
		"128": {
			jp: "工学部3号館128号室 (電気系セミナー室1) ",
			en: "Bldg. 3 Room 128 (Seminar 1)",
		},
		"VDEC306": {
			jp: "VDEC 306",
			en: "VDEC 306",
		},
		"VDEC402": {
			jp: "VDEC 402",
			en: "VDEC 402",
		},
		"Bldg13": {
			jp: "13号館一般実験室",
			en: "Bldg. 13",
		},
	}

	v, ok := meetingPlaces[key]
	if !ok {
		return &meetingPlace{
			jp: key,
			en: key,
		}
	}

	return v
}
