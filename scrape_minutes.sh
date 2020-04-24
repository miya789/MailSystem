#!/bin/sh -f

CONFIG_DIR="config"
TMP_DIR="tmp"

LOG_FILE="${TMP_DIR}/log.txt"
printf "[SCRAPE LOG] `date "+%Y/%m/%d-%H:%M:%S"`\n" >> ${LOG_FILE}
printf "  Loading env file...\n\n" >> ${LOG_FILE}

. "${CONFIG_DIR}/.env"

####################################################################################################################
#                                       ここで，schdule.csvを元に，議事録を作るべきdateを作成                          #
####################################################################################################################

SCHEDULE_FILE="${CONFIG_DIR}/lab_schedules.csv"
COUNT=`grep '' $SCHEDULE_FILE | wc -l | awk '{printf "%d", $1}'`
TODAY=`date +%Y/%m/%d`
printf "Checking if there is the meeting on today (${TODAY})...\n" | sed "s/^/  /g"  >> ${LOG_FILE}
should_generate_minute=0
i=2 # To skip header
while [ $i -le $COUNT ] && [ $should_generate_minute -eq 0 ]; do
  line=`cat $SCHEDULE_FILE | head -$i | tail -1`
  printf "[${i}/${COUNT}]: ${line}\n" | sed "s/^/    /g" >> ${LOG_FILE}
  _date=`echo "${line}" | cut -d',' -f2`
  DATE=`date -d "${_date}" +%Y/%m/%d`
  if [ "$DATE" = "$TODAY" ]; then
    MEETING_SUBJECT=`echo "$line" | cut -d',' -f1`
    MEETING_TIME=`echo "$line" | cut -d',' -f3`
    MEETING_PLACE=`echo "$line" | cut -d',' -f6`
    MEETING_ZOOM_URL=`echo "$line" | cut -d',' -f7`
    printf "We have the meeting from ${MEETING_TIME} on ${DATE} at ${MEETING_PLACE}.\n\n" | sed "s/^/  /g" >> ${LOG_FILE}
    should_generate_minute=1
  fi
  i=$(expr $i + 1)
done

# 予定の有無を判定
if [ $should_generate_minute -eq 0 ]; then
  printf "There is no meeting on ${TODAY}.\n\n" >> ${LOG_FILE}
  exit 0
fi

# Preparation tmp html files
TOP_HTML="${TMP_DIR}/top.html"
INDEX_HTML="${TMP_DIR}/index.html"
INDEX_EDIT_HTML="${TMP_DIR}/index_edit.html"
TARGET_HTML="${TMP_DIR}/target.html"
TARGET_EDIT_HTML="${TMP_DIR}/target_edit.html"
INDEX_ORIGINAL_TXT="${TMP_DIR}/index_original.txt"
index_msg_txt="${TMP_DIR}/index_msg.txt"
TARGET_ORIGINAL_TXT="${TMP_DIR}/target_original.txt"
target_msg_txt="${TMP_DIR}/target_msg.txt"

TODAY_FOR_MINUTES=`date +%Y%m%d`

CURL_OPTIONS="--socks5 ${PROXY} --digest -u ${USER}:${PASSWORD}"

# Scraping
echo "  Getting top_page.html..." >> ${LOG_FILE}
curl ${MINUTES_TOP_URL} ${CURL_OPTIONS} > ${TOP_HTML}
sleep 5
if [ "${MEETING_SUBJECT}" = "TeamMEMS meeting" ]; then
  INDEX_URL=`grep 'meeting 議事録</a>' ${TOP_HTML} | sed -n 's/^.*href="\([^"]*\)".*$/\1/p' | head -1 | tail -1`
elif [ "${MEETING_SUBJECT}" = "Executive meeting" ]; then
  INDEX_URL=`grep 'meeting 議事録</a>' ${TOP_HTML} | sed -n 's/^.*href="\([^"]*\)".*$/\1/p' | head -2 | tail -1`
fi
(
  echo "    INDEX_URL: ${INDEX_URL}"
  echo "  Finished!"
  echo
) >> ${LOG_FILE}

echo "  Getting index.html..." >> ${LOG_FILE}
curl ${INDEX_URL} ${CURL_OPTIONS} > ${INDEX_HTML}
sleep 5
TARGET_URL=`grep "${TODAY_FOR_MINUTES}" ${INDEX_HTML} | sed -n 's/^.*href=\"\([^"]*\)".*$/\1/p' | sed -e 's/\&amp\;/\&/g' | tail -1`

(
  echo "    TARGET_URL:${TARGET_URL}"
  echo "  Finished!"
  echo
) >> ${LOG_FILE}

while [ "${TARGET_URL}" = "" ]; do
  echo "  Cannot find today minutes, so generating page..." >> ${LOG_FILE}
  INDEX_EDIT_URL=`grep "Edit" ${INDEX_HTML} | sed -n 's/^.*href=\"\([^"]*\)".*$/\1/p' | sed -e 's/\&amp\;/\&/g' | head -1 | tail -1`
  (
    echo "    INDEX_EDIT_URL:${INDEX_EDIT_URL}"
    echo
    echo "  Getting index_edit.html..."
    echo
  ) >> ${LOG_FILE}
  curl ${INDEX_EDIT_URL} ${CURL_OPTIONS} > ${INDEX_EDIT_HTML}
  sleep 5



  # params作成現場
  INDEX_DIGEST=`cat "${INDEX_EDIT_HTML}" | grep digest | sed -n 's/^.* value=\"\([^"]*\).*/\1/p'`
  cat ${INDEX_EDIT_HTML} | sed -ne '/<textarea name=\"original/,/<\/textarea>/p' | sed 's/  \(<[^>]*\)/\1/g' | sed -e 's/<[^>]*>//g' | sed -e 's/\&amp\;/\&/g' | sed -e 's/\&gt\;/\>/g' > ${INDEX_ORIGINAL_TXT}

  cp ${INDEX_ORIGINAL_TXT} ${index_msg_txt}
  INSERTING_TXT="-[[${TODAY}>ミーティング議事録/${TODAY_FOR_MINUTES}]]" # 新しい議事録のURLなどを挿入する行
  sed -i -e "6a ${INSERTING_TXT}" ${index_msg_txt} | sed "s/^/    /g"

  INDEX_MSG_ENC=`cat ${index_msg_txt} | nkf -WwMQ | sed -e ':loop; N; $!b loop; s/=\n//g' | sed -z 's/\n/%0D%0A/g' | tr = % | tr -d '\n'`
  INDEX_ORIGINAL_ENC=`cat ${INDEX_ORIGINAL_TXT} | nkf -WwMQ | sed -e ':loop; N; $!b loop; s/=\n//g' | sed -z 's/\n/%0D%0A/g' | tr = % | tr -d '\n'`

  # 各パラメータの値の確認が必要
  ## encode_hit="ぷ"
  ## cmd="edit"
  ## digest="<毎回異なるハッシュ値らしきもの>"
  ## msg="<投稿内容>"
  ## original="<元の内容>"
  ## write="Update"
  if [ "${INDEX_MSG_ENC}" = "" ] || [ "${INDEX_ORIGINAL_ENC}" = "" ]; then
    printf "[ERROR] `date "+%Y/%m/%d-%H:%M:%S"`\n" >> ${LOG_FILE}
    printf "  Generating INDEX_PARAMS was failed!\n\n" >> ${LOG_FILE}
    exit 1
  fi
  INDEX_PARAMS="encode_hint=%E3%81%B7&cmd=edit&digest=${INDEX_DIGEST}&msg=${INDEX_MSG_ENC}&original=${INDEX_ORIGINAL_ENC}&write=Update"

  echo $INDEX_PARAMS > "${TMP_DIR}/`date "+%Y%m%d_%H%M%S"`_params.txt"
  cat $index_msg_txt > "${TMP_DIR}/`date "+%Y%m%d_%H%M%S"`_index_msg.txt"
  cat $INDEX_ORIGINAL_TXT > "${TMP_DIR}/`date "+%Y%m%d_%H%M%S"`_index_original.txt"



  curl ${INDEX_EDIT_URL} ${CURL_OPTIONS} -XPOST -d "${INDEX_PARAMS}"
  sleep 5

  curl ${INDEX_URL} ${CURL_OPTIONS} > ${INDEX_HTML}
  TARGET_URL=`grep "${TODAY_FOR_MINUTES}" ${INDEX_HTML} | sed -n 's/^.*href=\"\([^"]*\)".*$/\1/p' | sed -e 's/\&amp\;/\&/g' | tail -1`
done
(
  echo "  Preparaing page finished!"
  echo "  Got TARGET_URL: ${TARGET_URL}"
  echo
) >> ${LOG_FILE}

echo "  Accessing target page..." >> ${LOG_FILE}
curl ${TARGET_URL} ${CURL_OPTIONS} > ${TARGET_HTML}

TARGET_EDIT_URL=`grep "Edit" ${TARGET_HTML} | sed -n 's/^.*href=\"\([^"]*\)".*$/\1/p' | sed -e 's/\&amp\;/\&/g' | head -1 | tail -1`
if [ "${TARGET_EDIT_URL}" = "" ]; then
  echo "    There is not minute page."
  TARGET_EDIT_URL=TARGET_URL
  cp ${TARGET_HTML} ${TARGET_EDIT_HTML}
else
  echo "    Already there is minute page."
  curl ${TARGET_EDIT_URL} ${CURL_OPTIONS} > ${TARGET_EDIT_HTML}
fi
sleep 5

# params作成現場
TARGET_DIGEST=`cat "${TARGET_EDIT_HTML}" | grep digest | sed -n 's/^.* value=\"\([^"]*\).*/\1/p'`
# cat ${TARGET_EDIT_HTML} | sed -ne '/<textarea name=\"original/,/<\/textarea>/p' | sed 's/  \(<[^>]*\)/\1/g' | sed -e 's/<[^>]*>//g' | sed '1,2d' > ${TARGET_ORIGINAL_TXT} # 初めの二行はテンプレート
cat ${TARGET_EDIT_HTML} | sed -ne '/<textarea name=\"original/,/<\/textarea>/p' | sed 's/  \(<[^>]*\)/\1/g' | sed -e 's/<[^>]*>//g' > ${TARGET_ORIGINAL_TXT}

cp ${TARGET_ORIGINAL_TXT} ${target_msg_txt}
# INSERTING_TXT="Test\nThis is the sentence." # 書く内容を用意
# sed -i "1s/^/${INSERTING_TXT}\n/" ${target_msg_txt}

TARGET_MSG_ENC=`cat ${target_msg_txt} | nkf -WwMQ | sed -e ':loop; N; $!b loop; s/=\n//g' | sed -z 's/\n/%0D%0A/g' | sed -z 's/=2A/*/g' | sed -z 's/=2B/+/g' | sed -z 's/=25/%/g' | tr = % | tr -d '\n'`
TARGET_ORIGINAL_ENC=`cat ${TARGET_ORIGINAL_TXT} | nkf -WwMQ | sed -e ':loop; N; $!b loop; s/=\n//g' | sed -z 's/\n/%0D%0A/g' | sed -z 's/=2A/*/g' | sed -z 's/=2B/+/g' | sed -z 's/=25/%/g' | tr = % | tr -d '\n'`

# 各パラメータの値の確認が必要
## encode_hit="ぷ"
## cmd="edit"
## digest="<毎回異なるハッシュ値らしきもの>"
## msg="<投稿内容>"
## original="<元の内容>"
## write="Update"
TARGET_PARAMS="encode_hint=%E3%81%B7&digest=${TARGET_DIGEST}&msg=${TARGET_MSG_ENC}&original=${TARGET_ORIGINAL_ENC}&write=Update"
echo "    TARGET_PARAMS:${TARGET_PARAMS}"

curl ${TARGET_EDIT_URL} ${CURL_OPTIONS} -XPOST -d "${TARGET_PARAMS}" > result.txt
sleep 5
