#!/bin/sh -f

#-----------------------------------------
# change here
#-----------------------------------------
# cshなので <NAME>=<value> で設定

# dir="/home/miyazawa/ExecutiveMail/"
# dir="${HOME}/ExecutiveMail/"
dir="./"

echo $GRADE
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

# 日付計算
echo "#########################" >> ${LOG_FILE}
echo `date "+%Y/%m/%d-%H:%M:%S"` >> ${LOG_FILE}
# echo `date -u +"%Y-%m-%dT%H:%M:%SZ"` >> ${LOG_FILE}
echo "#########################" >> ${LOG_FILE}

# 休日データのロード
${HOLIDAYS_SCRIPT_FILE} > ${HOLIDAYS_FILE}
echo "Holiday File Regenerated." >> ${LOG_FILE}
echo "" >> ${LOG_FILE}

# 曜日の判定

youbi_judge () {
  echo $1
  plusdate=0
  yobi=`date "+%u"`
  Sat=1
  Sun=2
  tmpdate=`date "+%Y%m%d"`
  holidayflg=`grep ${tmpdate} ${HOLIDAYS_FILE}`

  echo "First... yobi:${yobi}, tmpdate:${tmpdate}, holidayflg:${holidayflg}" | column -t >> $1

  echo [ $yobi -eq $yobi ]
  if [ $yobi -eq $Sat ] || [ ${yobi} -eq $Sun ] || [ "${holidayflg}" != "" ]; then
    echo "Today is holiday." >> $1
    exit 0
  else
    plusdate=$(expr $plusdate + 1)
    yobi=`date -d "${plusdate} day" +%u`
    tmpdate=`date -d "${plusdate} days" +%Y%m%d`
    holidayflg=`grep ${tmpdate} ${HOLIDAYS_FILE}`

    echo "${plusdate}days... yobi:${yobi}, tmpdate:${tmpdate}, holidayflg:${holidayflg}" | column -t >> $1

    while [ ${yobi} -eq ${Sat} ] || [ ${yobi} -eq $Sun ] || [ "${holidayflg}" != ""  ]; do
      plusdate=$(expr $plusdate + 1)
      yobi=`date -d "${plusdate} day" +%u`
      tmpdate=`date -d "${plusdate} days" +%Y%m%d`
      holidayflg=`grep ${tmpdate} ${HOLIDAYS_FILE}`
      echo "${plusdate}days... yobi:${yobi}, tmpdate:${tmpdate}, holidayflg:${holidayflg}" | column -t >> $1
    done
  fi

  echo "Last... yobi:${yobi}, tmpdate:${tmpdate}, holidayflg:${holidayflg}" | column -t >> $1
  echo "" >> $1

  tomorrow=`date -d "${plusdate} days" "+%m/%d"`
  # tomorrow=`date -v+${plusdate}d "+%m/%d"`

  echo "Next weekday is ${tomorrow}" >> $1
}

youbi_judge ${LOG_FILE}

# メール送信判定

flg=0

COUNT=`grep '' ${SCHEDULE_FILE} | wc -l`

i=1

while [ $i -le $COUNT ]; do
	line="`cat $SCHEDULE_FILE | head -$i | tail -1`"
	DATE=`echo "$line" | cut -d' ' -f1`
	if [ $DATE = $tomorrow ]; then
		jikan=`echo "$line" | cut -d' ' -f2`
		place=`echo "$line" | cut -d' ' -f3`
		url=`echo "$line" | cut -d' ' -f4`
		echo "We have a meeting on ${tomorrow} ${jikan} at ${place}." >> ${LOG_FILE}
		flg=1
	fi
  i=$(expr $i + 1)
done

if [ $flg -eq 0 ]; then
	echo "There is no meeting on ${tomorrow}." >> ${LOG_FILE}
	exit 0
fi

# 文面作成用変数設定
placeEN=$place
placeJP=$place

case $place in
	113 )
		placeEN="Bldg. 3 Room 113 (Seminar 3)"
		placeJP="工学部3号館 113号室 (電気系セミナー室3) "
		;;
	114 )
		placeEN="Bldg. 3 Room 114 (Seminar 2)"
		placeJP="工学部3号館 114号室 (電気系セミナー室2) "
		;;
	128 )
		placeEN="Bldg. 3 Room 128 (Seminar 1)"
		placeJP="工学部3号館128号室 (電気系セミナー室1) "
		;;
	VDEC306 )
		placeEN="VDEC 306"
		placeJP="VDEC 306"
		;;
	VDEC402 )
		placeEN="VDEC 402"
		placeJP="VDEC 402"
    ;;
	Bldg13 )
		placeEN="Bldg. 13"
		placeJP="13号館一般実験室"
		;;
	default )
		placeEN=$place
		placeJP=$place
		echo "Unusual place: ${place}" >> ${LOG_FILE}
		;;
esac

case ${yobi} in
	1 )
		yobiJP="月"
		yobiEN="Mon"
    ;;
	2 )
		yobiJP="火"
		yobiEN="Tue"
    ;;
	3 )
		yobiJP="水"
		yobiEN="Wed"
    ;;
	4 )
		yobiJP="木"
		yobiEN="Thu"
    ;;
	5 )
		yobiJP="金"
		yobiEN="Fri"
    ;;
	6 )
		yobiJP="土"
		yobiEN="Sat"
    ;;
	7 )
		yobiJP="日"
		yobiEN="Sun"
    ;;
esac

month=`date -d "${plusdate} days" "+%m" | bc`
# month=`date -v+${plusdate}d "+%m" | bc`
hizuke=`date -d "${plusdate} days" "+%d" | bc`
# hizuke=`date -v+${plusdate}d "+%d" | bc`

tomorrowTitle="${month}/${hizuke}(${yobiEN})"
tomorrowJP="${month}/${hizuke}(${yobiJP})"
# tomorrowEN=`date -v+${plusdate}d "+%A, %B "`${hizuke}
tomorrowEN=`date -d "${plusdate} days" "+%A, %B "`${hizuke}

subject="The next Executive Meeting【${tomorrowTitle} ${jikan} - @${placeJP}】"
subjectEnc=`echo ${subject} | nkf --mime --ic=UTF-8 --oc=UTF-8`

# 文面作成及び送信--
if [ -e ${TMP} ]; then
	rm -rf ${TMP}
fi
touch ${TMP}
(
  echo "From: ${from}"
  echo "To: ${to}"
  # echo "Bcc: ${bcc}"
  echo "Subject: ${subjectEnc}"
  echo "Content-Type: text/plain; charset=UTF-8"
  echo "Content-Transfer-Encoding: 8bit"
  echo "MIME-Version: 1.0"
  echo 
  echo "${GRADE}の${NAME}です。"
  echo "次回のExecutive Meetingは${tomorrowJP} ${jikan} - @${placeJP}で行われます。"
  echo "Zoom URL: ${url}."
  echo "よろしくお願いします。"
  echo
  echo "Dear Executive members:"
  echo
  echo "The next Executive Meeting is going to be held at the ${placeEN} from ${jikan} on"
  echo "${tomorrowEN}. Please attend the meeting."
  echo "Thank you."
  echo
  echo '--'
  cat ${SIGNATURE_FILE}
) >> ${TMP}

#BCC使わなければこっちが安全
#cat ${TMP} | $pathsendmail -i -f ${from} ${to}

# cat ${TMP} | $pathsendmail -i -t

echo "Contents:" >> ${LOG_FILE}
cat ${TMP} | sed "s/^/  /g" >> ${LOG_FILE}
echo "" >> ${LOG_FILE}

# rm -f ${TMP}
