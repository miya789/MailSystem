#!/bin/sh

CONFIG_DIR="config"
TMP_DIR="tmp"
if [ ! -d ${TMP_DIR} ]; then
  mkdir ${TMP_DIR}
fi

LOG_FILE="${TMP_DIR}/log.txt"
if [ ! -e ${LOG_FILE} ]; then
  touch ${LOG_FILE}
fi

printf "[SCRAPE_ISSUES LOG] `date "+%Y/%m/%d-%H:%M:%S"`\n" >> ${LOG_FILE}
printf "  Loading env file...\n\n" >> ${LOG_FILE}

. "${CONFIG_DIR}/.env"

usage_exit() {
  echo "Usage: $0 [-p] ..." 1>&2
  exit 1
}
PROXY_OPT=""
while getopts ph FLAG; do
  case $FLAG in
    p ) PROXY_OPT="--socks5 ${PROXY}"
      ;;
    h ) usage_exit
      ;;
    \? ) usage_exit
      ;;
  esac
done
shift $((OPTIND - 1))

NANOTECH_HELP_TXT="${TMP_DIR}/nanotech_help.txt"
RECEPTION_TXT="${TMP_DIR}/reception.txt"
CURL_OPTIONS="--digest -u ${USER}:${PASSWORD} -v ${PROXY_OPT}"
WAIT_TIME=3

# Scraping
echo "  Getting issues...\n" >> ${LOG_FILE}

data_file="${TMP_DIR}/data.txt";

if [ -e ${data_file} ]; then
  rm -rf ${data_file}
fi
touch ${data_file}

# Calculate date
# send_mail_to_executive.sh から拝借
SCHEDULE_FILENAME="config/executive_mail.csv"
SCHEDULE_FILE="${dir}${SCHEDULE_FILENAME}"

TODAY=`date +%Y/%m/%d`

printf "Checking if there is the meeting on ${TODAY}...\n" | sed "s/^/  /g"  >> ${LOG_FILE}
should_scrape_issues=0
COUNT=`grep '' ${SCHEDULE_FILE} | wc -l | awk '{printf "%d", $1}'`
schedule_i=2
while [ $schedule_i -le $COUNT ] && [ $should_scrape_issues -eq 0 ]; do
  line=`cat $SCHEDULE_FILE | head -$schedule_i | tail -1`
  printf "[${schedule_i}/${COUNT}]: ${line}\n" | sed "s/^/    /g" >> ${LOG_FILE}
  DATE=`echo "$line" | cut -d',' -f1`
  if [ $DATE = $TODAY ]; then
    MEETING_TIME=`echo "$line" | cut -d',' -f2`
    MEETING_PLACE=`echo "$line" | cut -d',' -f3`
    printf "We have the meeting from ${MEETING_TIME} on ${DATE} at ${MEETING_PLACE}.\n" | sed "s/^/  /g" >> ${LOG_FILE}
    should_scrape_issues=1
  fi
  schedule_i=$(expr $schedule_i + 1)
done

if [ $should_scrape_issues -eq 0 ]; then
  printf "There is no meeting on today (${TODAY}).\n\n" | sed "s/^/  /g" >> ${LOG_FILE}
  exit 0
fi

(
  echo "[[Executive Meeting]]"
  echo ""
  echo "#contents"
  echo ""
  echo "*${DATE} ${MEETING_TIME}- @${MEETING_PLACE}"
  echo "- 出席"
  echo "--"
  echo ""
  echo "**共有事項"
  echo ""
  echo "**ナノテクヘルプ"
) >> ${data_file}

curl ${RECEPTION_URL} ${CURL_OPTIONS} > "${RECEPTION_TXT}"
curl ${NANOTECH_HELP_URL} ${CURL_OPTIONS} > "${NANOTECH_HELP_TXT}"
curl ${NANOTECH_HELP_URL2} ${CURL_OPTIONS} >> "${NANOTECH_HELP_TXT}"

RECEPTION_COUNT=`grep '<td class=\"subject\"><a href=\"\/issues' ${RECEPTION_TXT} | wc -l | awk '{printf "%d", $1}'`
i=$RECEPTION_COUNT
while [ $i -ne 0 ]; do
  sed -n 's/<td class=\"subject\"><a href=\"\/issues\/\(.*\)\">\(.*\)<\/a><\/td>$/- [[#\1 \2>http:\/\/mozart.if.t.u-tokyo.ac.jp:3000\/issues\/\1]]/p' ${RECEPTION_TXT} | tail -${i} | head -1 | sed 's/^    //g' >> ${data_file}
  i=$(expr $i - 1)
done

echo "" >> ${data_file}

NANOTECH_HELP_COUNT=`grep '<td class=\"subject\"><a href=\"\/issues' ${NANOTECH_HELP_TXT} | wc -l | awk '{printf "%d", $1}'`
i=$NANOTECH_HELP_COUNT
while [ $i -ne 0 ]; do
  sed -n 's/<td class=\"subject\"><a href=\"\/issues\/\(.*\)\">\(.*\)<\/a><\/td>$/- [[#\1 \2>http:\/\/mozart.if.t.u-tokyo.ac.jp:3000\/issues\/\1]]/p' ${NANOTECH_HELP_TXT} | tail -${i} | head -1 | sed 's/^    //g' >> ${data_file}
  i=$(expr $i - 1)
done

(
  echo "**装置メンテ"
  echo ""
  echo "**設備"
  echo ""
  echo "**講習会"
  echo ""
  echo "**その他"
  echo ""
  echo "**今後の予定"
  echo ""
  echo "***Executive Meetings"
  echo ",Date & Time,Location,Contents"
) >> ${data_file}

# Calculate date
# send_mail_to_executive.sh から拝借
# schedule_i を使い回してる

CALENDAR_FILENAME="config/executive_calendar.csv"
CALENDAR_FILE="${dir}${CALENDAR_FILENAME}"

while [ $schedule_i -le $COUNT ]; do
  line=`cat $CALENDAR_FILE | head -$schedule_i | tail -1`
  printf "[${schedule_i}/${COUNT}]: ${line}\n" | sed "s/^/    /g" >> ${LOG_FILE}
  DATE=`echo "$line" | cut -d',' -f1`
  MEETING_START_TIME=`echo "$line" | cut -d',' -f2`
  MEETING_END_TIME=`echo "$line" | cut -d',' -f4`
  MEETING_CONTENTS=`echo "$line" | cut -d',' -f6`
  day_of_week=`date -d ${DATE} "+%a"`
  echo ",${DATE}(${day_of_week}),${MEETING_START_TIME} ~ ${MEETING_END_TIME},${MEETING_CONTENTS}" | sed "s/^/      /g" >> ${LOG_FILE}
  echo ",${DATE}(${day_of_week}),${MEETING_START_TIME} ~ ${MEETING_END_TIME},${MEETING_CONTENTS}" >> ${data_file}
  schedule_i=$(expr $schedule_i + 1)
done
