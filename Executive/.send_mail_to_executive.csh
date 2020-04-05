#!/bin/csh -f

setenv LC_ALL "c date" # for Japanese env

# set dir="/home/miyazawa/ExecutiveMail/"
# set dir="${HOME}/ExecutiveMail/"
set dir="./"

echo "Loading environment file..."
source "${dir}.private_info"

# ファイル名の準備
set TMP_FILENAME                    = "tmp.txt"
set SCHEDULE_FILENAME               = "schedule.txt"
set SIGNATURE_FILENAME              = "signature.txt"
set PUBLIC_HOLIDAYS_FILENAME        = "public_holidays.txt"
set PUBLIC_HOLIDAYS_SCRIPT_FILENAME = "public_holidays.sh"
set LOG_FILENAME                    = "log.txt"
set TMP                         = "${dir}${TMP_FILENAME}"
set SCHEDULE_FILE               = "${dir}${SCHEDULE_FILENAME}"
set SIGNATURE_FILE              = "${dir}${SIGNATURE_FILENAME}"
set PUBLIC_HOLIDAYS_FILE        = "${dir}${PUBLIC_HOLIDAYS_FILENAME}"
set PUBLIC_HOLIDAYS_SCRIPT_FILE = "${dir}${PUBLIC_HOLIDAYS_SCRIPT_FILENAME}"
set LOG_FILE                    = "${dir}${LOG_FILENAME}"
set SENDMAIL_PATH = "/usr/sbin/sendmail"

# 日付計算
echo "" >> ${LOG_FILE}
echo "[MAIL LOG] `date "+%Y/%m/%d-%H:%M:%S"`" >> ${LOG_FILE}

# 休日データのロード
${PUBLIC_HOLIDAYS_SCRIPT_FILE} > ${PUBLIC_HOLIDAYS_FILE}
echo "Holiday File Regenerated." | sed "s/^/  /g" >> ${LOG_FILE}
echo "" >> ${LOG_FILE}

# 曜日の判定
@ plusdate=0
set day_of_week_num = `date "+%u"`
set date            = `date "+%Y%m%d"`
set is_holiday      = `grep ${date} ${PUBLIC_HOLIDAYS_FILE}`

echo "Today:" | sed "s/^/  /g" | column -t -s, >> ${LOG_FILE}
echo "day of week(No.): ${day_of_week_num}, date: ${date}, is_holiday: ${is_holiday}" | sed "s/^/    /g" >> ${LOG_FILE}

if ((${day_of_week_num} == 6) || (${day_of_week_num} == 7) || ("${is_holiday}" != "")) then
  echo "Today is a holiday, so finished." | sed "s/^/  /g" >> ${LOG_FILE}
  exit 0
else
  echo "Today is not a holiday, so continuing..." | sed "s/^/  /g" >> ${LOG_FILE}
  echo "" >> ${LOG_FILE}

  # 次の平日の調査
  echo "Searching the next weekday..." | sed "s/^/  /g" >> ${LOG_FILE}
  @ plusdate++
  # day_of_week_num >>>
  if ( "${OSTYPE}" == "FreeBSD" ) then
    set day_of_week_num=`date -v+${plusdate}d "+%u"`
  else if ( "${OSTYPE}" == "linux-gnu" ) then
    set day_of_week_num=`date -d "${plusdate} day" +%u`
  endif
  # day_of_week_num <<<
  # date >>>
  if ( "${OSTYPE}" == "FreeBSD" ) then
    set date=`date -v+${plusdate}d "+%Y%m%d"`
  else if ( "${OSTYPE}" == "linux-gnu" ) then
    set date=`date -d "${plusdate} days" +%Y%m%d`
  endif
  # date <<<
  set is_holiday=`grep ${date} ${PUBLIC_HOLIDAYS_FILE}`
  echo "${plusdate} day later:" | sed "s/^/  /g" | column -t -s, >> ${LOG_FILE}
  echo "day of week(No.): ${day_of_week_num}, date: ${date}, is_holiday: ${is_holiday}" | sed "s/^/    /g" >> ${LOG_FILE}

  # 次の平日に辿り着く迄逃れられない！
  while((${day_of_week_num} == 6) || (${day_of_week_num} == 7) || ("${is_holiday}" != ""))
    @ plusdate++
    # day_of_week_num >>>
    if ( "${OSTYPE}" == "FreeBSD" ) then
      set day_of_week_num=`date -v+${plusdate}d "+%u"`
    else if ( "${OSTYPE}" == "linux-gnu" ) then
      set day_of_week_num=`date -d "${plusdate} day" +%u`
    endif
    # day_of_week_num <<<
    # date >>>
    if ( "${OSTYPE}" == "FreeBSD" ) then
      set date=`date -v+${plusdate}d "+%Y%m%d"`
    else if ( "${OSTYPE}" == "linux-gnu" ) then
      set date=`date -d "${plusdate} days" +%Y%m%d`
    endif
    # date <<<
    set is_holiday=`grep ${date} ${PUBLIC_HOLIDAYS_FILE}`
    echo "${plusdate} days later:" | sed "s/^/  /g" | column -t -s, >> ${LOG_FILE}
    echo "day of week(No.): ${day_of_week_num}, date: ${date}, is_holiday: ${is_holiday}" | sed "s/^/    /g" >> ${LOG_FILE}
  end
endif

echo "Finished!" | sed "s/^/  /g" >> ${LOG_FILE}
echo "The next weekday:" | sed "s/^/  /g" | column -t -s, >> ${LOG_FILE}
echo "day of week(No.): ${day_of_week_num}, date: ${date}, is_holiday: ${is_holiday}" | sed "s/^/    /g" >> ${LOG_FILE}

if ( "${OSTYPE}" == "FreeBSD" ) then
  set NEXT_WEEKDAY=`date -v+${plusdate}d "+%m/%d"`
else if ( "${OSTYPE}" == "linux-gnu" ) then
  set NEXT_WEEKDAY=`date -d "${plusdate} days" "+%m/%d"`
endif

echo "" >> ${LOG_FILE}
echo "Checking if there is the meeting on ${NEXT_WEEKDAY}..." | sed "s/^/  /g" >> ${LOG_FILE}
# メール送信判定
@ should_send_mail=0
set COUNT=`grep '' ${SCHEDULE_FILE} | wc -l`
@ i=1
while ( $i <= $COUNT )
  set line=`cat $SCHEDULE_FILE | head -$i | tail -1`
  echo "[${i}/${COUNT}]: ${line}" | sed "s/^/    /g" >> ${LOG_FILE}
  set DATE=`echo "$line" | cut -d' ' -f1`
  if ( $DATE == $NEXT_WEEKDAY ) then
    set MEETING_TIME=`echo "$line" | cut -d' ' -f2`
    set MEETING_PLACE=`echo "$line" | cut -d' ' -f3`
    set MEETING_ZOOM_URL=`echo "$line" | cut -d' ' -f4`
    echo "We have the meeting from ${MEETING_TIME} on ${NEXT_WEEKDAY} at ${MEETING_PLACE}." | sed "s/^/  /g" >> ${LOG_FILE}
    @ should_send_mail=1
  endif
  @ i++
end
if ( $should_send_mail == 0 ) then
  echo "There is no meeting on ${NEXT_WEEKDAY}." >> ${LOG_FILE}
  exit 0
endif

# 場所の表記変換
switch($MEETING_PLACE)
  case 113 :
    set MEETING_PLACE_EN="Bldg. 3 Room 113 (Seminar 3)"
    set MEETING_PLACE_JP="工学部3号館 113号室 (電気系セミナー室3) "
    breaksw
  case 114 :
    set MEETING_PLACE_EN="Bldg. 3 Room 114 (Seminar 2)"
    set MEETING_PLACE_JP="工学部3号館 114号室 (電気系セミナー室2) "
    breaksw
  case 128 :
    set MEETING_PLACE_EN="Bldg. 3 Room 128 (Seminar 1)"
    set MEETING_PLACE_JP="工学部3号館128号室 (電気系セミナー室1) "
    breaksw
  case VDEC306 :
    set MEETING_PLACE_EN="VDEC 306"
    set MEETING_PLACE_JP="VDEC 306"
    breaksw
  case VDEC402 :
    set MEETING_PLACE_EN="VDEC 402"
    set MEETING_PLACE_JP="VDEC 402"
  case Bldg13 :
    set MEETING_PLACE_EN="Bldg. 13"
    set MEETING_PLACE_JP="13号館一般実験室"
    breaksw
  default :
    set MEETING_PLACE_EN=$MEETING_PLACE
    set MEETING_PLACE_JP=$MEETING_PLACE
    echo "Unusual place: ${MEETING_PLACE}" >> ${LOG_FILE}
    breaksw
endsw

# 曜日の表記変換
switch(${day_of_week_num})
  case 1 :
    set day_of_week_JP="月"
    set day_of_week_EN="Mon"
    breaksw
  case 2 :
    set day_of_week_JP="火"
    set day_of_week_EN="Tue"
       breaksw
  case 3 :
    set day_of_week_JP="水"
    set day_of_week_EN="Wed"
       breaksw
  case 4 :
    set day_of_week_JP="木"
    set day_of_week_EN="Thu"
       breaksw
  case 5 :
    set day_of_week_JP="金"
    set day_of_week_EN="Fri"
       breaksw
  case 6 :
    set day_of_week_JP="土"
    set day_of_week_EN="Sat"
       breaksw
  case 7 :
    set day_of_week_JP="日"
    set day_of_week_EN="Sun"
       breaksw
endsw

if ( "${OSTYPE}" == "FreeBSD" ) then
  set MONTH=`date -v+${plusdate}d "+%m" | bc`
  set DAY=`date -v+${plusdate}d "+%d" | bc`
else if ( "${OSTYPE}" == "linux-gnu" ) then
  set MONTH=`date -d "${plusdate} days" "+%m" | bc`
  set DAY=`date -d "${plusdate} days" "+%d" | bc`
endif

set DATE_FOR_TITLE="${MONTH}/${DAY}(${day_of_week_EN})"
set DATE_FOR_CONTENTS_JP="${MONTH}/${DAY}(${day_of_week_JP})"
if ( "${OSTYPE}" == "FreeBSD" ) then
  set DATE_FOR_CONTENTS_EN=`date -v+${plusdate}d "+%A, %B "`${DAY}
else if ( "${OSTYPE}" == "linux-gnu" ) then
  set DATE_FOR_CONTENTS_EN=`date -d "${plusdate} days" "+%A, %B "`${DAY}
endif

set SUBJECT="The next Executive Meeting【${DATE_FOR_TITLE} ${MEETING_TIME} - @${MEETING_PLACE_JP}】"
set SUBJECT_ENC=`echo ${SUBJECT} | nkf --mime --ic=UTF-8 --oc=UTF-8`

# メール文面ファイル(temp.txt)執筆
if (-e ${TMP}) then
  rm -rf ${TMP}
endif
touch ${TMP}

echo "From: ${from}" >> ${TMP}
echo "To: ${to}" >> ${TMP}
# echo "Bcc: ${bcc}" >> ${TMP}
echo "Subject: ${SUBJECT_ENC}" >> ${TMP}
echo "Content-Type: text/plain; charset=UTF-8" >> ${TMP}
echo "Content-Transfer-Encoding: 8bit" >> ${TMP}
echo "MIME-Version: 1.0" >> ${TMP}
echo  >> ${TMP}

echo "Executiveの皆様" >> ${TMP}
echo "" >> ${TMP}
echo "${GRADE}の${NAME_JP}です．" >> ${TMP}
echo "次回のExecutive Meetingは${DATE_FOR_CONTENTS_JP} ${MEETING_TIME} - @${MEETING_PLACE_JP}で行われます．" >> ${TMP}
if ("$MEETING_ZOOM_URL" != "") then
  echo "(Zoom URL: ${MEETING_ZOOM_URL})" >> ${TMP}
endif
echo "宜しくお願い致します．" >> ${TMP}
echo "" >> ${TMP}
echo "" >> ${TMP}
echo "Dear Executive members," >> ${TMP}
echo "" >> ${TMP}
echo "I'm ${GRADE} ${NAME_EN}." >> ${TMP}
echo "The next Executive Meeting is going to be held at the ${MEETING_PLACE_EN} from ${MEETING_TIME} on ${DATE_FOR_CONTENTS_EN}." >> ${TMP}
if ("$MEETING_ZOOM_URL" != "") then
  echo "(Zoom URL: ${MEETING_ZOOM_URL})" >> ${TMP}
endif
echo "Please attend the meeting." >> ${TMP}
echo "Thank you." >> ${TMP}
echo "" >> ${TMP}
echo "--" >> ${TMP}
cat ${SIGNATURE_FILE} >> ${TMP}

# メール文面の送信
# cat ${TMP} | $SENDMAIL_PATH -i -f ${from} ${to} # BCC使わなければこっちが安全
cat ${TMP} | $SENDMAIL_PATH -i -t

# メール文面のログ吐き出し
echo "" >> ${LOG_FILE}
echo "The sent mail is as follows..." | sed "s/^/  /g" >> ${LOG_FILE}
cat ${TMP} | sed "s/^/    /g" >> ${LOG_FILE}
echo "" >> ${LOG_FILE}

# メール文面ファイルの削除
rm -f ${TMP}
