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

curl --socks5 $PROXY $MINUTES_TOP_URL --digest -u $USER:$PASSWORD > top_page.html
sleep 1
TeamMEMS=`grep 'meeting 議事録</a>' top_page.html | sed -n 's/^.*href="\([^"]*\)".*$/\1/p' | head -1 | tail -1`
Executive=`grep 'meeting 議事録</a>' top_page.html | sed -n 's/^.*href="\([^"]*\)".*$/\1/p' | head -2 | tail -1`

echo "TeamMEMS: ${TeamMEMS}"
echo "Executive: ${Executive}"

curl --socks5 $PROXY $TeamMEMS --digest -u $USER:$PASSWORD > teamMEMS.html

date=`date +%Y/%m/%d` # 危険なので暫定的に

next_url=`grep "${date}" teamMEMS.html | sed -n 's/^.*href=\"\([^"]*\)".*$/\1/p' | tail -1`
sleep 1

curl --socks5 $PROXY $next_url --digest -u $USER:$PASSWORD > teamMEMS_contents.html
edit_url=`grep "Edit" teamMEMS_contents.html | sed -n 's/^.*href=\"\([^"]*\)".*$/\1/p' | sed -e 's/\&amp\;/\&/g' | head -1 | tail -1`
sleep 1

echo $test_edit
curl --socks5 $PROXY $edit_url --digest -u $USER:$PASSWORD > teamMEMS_edit.html

digest=`cat teamMEMS_edit.html | grep digest | sed -n 's/^.* value=\"\([^"]*\).*/\1/p'`
cat teamMEMS_edit.html | sed -ne '/<textarea name=\"original/,/<\/textarea>/p' | sed 's/  \(<[^>]*\)/\1/g' | sed -e 's/<[^>]*>//g' | sed '1,2d' > original.txt # 初めの二行はテンプレート

echo digest
echo $digest
echo original
cat original.txt

####################################################################################################################
#                                       ここで，tmp.txtを元に，msg.txtを作成                                         #
####################################################################################################################

msg_enc=`cat msg.txt | nkf -WwMQ | sed -e ':loop; N; $!b loop; s/=\n//g' | sed -z 's/\n/%0D%0A/g' | tr = % | tr -d '\n'`
original_enc=`cat original.txt | nkf -WwMQ | sed -e ':loop; N; $!b loop; s/=\n//g' | sed -z 's/\n/%0D%0A/g' | tr = % | tr -d '\n'`

# echo "${original_enc}" | nkf -W --url-input | sed -z 's/%0A/\n/g' > tmp2.txt
# echo $original_enc
# cat tmp2.txt

# 各パラメータの値の確認が必要
params="encode_hint=%E3%81%B7"\
"&cmd=edit"\
"&digest=${digest}"\
"&msg=${msg_enc}"\
"&original=${original_enc}"\
"&write=Update"
echo params
# "&template_page="\

echo ${params} > params.txt

curl --socks5 $PROXY $edit_url --digest -u $USER:$PASSWORD -XPOST -d "$params" > teamMEMS_test_edit.html
