package memswiki

import (
	"LabMeeting/pkg/redmine"
	"LabMeeting/pkg/schedule"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var (
	username            string
	password            string
	wikiTemplateMembers string
	proxyURL            *url.URL
	memsWiki            *url.URL
)

func init() {
	if err := godotenv.Load("config/.env"); err != nil {
		log.Println(fmt.Errorf("Failed to read \".env\""))
		return
	}
	username = os.Getenv("WIKI_USERNAME")
	password = os.Getenv("WIKI_PASSWORD")
	wikiTemplateMembers = os.Getenv("WIKI_TEMPLATE_MEMBERS")
	proxyURL, _ = url.Parse(os.Getenv("PROXY_URL"))
	memsWiki, _ = url.Parse(os.Getenv("MEMSWIKI_URL"))

	return
}

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
// ` + wikiTemplateMembers + `

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

	log.Printf("The generated template is as follows...\n--------------------------------------------------\n%s\n--------------------------------------------------\n", templateText)

	return templateText, nil
}

var stdHeader = map[string]string{
	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	"Accept-Encoding":           "gzip, deflate",
	"Accept-Language":           "ja,en-US;q=0.7,en;q=0.3",
	"Cache-Contorl":             "no-cache",
	"Connection":                "keep-alive",
	"Content-Length":            "55",
	"Host":                      "mozart.if.t.u-tokyo.ac.jp",
	"Origin":                    "http://mozart.if.t.u-tokyo.ac.jp",
	"Pragma":                    "no-cache",
	"Upgrade-Insecure-Requests": "1",
	"User-Agent":                "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:86.0) Gecko/20100101 Firefox/86.0",
}

func WriteMinute(date int, msg string) error {
	log.SetFlags(log.Lshortfile) // メチャクチャ難しかったのでログを詳細に出している
	log.Printf("Accessing %s...\n", memsWiki.Scheme+"://"+memsWiki.Host)

	page := "Executive Meeting/" + strconv.Itoa(date) // 作成する議事録ページのアドレス

	// 指定されたページを新規作成するページの準備
	payloadNewPage := map[string]string{
		"encode_hint": "ぷ",
		"plugin":      "newpage",
		"refer":       "FrontPage",
		"page":        page,
	}
	valuesNewPage := url.Values{}
	for k, v := range payloadNewPage {
		valuesNewPage.Add(k, v)
	}
	stdHeader["Content-Type"] = "application/x-www-form-urlencoded"
	body, err := digestPost(http.MethodPost, memsWiki.Scheme+"://"+memsWiki.Host, "/memswiki/index.php", stdHeader, strings.NewReader(valuesNewPage.Encode()))
	if err != nil {
		return fmt.Errorf("Failed to WriteMinutes(): %w", err)
	}
	original, err := getOriginal(body)
	if err != nil {
		return fmt.Errorf("Failed to WriteMinutes(): %w", err)
	}
	digest, err := getDigest(body)
	if err != nil {
		return fmt.Errorf("Failed to WriteMinutes(): %w", err)
	}

	log.Println(original)
	log.Println(digest)
	log.Println(msg)

	// 指定されたページに議事録を登録
	payloadPost := map[string]string{
		"encode_hint": "ぷ",
		"cmd":         "edit",
		"page":        page,
		"digest":      digest,
		"msg":         msg,
		"write":       "Update",
		"original":    original,
	}
	valuesPost := url.Values{}
	for k, v := range payloadPost {
		valuesPost.Add(k, v)
	}
	// stdHeader["Content-Type"] = "application/x-www-form-urlencoded"
	// body, err = digestPost(http.MethodPost, memsWiki.Scheme+"://"+memsWiki.Host, "/memswiki/index.php", stdHeader, strings.NewReader(valuesPost.Encode()))
	// if err != nil {
	// 	return fmt.Errorf("Failed to WriteMinutes(): %w", err)
	// }

	return nil
}
