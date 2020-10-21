#!/bin/sh

CONFIG_DIR="config"
TMP_DIR="tmp"

LOG_FILE="${TMP_DIR}/log.txt"
printf "[SCRAPE_ISSUES LOG] `date "+%Y/%m/%d-%H:%M:%S"`\n" >> ${LOG_FILE}
printf "  Loading env file...\n\n" >> ${LOG_FILE}

. "${CONFIG_DIR}/.env"


NANOTECH_HELP_TXT="${TMP_DIR}/nanotech_help.txt"
RECEPTION_TXT="${TMP_DIR}/reception.txt"
CURL_OPTIONS="--socks5 ${PROXY} --digest -u ${USER}:${PASSWORD} -v"
WAIT_TIME=3

# Scraping
echo "  Getting issues..." >> ${LOG_FILE}

data_file="${TMP_DIR}/data.txt";

if [ -e ${data_file} ]; then
  rm -rf ${data_file}
fi
touch ${data_file}

(
  echo "[[Executive Meeting]]"
  echo ""
  echo "#contents"
  echo ""
  echo "*2020/4/16 10:00- @セミナー室2"
  echo "- 出席"
  echo "--"
  echo ""
  echo "**共有事項"
  echo ""
  echo "**ナノテクヘルプ"
) >> ${data_file}

curl ${NANOTECH_HELP_URL} ${CURL_OPTIONS} > "${NANOTECH_HELP_TXT}"
curl ${NANOTECH_HELP_URL2} ${CURL_OPTIONS} >> "${NANOTECH_HELP_TXT}"
curl ${RECEPTION_URL} ${CURL_OPTIONS} > "${RECEPTION_TXT}"

NANOTECH_HELP_COUNT=`grep '<td class=\"subject\"><a href=\"\/issues' ${NANOTECH_HELP_TXT} | wc -l | awk '{printf "%d", $1}'`
i=$NANOTECH_HELP_COUNT
while [ $i -ne 0 ]; do
  sed -n 's/<td class=\"subject\"><a href=\"\/issues\/\(.*\)\">\(.*\)<\/a><\/td>$/- [[#\1 \2>http:\/\/mozart.if.t.u-tokyo.ac.jp:3000\/issues\/\1]]/p' ${NANOTECH_HELP_TXT} | tail -${i} | head -1 | sed 's/^    //g' >> ${data_file}
  i=$(expr $i - 1)
done

RECEPTION_COUNT=`grep '<td class=\"subject\"><a href=\"\/issues' ${RECEPTION_TXT} | wc -l | awk '{printf "%d", $1}'`
i=$RECEPTION_COUNT
while [ $i -ne 0 ]; do
  sed -n 's/<td class=\"subject\"><a href=\"\/issues\/\(.*\)\">\(.*\)<\/a><\/td>$/- [[#\1 \2>http:\/\/mozart.if.t.u-tokyo.ac.jp:3000\/issues\/\1]]/p' ${RECEPTION_TXT} | tail -${i} | head -1 | sed 's/^    //g' >> ${data_file}
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
  echo ",2020/05/07(木), 10:00 ~ 12:00,セミナー室2"
) >> ${data_file}
