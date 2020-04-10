#!/bin/sh

# after daily check, generating sh will work
# as well as this, send_mail sh wil be changed

dir="${HOME}/MailSystem/output/"

LAB_SCHEDULES_FILE="config/lab_schedules.csv"

EXECUTIVE_MEETINGS_FILE_EN="${dir}executive_meetings_en.csv"
EXECUTIVE_MEETINGS_FILE_JP="${dir}executive_meetings_jp.csv"
TEAMMEMS_MEETINGS_FILE_EN="${dir}teamMEMS_meetings_en.csv"
TEAMMEMS_MEETINGS_FILE_JP="${dir}teamMEMS_meetings_jp.csv"
OTHER_SCHEDULES_FILE="${dir}other_schedules.csv"
LAB_SCHEDULES_FORMATTED_FILE="${dir}lab_schedules_formatted.csv"

if [ -e ${EXECUTIVE_MEETINGS_FILE_EN} ]; then rm -rf ${EXECUTIVE_MEETINGS_FILE_EN}; fi
if [ -e ${EXECUTIVE_MEETINGS_FILE_JP} ]; then rm -rf ${EXECUTIVE_MEETINGS_FILE_JP}; fi
if [ -e ${TEAMMEMS_MEETINGS_FILE_EN} ]; then rm -rf ${TEAMMEMS_MEETINGS_FILE_EN}; fi
if [ -e ${TEAMMEMS_MEETINGS_FILE_JP} ]; then rm -rf ${TEAMMEMS_MEETINGS_FILE_JP}; fi
if [ -e ${OTHER_SCHEDULES_FILE} ]; then rm -rf ${OTHER_SCHEDULES_FILE}; fi
touch ${EXECUTIVE_MEETINGS_FILE_EN}
touch ${EXECUTIVE_MEETINGS_FILE_JP}
touch ${TEAMMEMS_MEETINGS_FILE_EN}
touch ${TEAMMEMS_MEETINGS_FILE_JP}
touch ${OTHER_SCHEDULES_FILE}

COUNT=`grep '' ${LAB_SCHEDULES_FILE} | wc -l | awk '{printf "%d", $1}'`
for i in `seq ${COUNT}`; do
  line=`cat ${LAB_SCHEDULES_FILE} | head -$i | tail -1`
  subject=`echo ${line}     | cut -d"," -f1`
  start_date=`echo ${line}  | cut -d"," -f2`
  start_time=`echo ${line}  | cut -d"," -f3`
  end_date=`echo ${line}    | cut -d"," -f4`
  end_time=`echo ${line}    | cut -d"," -f5`
  location=`echo ${line}    | cut -d ',' -f6 | sed -e "s/[^>]*\(113\|114\|128\|VDEC306\|VDEC402\|Bldg13\)[^>]*/\1/g"`

  # printed_line="${start_date} ${start_time} ${location} ${ZOOM_URL}"

  case $location in
    113 )
      formatted_location="113"
      location_en="Bldg. 3 Room 113 (Seminar 3)"
      location_jp="工学部3号館 113号室 (電気系セミナー室3)"
      ;;
    114 )
      formatted_location="114"
      location_en="Bldg. 3 Room 114 (Seminar 2)"
      location_jp="工学部3号館 114号室 (電気系セミナー室2)"
      ;;
    128 )
      formatted_location="128"
      location_en="Bldg. 3 Room 128 (Seminar 1)"
      location_jp="工学部3号館128号室 (電気系セミナー室1)"
      ;;
    VDEC306 )
      formatted_location="VDEC306"
      location_en="VDEC 306"
      location_jp="VDEC 306"
      ;;
    VDEC402 )
      formatted_location="VDEC402"
      location_en="VDEC 402"
      location_jp="VDEC 402"
      ;;
    Bldg13 )
      formatted_location="Bldg13"
      location_en="Bldg. 13"
      location_jp="13号館一般実験室"
      ;;
    * )
      formatted_location="${location}"
      location_en=$location
      location_jp=$location
      ;;
  esac

  if [ "${subject}" = "Subject" ]; then
    echo "This is a header..."
    echo $line >> $EXECUTIVE_MEETINGS_FILE_EN
    echo $line >> $EXECUTIVE_MEETINGS_FILE_JP
    echo $line >> $TEAMMEMS_MEETINGS_FILE_EN
    echo $line >> $TEAMMEMS_MEETINGS_FILE_JP
    echo $line >> $OTHER_SCHEDULES_FILE
    echo $line >> $LAB_SCHEDULES_FORMATTED_FILE
    continue;
  fi

  start_date=`date -d "${start_date}" +%Y/%m/%d`
  end_date=`date -d "${end_date}" +%Y/%m/%d`
  start_time=`date -d "${start_time}" +%H:%M`
  end_time=`date -d "${end_time}" +%H:%M`
  
  echo "${subject},${start_date},${start_time},${end_date},${end_time},${formatted_location}" >> $LAB_SCHEDULES_FORMATTED_FILE
  if [ "${subject}" = "Executive meeting" ]; then
    echo "generating minutes of Executive meeting on ${start_date}..."
    echo "${subject},${start_date},${start_time},${end_date},${end_time},${location_en}" >> $EXECUTIVE_MEETINGS_FILE_EN
    echo "${subject},${start_date},${start_time},${end_date},${end_time},${location_jp}" >> $EXECUTIVE_MEETINGS_FILE_JP
  elif [ "${subject}" = "TeamMEMS meeting" ]; then
    echo "generating minutes of TeamMEMS meeting on ${start_date}..."
    echo "${subject},${start_date},${start_time},${end_date},${end_time},${location_en}" >> $TEAMMEMS_MEETINGS_FILE_EN
    echo "${subject},${start_date},${start_time},${end_date},${end_time},${location_jp}" >> $TEAMMEMS_MEETINGS_FILE_JP
  else
    echo "not working..."
    echo "${subject},${start_date},${start_time},${end_date},${end_time},${location_jp}" >> $OTHER_SCHEDULES_FILE
  fi
  echo "subject: ${subject}" | sed "s/^/  /g"
done
