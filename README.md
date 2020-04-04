# MitaLab Mail System
- by @miya_789 (EEIC2017)

ファイル構成

```
.
├── README.md
├── Executive
│   ├── .private_info: 環境設定ファイル(Git管理外)
│   ├── README.txt
│   ├── ExecutiveMail.sh
│   ├── ExecutiveMail_sh.sh
│   ├── holidays.sh
│   ├── log.txt: ログ(Git管理外)
│   └── signature.txt
└── TeamMEMS
    ├── .private_info: 環境設定ファイル(Git管理外)
    ├── README.txt
    ├── TeamMEMSmail.sh
    ├── holidays.sh
    ├── log.txt: ログ(Git管理外)
    └── signature.txt

```

以下転記
```
######################################################################
Executive Meeting 自動リマインダ
By. 竹城 雄大(EEIC2014/EEIS2016)
######################################################################

-0: はじめに
このスクリプトは、Executive Meetingのリマインダを毎度毎度忘れずに流すのが面倒な
ため、これを自動化するために作製されたものである。
以下の通りにセットアップすれば良い。

-1: ファイル構成
-ExecutiveMailディレクトリ

--ExecutiveMail.sh	(本体)
---このスクリプトの最初の部分を自分の設定用に変更

--signature.txt		(自分の署名)
---自分の署名にそのまま置き換える(署名と本文の区切り文字不要)

--schedule.txt		(ミーティング日程設定ファイル)
---これをミーティング日程が決まるたびに置き換える。
---例えば次のようなデータが、それぞれの日程について改行区切りで入っている
---e.g.) 02/09 10:00 114
---各日程については、半角スペース区切りで「mm/dd 開始時刻 場所」と入力
----mm/dd: 日付。月と日それぞれ2文字ずつ(0埋め)入力 e.g.)03/07、12/05、10/14
----開始時刻: そのまま:区切りで入力。2文字ずつ埋める必要はかならずしもない。
    e.g. 10:00、9:30
----場所：基本的には以下の候補いずれかを「そのまま」コピー
         113 114 128 VDEC306 Bldg13
-----これ以外の文字列の場合は、ここに書いたものがそのまま日本語でも英語でも使わ
     れる。
-----Executivemail.sh を読めばわかるが、英語と日本語それぞれの場所文字列を
     出すためのキーワードとなっている。他の場所が必要な場合は、該当箇所を追加
     すればよい。

--holidays.sh
---祝日取得用のスクリプト。

--README.txt		(このファイル)

--(holidays.txt)
---holidays.shで自動生成される祝日データ。コピー時はなくて良い。

-2: セットアップ
--0. ExecutiveMail.shとsignature.txtを自分用の設定に置き換える
--1. ExecutiveMailディレクトリをmozartの自分のホームフォルダにコピー
--2. >cd ExecutiveMail
--3. >chmod a+x ExecutiveMail.sh
--4. >pwd 
--5. pwdで出てきた結果をコピー
--6. >crontab -e 
--7. iキーを押す(insertモード)
--8. 00 09 * * * (pwdの結果)/ExecutiveMail.sh　と入力
     eg) 00 09 * * * /home/takeshiro/ExecutiveMail/ExecutiveMail.sh
--9. :wq と入力(エディタから抜ける)

-3: 日常的に行うこと
Executive Meetingで日程が決まるたびに、schedule.txtを更新する。
決まった直後に行うのが忘れなくて良い。
なお、うまく送信できなかった時は自分宛にcronからメールが届くはずである。
その時は自分でメールを書いて送信すれば良い。

-3: 引き継ぎ後
--cronの設定を消去する。具体的には
---1. >crontab -e
---2. 設定時に書いた一行をそのまま消す。
```