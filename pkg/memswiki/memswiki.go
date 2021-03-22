package memswiki

import (
	"LabMeeting/pkg/redmine"
	"LabMeeting/pkg/schedule"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/joho/godotenv"
)

// only executive
func WriteTemplate(receptionIssues, nanotechHelpIssues []redmine.Issue, calendarSchdules []*schedule.CalendarSchedule) (string, error) {
	if err := godotenv.Load("../config/.env"); err != nil {
		log.Println(fmt.Errorf("Failed to read \"../config/.env\""))
	}
	MEMSWIKI_LINK_URL := os.Getenv("MEMSWIKI_LINK_URL")
	var (
		receptionLinks    = ""
		nanotechHelpLinks = ""
		dates             = ""
	)
	for _, receptionIssue := range receptionIssues {
		receptionLinks += "- " + "[[#" + strconv.Itoa(receptionIssue.ID) + " " + receptionIssue.Subject + ">" + MEMSWIKI_LINK_URL + strconv.Itoa(receptionIssue.ID) + "]]\n"
	}
	for _, nanotechHelpIssue := range nanotechHelpIssues {
		nanotechHelpLinks += "- " + "[[#" + strconv.Itoa(nanotechHelpIssue.ID) + " " + nanotechHelpIssue.Subject + ">" + MEMSWIKI_LINK_URL + strconv.Itoa(nanotechHelpIssue.ID) + "]]\n"
	}
	for _, calendarSchdule := range calendarSchdules {
		if regexp.MustCompile(`https://`).MatchString(calendarSchdule.Location) && regexp.MustCompile(`zoom`).MatchString(calendarSchdule.Location) {
			dates += "," + calendarSchdule.StartDate + " " + calendarSchdule.StartTime + "~" + calendarSchdule.EndTime + "," + "Zoom" + "," + calendarSchdule.Subject + "\n"
		} else {
			dates += "," + calendarSchdule.StartDate + " " + calendarSchdule.StartTime + "~" + calendarSchdule.EndTime + "," + calendarSchdule.Location + "," + calendarSchdule.Subject + "\n"
		}
	}

	templateText := `[[Executive Meeting]]

#contents

*` + "2020/09/24" + ` ` + "10:00" + `- @` + "Zoom" + `
- 出席
// 全メンバーは以下

// 進行は基本的に以下の順だが，変更があればそれに合わせて変えよう
// 1. 共有事項
// 2. [[内部プロジェクト>http://mozart.if.t.u-tokyo.ac.jp:3000/issues/gantt]]
// 3. [[装置 (ガントチャートから)>http://mozart.if.t.u-tokyo.ac.jp:3000/issues/gantt]]
// 4. [[ナノテク受付窓口>http://mozart.if.t.u-tokyo.ac.jp:3000/projects/reception/issues?query_id=3]]
// 5. [[ナノテクへルプ>http://mozart.if.t.u-tokyo.ac.jp:3000/projects/nanotech_help/issues?query_id=4]]

**共有事項
// 書き方のテンプレ
// - XXX
// -- YYY

**ナノテクへルプ受付窓口
` + receptionLinks + `
**ナノテクヘルプ
` + nanotechHelpLinks + `
**装置メンテ

**設備

**講習会

**その他

**今後の予定
***Executive Meetings
,Date & Time,Location,Contents
` + dates

	log.Printf("\n--------------------------------------------------\n%s\n--------------------------------------------------\n", templateText)

	return templateText, nil
}
