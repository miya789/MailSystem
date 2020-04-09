#!/bin/sh -f

dir="./"

. "${dir}.env"

curl --socks5 $PROXY $MINUTES_TOP_URL --digest -u $USER:$PASSWORD > top_page.html
sleep 1
TeamMEMS=`grep 'meeting 議事録</a>' top_page.html | sed -n 's/^.*href="\([^"]*\)".*$/\1/p' | head -1 | tail -1`
Executive=`grep 'meeting 議事録</a>' top_page.html | sed -n 's/^.*href="\([^"]*\)".*$/\1/p' | head -2 | tail -1`

echo "TeamMEMS: ${TeamMEMS}"
echo "Executive: ${Executive}"

curl --socks5 $PROXY $TeamMEMS --digest -u $USER:$PASSWORD > teamMEMS.html
next_url=`grep "2020/04/16" teamMEMS.html | sed -n 's/^.*href=\"\([^"]*\)".*$/\1/p' | tail -1`
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
