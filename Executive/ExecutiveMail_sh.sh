#!/bin/sh -f

# 1.1 日本語環境で動作する場合用
export LC_ALL="c date"

# 1.2 ディレクトリ設定
# dir="${HOME}/ExecutiveMail/" # こっちの方が安全かもしれない
dir="./"

# 1.3 環境設定の読み込み
. "${dir}.private_info_sh"

# 1.4 ファイル名の準備
TMP_FILENAME="tmp.txt"
SCHEDULE_FILENAME="schedule.txt"
SIGNATURE_FILENAME="signature.txt"
HOLIDAYS_FILENAME="holidays.txt"
HOLIDAYS_SCRIPT_FILENAME="holidays.sh"
LOG_FILENAME="log2.txt"
TMP="${dir}${TMP_FILENAME}"
SCHEDULE_FILE="${dir}${SCHEDULE_FILENAME}"
SIGNATURE_FILE="${dir}${SIGNATURE_FILENAME}"
HOLIDAYS_FILE="${dir}${HOLIDAYS_FILENAME}"
HOLIDAYS_SCRIPT_FILE="${dir}${HOLIDAYS_SCRIPT_FILENAME}"
LOG_FILE=${dir}${LOG_FILENAME}
pathsendmail="/usr/sbin/sendmail"

# 1.5 OSが異なる環境でも動作確認を行う為，日付差分のoptionを生成する関数
generate_diff_option () {
  if [ "${OSTYPE}" = "FreeBSD" ]; then
    echo "-v+${plusdate}d"
  elif [ "${OSTYPE}" = "linux-gnu" ]; then
    echo "-d \"${plusdate} days\""
  fi
}

# 1.6 ログ用の時間を記録
echo "[MAIL LOG] `date "+%Y/%m/%d-%H:%M:%S"`" >> ${LOG_FILE}

# 1.7 最新休日情報のロード
${HOLIDAYS_SCRIPT_FILE} > ${HOLIDAYS_FILE}
echo "Holiday File Regenerated.\n" | sed "s/^/  /g" >> ${LOG_FILE}

# 2.1 曜日の判定
Sat=6
Sun=7
plusdate=0
day_of_week_num=`date "+%u"`
date=`date "+%Y%m%d"`
is_holiday=`grep ${date} ${HOLIDAYS_FILE}`

# 2.2 本日の詳細
(
  echo "Today:" | sed "s/^/  /g" | column -t -s,
  echo "day of week(No.): ${day_of_week_num}, date: ${date}, is_holiday: ${is_holiday}" | sed "s/^/    /g"
) >> ${LOG_FILE}

# 2.3 「本日が休日か」判定
if [ $day_of_week_num -eq $Sat ] || [ $day_of_week_num -eq $Sun ] || [ "${is_holiday}" != "" ]; then
  echo "Today is a holiday, so finished.\n" | sed "s/^/  /g" >> ${LOG_FILE}
  exit 0
else
  echo "Today is not a holiday, so continuing...\n" | sed "s/^/  /g" >> ${LOG_FILE}
fi

# 2.4 次の平日の探索
# echo "" >> ${LOG_FILE}
echo "Searching the next weekday..." | sed "s/^/  /g" >> ${LOG_FILE}
plusdate=$(expr $plusdate + 1)
day_of_week_num=`eval "date $(generate_diff_option ${plusdate}) +%u"`
date=`eval "date $(generate_diff_option ${plusdate}) +%Y%m%d"`
is_holiday=`grep ${date} ${HOLIDAYS_FILE}`
echo "${plusdate} day later:" | sed "s/^/  /g" | column -t -s, >> ${LOG_FILE}
echo "day of week(No.): ${day_of_week_num}, date: ${date}, is_holiday: ${is_holiday}" | sed "s/^/    /g" >> ${LOG_FILE}

## 次の平日に辿り着くまでループ
while [ $day_of_week_num -eq ${Sat} ] || [ $day_of_week_num -eq $Sun ] || [ "${holidayflg}" != ""  ]; do
  plusdate=$(expr $plusdate + 1)
  day_of_week_num=`eval "date $(generate_diff_option ${plusdate}) +%u"`
  date=`eval "date $(generate_diff_option ${plusdate}) +%Y%m%d"`
  is_holiday=`grep ${date} ${HOLIDAYS_FILE}`
  echo "${plusdate} days later:" | sed "s/^/  /g" | column -t -s, >> ${LOG_FILE}
  echo "day of week(No.): ${day_of_week_num}, date: ${date}, is_holiday: ${is_holiday}" | sed "s/^/    /g" >> ${LOG_FILE}
done

# 2.6 発見した次の翌日の詳細
NEXT_WEEKDAY=`eval "date $(generate_diff_option ${plusdate}) +%m/%d"`
(
  echo "Finished!" | sed "s/^/  /g"
  echo "The next weekday:" | sed "s/^/  /g" | column -t -s,
  echo "day of week(No.): ${day_of_week_num}, date: ${date}, is_holiday: ${is_holiday}" | sed "s/^/    /g"
  echo ""
) >> ${LOG_FILE}

# 3.1 該当日付の予定確認
echo "Checking if there is the meeting on ${NEXT_WEEKDAY}..." | sed "s/^/  /g"  >> ${LOG_FILE}
should_send_mail=0
COUNT=`grep '' ${SCHEDULE_FILE} | wc -l`
i=1
while [ $i -le $COUNT ]; do
  line=`cat $SCHEDULE_FILE | head -$i | tail -1`
  echo "[${i}/${COUNT}]: ${line}" | sed "s/^/    /g" >> ${LOG_FILE}
  DATE=`echo "$line" | cut -d' ' -f1`
  if [ $DATE = $NEXT_WEEKDAY ]; then
    MEETING_TIME=`echo "$line" | cut -d' ' -f2`
    MEETING_PLACE=`echo "$line" | cut -d' ' -f3`
    MEETING_ZOOM_URL=`echo "$line" | cut -d' ' -f4`
    echo "We have the meeting from ${MEETING_TIME} on ${DATE} at ${MEETING_PLACE}." | sed "s/^/  /g" >> ${LOG_FILE}
    should_send_mail=1
  fi
  i=$(expr $i + 1)
done

# 3.2 予定の有無を判定
if [ $should_send_mail -eq 0 ]; then
	echo "There is no meeting on ${DATE}.\n" >> ${LOG_FILE}
	exit 0
fi

# 3.3 場所の表記変換
case $MEETING_PLACE in
	113 )
		MEETING_PLACE_EN="Bldg. 3 Room 113 (Seminar 3)"
		MEETING_PLACE_JP="工学部3号館 113号室 (電気系セミナー室3) "
		;;
	114 )
		MEETING_PLACE_EN="Bldg. 3 Room 114 (Seminar 2)"
		MEETING_PLACE_JP="工学部3号館 114号室 (電気系セミナー室2) "
		;;
	128 )
		MEETING_PLACE_EN="Bldg. 3 Room 128 (Seminar 1)"
		MEETING_PLACE_JP="工学部3号館128号室 (電気系セミナー室1) "
		;;
	VDEC306 )
		MEETING_PLACE_EN="VDEC 306"
		MEETING_PLACE_JP="VDEC 306"
		;;
	VDEC402 )
		MEETING_PLACE_EN="VDEC 402"
		MEETING_PLACE_JP="VDEC 402"
    ;;
	Bldg13 )
		MEETING_PLACE_EN="Bldg. 13"
		MEETING_PLACE_JP="13号館一般実験室"
		;;
	* )
		MEETING_PLACE_EN=$MEETING_PLACE
		MEETING_PLACE_JP=$MEETING_PLACE
		echo "Unusual place: ${MEETING_PLACE}" | sed "s/^/  /g" >> ${LOG_FILE}
		;;
esac

# 3.4 曜日の表記変換
case ${day_of_week_num} in
	1 )
		day_of_week_JP="月"
		day_of_week_EN="Mon"
    ;;
	2 )
		day_of_week_JP="火"
		day_of_week_EN="Tue"
    ;;
	3 )
		day_of_week_JP="水"
		day_of_week_EN="Wed"
    ;;
	4 )
		day_of_week_JP="木"
		day_of_week_EN="Thu"
    ;;
	5 )
		day_of_week_JP="金"
		day_of_week_EN="Fri"
    ;;
	6 )
		day_of_week_JP="土"
		day_of_week_EN="Sat"
    ;;
	7 )
		day_of_week_JP="日"
		day_of_week_EN="Sun"
    ;;
  * )
		echo "The day of week(No.${day_of_week_num}) is invalid error." | sed "s/^/  /g" >> ${LOG_FILE}
    exit 1
		;;
esac

# 4.1 日付の表記用意
MONTH=`eval "date $(generate_diff_option ${plusdate}) +%m" | bc`
DAY=`eval "date $(generate_diff_option ${plusdate}) +%d" | bc`

DATE_FOR_TITLE="${MONTH}/${DAY}(${day_of_week_EN})"
DATE_FOR_CONTENTS_JP="${MONTH}/${DAY}(${day_of_week_JP})"
DATE_FOR_CONTENTS_EN=`eval "date "$(generate_diff_option ${plusdate})" +'%A, %B'"`${DAY}

# 4.2 件名の作成及びエンコード
SUBJECT="The next Executive Meeting【${DATE_FOR_TITLE} ${MEETING_TIME} - @${MEETING_PLACE_JP}】"
SUBJECT_ENC=`echo ${SUBJECT} | nkf --mime --ic=UTF-8 --oc=UTF-8`

# 4.3 文面ファイル(temp.txt)の用意
if [ -e ${TMP} ]; then
  rm -rf ${TMP}
fi
touch ${TMP}

# 4.4 文面ファイル(tmp.txt)の執筆
(
  echo "From: ${from}"
  echo "To: ${to}"
  if [ "$BCC" != "" ]; then
    echo "Bcc: ${BCC}"
  fi
  echo "Subject: ${SUBJECT_ENC}"
  echo "Content-Type: text/plain; charset=UTF-8"
  echo "Content-Transfer-Encoding: 8bit"
  echo "MIME-Version: 1.0"
  echo 

  echo "Executiveの皆様"
  echo ""
  echo "${GRADE}の${NAME_JP}です．"
  echo "次回のExecutive Meetingは${DATE_FOR_CONTENTS_JP} ${MEETING_TIME} - @${MEETING_PLACE_JP}で行われます．"
  if [ "$MEETING_ZOOM_URL" != "" ]; then
    echo "(Zoom URL: ${MEETING_ZOOM_URL})"
  fi
  echo "宜しくお願い致します．"
  echo ""
  echo ""
  echo "Dear Executive members,"
  echo ""
  echo "I'm ${GRADE} ${NAME_EN}."
  echo "The next Executive Meeting is going to be held at the ${MEETING_PLACE_EN} from ${MEETING_TIME} on ${DATE_FOR_CONTENTS_EN}."
  if [ "$MEETING_ZOOM_URL" != "" ]; then
    echo "(Zoom URL: ${MEETING_ZOOM_URL})"
  fi
  echo "Please attend the meeting."
  echo "Thank you."
  echo ""
  echo "--"
  cat ${SIGNATURE_FILE}
) >> ${TMP}

# 4.5 メールの送信
# cat ${TMP} | $SENDMAIL_PATH -i -f ${from} ${to} # BCC使わなければこっちが安全
cat ${TMP} | $SENDMAIL_PATH -i -t

# 4.6 文面ファイル(tmp.txt)をログへ吐き出し
(
  echo ""
  echo "The sent mail is as follows..." | sed "s/^/  /g"
  cat ${TMP} | sed "s/^/    /g"
  echo ""
) >> ${LOG_FILE}

# 4.7 文面ファイル(tmp.txt)の削除
rm -f ${TMP}
