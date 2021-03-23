package lab_cmd

import (
	"LabMeeting/pkg/lab_mail"
	"LabMeeting/pkg/meeting_type"
	"LabMeeting/pkg/memswiki"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func SendMinutes() {
	// 議事録に生成するページのアドレスを指定
	// 上書きはできない筈だが，存在するアドレスには注意すること
	var num int
	fmt.Println("\"Executive Meeting/[日付]\" の[日付]を入力して下さい．")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		n, err := strconv.Atoi(scanner.Text())
		num = n
		if err != nil || len(scanner.Text()) != 8 {
			fmt.Println("8桁の数値を入力してください．(例: 20210101)")
		} else {
			// 確認の案内メッセージ
			if scanner.Text() != time.Now().Format("20060102") {
				fmt.Printf("\"Executive Meeting/%s\" を作成します．今日の日付ではありませんがよろしいですか? [y/N]\n", scanner.Text())
			} else {
				fmt.Printf("\"Executive Meeting/%s\" を作成します．よろしいですか? [y/N]\n", scanner.Text())
			}

			for scanner.Scan() {
				if strings.TrimSpace(strings.ToLower(scanner.Text())) == "y" {
					break
				} else if strings.TrimSpace(strings.ToLower(scanner.Text())) == "n" {
					fmt.Println("スクリプトを停止します．")
					os.Exit(0)
				}
			}
			fmt.Printf("\"Executive Meeting/%d\" を作成します．\n", num)
			break

		}
	}

	// 議事録として登録するファイルの読み込み
	filePath := "../config/minutes.txt"
	fmt.Println("読み込むテキストファイルを指定してください．(Default: ../config/minutes.txt)")
	scanner = bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Text() != "" {
			filePath = scanner.Text()
			break
		} else {
			break
		}
	}
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Println(fmt.Errorf("Failed to Read(): %w", err))
		return
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
	}

	// 議事録をWikiへ登録
	msg := string(b)
	fmt.Printf("\"Executive Meeting/%d\" に書き込む内容を表示します．\n%s\n", num, msg)
	fmt.Printf("\"Executive Meeting/%d\" に以上の内容を本当に書き込んでよろしいですか? [y/N]\n", num)
	scanner = bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if strings.TrimSpace(strings.ToLower(scanner.Text())) == "y" {
			break
		} else if strings.TrimSpace(strings.ToLower(scanner.Text())) == "n" {
			fmt.Println("スクリプトを停止します．")
			os.Exit(0)
		}
	}
	executive_list := "http://mozart.if.t.u-tokyo.ac.jp/memswiki/index.php?Executive%20Meeting"
	if err := memswiki.WriteMinute(num, msg); err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("\x1b[31m今回作成した記事へのリンクを一覧ページへ追加するのは手動で行ったください．\n%s\x1b[0m\n", executive_list)

	// 議事録をメールへ送信
	fmt.Println("メールにも流しますか? [y/N]")
	scanner = bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if strings.TrimSpace(strings.ToLower(scanner.Text())) == "y" {
			break
		} else if strings.TrimSpace(strings.ToLower(scanner.Text())) == "n" {
			fmt.Println("スクリプトを停止します．")
			os.Exit(0)
		}
	}
	lab_mail.SendMinutesMail(meeting_type.Executive, strconv.Itoa(num), msg)

	return
}
