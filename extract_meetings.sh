#!/bin/sh

# after daily check, generating sh will work
# as well as this, send_mail sh wil be changed

file="config/meeting_schedule_2020_3.csv"
# file="../0_Calendar/meeting_schedule_2020.csv"

SCHEDULE_TEAMMEMS_FILE="schedule_teamMEMS.txt"
SCHEDULE_EXECUTIVE_FILE="schedule_executive.txt"

COUNT=`grep '' ${file} | wc -l | awk '{printf "%d", $1}'`

if [ -e ${SCHEDULE_EXECUTIVE_FILE} ]; then
  rm -rf ${SCHEDULE_EXECUTIVE_FILE}
fi
touch ${SCHEDULE_EXECUTIVE_FILE}

if [ -e ${SCHEDULE_TEAMMEMS_FILE} ]; then
  rm -rf ${SCHEDULE_TEAMMEMS_FILE}
fi
touch ${SCHEDULE_TEAMMEMS_FILE}

for i in `seq ${COUNT}`; do
  line=`cat ${file} | head -$i | tail -1`
  Subject=`echo ${line}   | cut -d"," -f1`
  StartDate=`echo ${line} | cut -d"," -f2`
  StartTime=`echo ${line} | cut -d"," -f3`
  EndDate=`echo ${line}   | cut -d"," -f4`
  EndTime=`echo ${line}   | cut -d"," -f5`
  Location=`echo ${line}  | cut -d"," -f6`

  # printed_line="${StartDate} ${StartTime} ${Location} ${ZOOM_URL"

  if [ "${Subject}" = "Subject" ]; then
    echo "This is a header..."
    continue;
  elif [ "${Subject}" = "Executive meeting" ]; then
    echo "generating minutes of Executive meeting on ${StartDate}..."
    StartDate=`date -d "${StartDate}" +%Y/%m/%d`
    echo "${StartDate},${StartTime},${Location}" >> $SCHEDULE_EXECUTIVE_FILE
  elif [ "${Subject}" = "TeamMEMS meeting" ]; then
    echo "generating minutes of TeamMEMS meeting on ${StartDate}..."
    StartDate=`date -d "${StartDate}" +%Y/%m/%d`
    echo "${StartDate},${StartTime},${Location}" >> $SCHEDULE_TEAMMEMS_FILE
  else
    echo "not working..."
  fi
  echo "Subject: ${Subject}" | sed "s/^/  /g"
done
