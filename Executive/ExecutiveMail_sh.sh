#!/bin/sh -f

export LC_ALL="c date" # for Japanese env

# dir="/home/miyazawa/ExecutiveMail/"
# dir="${HOME}/ExecutiveMail/"
dir="./"

echo "Loading environment file..."
. "${dir}.private_info_sh"

# ファイル名の準備
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
pathsendmail = "/usr/sbin/sendmail"

generate_diff_option () {
  if [ "${OSTYPE}" = "FreeBSD" ]; then
    echo "-v+${plusdate}d"
  elif [ "${OSTYPE}" = "linux-gnu" ]; then
    echo "-d \"${plusdate} days\""
  fi
}

# 日付計算
echo "" >> ${LOG_FILE}
echo "[MAIL LOG] `date "+%Y/%m/%d-%H:%M:%S"`" >> ${LOG_FILE}

# 休日データのロード
${HOLIDAYS_SCRIPT_FILE} > ${HOLIDAYS_FILE}
echo "Holiday File Regenerated." | sed "s/^/  /g" >> ${LOG_FILE}
echo "" >> ${LOG_FILE}

# 曜日の判定
plusdate=0
day_of_week_num=`date "+%u"`
Sat=1
Sun=2
date=`date "+%Y%m%d"`
is_holiday=`grep ${date} ${HOLIDAYS_FILE}`

echo "Today:" | sed "s/^/  /g" | column -t -s, >> ${LOG_FILE}
echo "day of week(No.): ${day_of_week_num}, date: ${date}, is_holiday: ${is_holiday}" | sed "s/^/    /g" >> ${LOG_FILE}

if [ $day_of_week_num -eq $Sat ] || [ $day_of_week_num -eq $Sun ] || [ "${is_holiday}" != "" ]; then
  echo "Today is a holiday, so finished." | sed "s/^/  /g" >> ${LOG_FILE}
  exit 0
else
  echo "Today is not a holiday, so continuing..." | sed "s/^/  /g" >> ${LOG_FILE}
  echo "" >> ${LOG_FILE}

  # 次の平日の調査
  echo "Searching the next weekday..." | sed "s/^/  /g" >> ${LOG_FILE}
  plusdate=$(expr $plusdate + 1)
  day_of_week_num=`eval "date $(generate_diff_option ${plusdate}) +%u"`
  date=`eval "date $(generate_diff_option ${plusdate}) +%Y%m%d"`
  is_holiday=`grep ${date} ${HOLIDAYS_FILE}`
  echo "${plusdate} day later:" | sed "s/^/  /g" | column -t -s, >> ${LOG_FILE}
  echo "day of week(No.): ${day_of_week_num}, date: ${date}, is_holiday: ${is_holiday}" | sed "s/^/    /g" >> ${LOG_FILE}

  # 次の平日に辿り着く迄逃れられない！
  while [ $day_of_week_num -eq ${Sat} ] || [ $day_of_week_num -eq $Sun ] || [ "${holidayflg}" != ""  ]; do
    plusdate=$(expr $plusdate + 1)
    day_of_week_num=`eval "date $(generate_diff_option ${plusdate}) +%u"`
    date=`eval "date $(generate_diff_option ${plusdate}) +%Y%m%d"`
    is_holiday=`grep ${date} ${HOLIDAYS_FILE}`
    echo "${plusdate} days later:" | sed "s/^/  /g" | column -t -s, >> ${LOG_FILE}
    echo "day of week(No.): ${day_of_week_num}, date: ${date}, is_holiday: ${is_holiday}" | sed "s/^/    /g" >> ${LOG_FILE}
  done
fi

echo "Finished!" | sed "s/^/  /g" >> ${LOG_FILE}
echo "The next weekday:" | sed "s/^/  /g" | column -t -s, >> ${LOG_FILE}
echo "day of week(No.): ${day_of_week_num}, date: ${date}, is_holiday: ${is_holiday}" | sed "s/^/    /g" >> ${LOG_FILE}

NEXT_WEEKDAY=`eval "date $(generate_diff_option ${plusdate}) +%m/%d"`

echo "" >> ${LOG_FILE}
echo "Checking if there is the meeting on ${NEXT_WEEKDAY}..." | sed "s/^/  /g" >> ${LOG_FILE}
# メール送信判定
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
    echo "We have the meeting from ${MEETING_TIME} on ${NEXT_WEEKDAY} at ${MEETING_PLACE}." | sed "s/^/  /g" >> ${LOG_FILE}
    should_send_mail=1
  fi
  i=$(expr $i + 1)
done

if [ $should_send_mail -eq 0 ]; then
	echo "There is no meeting on ${NEXT_WEEKDAY}." >> ${LOG_FILE}
	exit 0
fi

# 文面作成用変数設定
MEETING_PLACE_EN=$MEETING_PLACE
MEETING_PLACE_JP=$MEETING_PLACE

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
	default )
		MEETING_PLACE_EN=$MEETING_PLACE
		MEETING_PLACE_JP=$MEETING_PLACE
		echo "Unusual place: ${MEETING_PLACE}" >> ${LOG_FILE}
		;;
esac

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
esac

MONTH=`eval "date $(generate_diff_option ${plusdate}) +%m" | bc`
DAY=`eval "date $(generate_diff_option ${plusdate}) +%d" | bc`

DATE_FOR_TITLE="${MONTH}/${DAY}(${day_of_week_EN})"
DATE_FOR_CONTENTS_JP="${MONTH}/${DAY}(${day_of_week_JP})"
DATE_FOR_CONTENTS_EN=`"date $(generate_diff_option ${plusdate}) \"+%A, %B \""`${DAY}

SUBJECT="The next Executive Meeting【${DATE_FOR_TITLE} ${MEETING_TIME} - @${MEETING_PLACE_JP}】"
SUBJECT_ENC=`echo ${SUBJECT} | nkf --mime --ic=UTF-8 --oc=UTF-8`

# メール文面ファイル(temp.txt)執筆
if [ -e ${TMP} ]; then
  rm -rf ${TMP}
fi
touch ${TMP}
(
  echo "From: ${from}"
  echo "To: ${to}"
  # echo "Bcc: ${bcc}"
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

# メール文面の送信
# cat ${TMP} | $SENDMAIL_PATH -i -f ${from} ${to} # BCC使わなければこっちが安全
# cat ${TMP} | $SENDMAIL_PATH -i -t

# メール文面のログ吐き出し
echo "" >> ${LOG_FILE}
echo "The sent mail is as follows..." | sed "s/^/  /g" >> ${LOG_FILE}
cat ${TMP} | sed "s/^/    /g" >> ${LOG_FILE}
echo "" >> ${LOG_FILE}

# メール文面ファイルの削除
# rm -f ${TMP}
