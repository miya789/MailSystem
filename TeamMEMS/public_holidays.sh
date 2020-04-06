#! /bin/sh

# このURLは
# Googleカレンダーの「カレンダー設定」→「日本の祝日」→「ICAL」から取得可能
# (2017/01/12現在)
url='https://calendar.google.com/calendar/ical/ja.japanese%23holiday%40group.v.calendar.google.com/public/basic.ics'

curl -s "$url"                           |
# /usr/local/bin/curl -s "$url"            |
sed -n '/^BEGIN:VEVENT/,/^END:VEVENT/p'  |
awk '/^BEGIN:VEVENT/{                    # iCalendar(RFC 5545)形式から
       rec++;                            # 日付と名称だけ抽出
     }                                   #
     match($0,/^DTSTART.*DATE:/){        # DTSTART行は日付であるから
       print rec,1,substr($0,RLENGTH+1); # 「レコード番号 "1" 日付」に
     }                                   #
     match($0,/^SUMMARY:/){              # SUMMARY行は名称であるから
       s=substr($0,RLENGTH+1);           # 「レコード番号 "2" 名称」に
       gsub(/ /,"_",s);                  #
       print rec,2,s;                    #
     }'                                  |
sort -k1n,1 -k2n,2                       | # レコード番号>列種別 にソート
awk '$2==1{printf("%d ",$3);}            # # 1レコード1行にする
     $2==2{print $3;       }             #
     '                                   |
sort                                     # 日付順にソートして出力
