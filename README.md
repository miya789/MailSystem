# MitaLab Mail System
- ラボのメールの自動送信のリポジトリです．

## ファイル構成
- (・)はGit管理外
```
.
├── README.md
├── Executive
│   ├── (.private_info): 環境設定ファイル
│   ├── README.txt
│   ├── ExecutiveMail.sh
│   ├── ExecutiveMail_sh.sh
│   ├── holidays.sh
│   ├── (holidays.txt)
│   ├── (log.txt): ログ
│   ├── (schedule.txt): 予定ファイル
│   ├── (signature.txt)
│   └── (tmp.txt)
└── TeamMEMS
    ├── (.private_info): 環境設定ファイル
    ├── README.txt
    ├── TeamMEMSmail.sh
    ├── holidays.sh
    ├── (holidays.txt)
    ├── (log.txt): ログ
    ├── (schedule.txt): 予定ファイル
    ├── (signature.txt)
    └── (tmp.txt)
```

### 各ファイルの使用方法
#### Executive
- ExecutiveMailディレクトリ
  - ExecutiveMail.sh	(本体)
    - このスクリプトの最初の部分を自分の設定用に変更
  - signature.txt		(自分の署名)
    - 自分の署名にそのまま置き換える(署名と本文の区切り文字不要)
  - schedule.txt		(ミーティング日程設定ファイル)
    - これをミーティング日程が決まるたびに置き換える。
    - 例えば次のようなデータが、それぞれの日程について改行区切りで入っている
    - e.g.) 02/09 10:00 114
    - 各日程については、半角スペース区切りで「mm/dd 開始時刻 場所」と入力
      - mm/dd: 日付。月と日それぞれ2文字ずつ(0埋め)入力 e.g.)03/07、12/05、10/14
      - 開始時刻: そのまま:区切りで入力。2文字ずつ埋める必要はかならずしもない。
      - e.g. 10:00、9:30
    - 場所：基本的には以下の候補いずれかを「そのまま」コピー
      - 113 114 128 VDEC306 Bldg13
      - これ以外の文字列の場合は、ここに書いたものがそのまま日本語でも英語でも使われる。
      - Executivemail.sh を読めばわかるが、英語と日本語それぞれの場所文字列を
  - 出すためのキーワードとなっている。他の場所が必要な場合は、該当箇所を追加すればよい。
  - holidays.sh
    - 祝日取得用のスクリプト。
  - README.txt		(このファイル)
  - (holidays.txt)
    - holidays.shで自動生成される祝日データ。コピー時はなくて良い。

## 使用手順
### セットアップ
1. 下記の何れかでmozartの自分のフォルダ(`$HOME`)で `git clone`
     - `git clone https://github.com/miya789/MailSystem.git` : デフォルト
     - `git clone git@github.com:miya789/MailSystem.git` : GitHubとSSH通信可能な人用
2. 各ディレクトリ( `Executive` , `TeamMEMS` )の `.private_info` を自分用の設定に置き換える
3. `cd Executive`
4. `chmod a+x ExecutiveMail.sh`
5. `pwd`
6. pwdで出てきた結果をコピー
7. `crontab -e`
8.  iキーを押す(insertモード)
9.  `00 09 * * * (pwdの結果)/ExecutiveMail.sh`　と入力
     eg) `00 09 * * * /$USER/Executive/ExecutiveMail.sh`
10. `:wq` と入力(エディタから抜ける)

### 普段
- Meetingで日程が決まるたびに、{Executive, TeamMEMS}/schedule.txtを更新する。
- 決まった直後に行うのが忘れなくて良い。
- なお、うまく送信できなかった時は自分宛にcronからメールが届くはずである。
  - その時は自分でメールを書いて送信すれば良い。

### 引継ぎ後
- 下記手順でcronの設定を消去する
  1. `crontab -e`
  2. 設定時に書いた一行をそのまま消す。
