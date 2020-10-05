# MitaLab Mail System

- ラボのメールの自動送信のリポジトリです．

## ファイル構成

- (・)は Git 管理外

```
.
├── README.md
├── config
│   ├── (executive_mail.csv)
│   ├── (executive_mail_zoom.csv)
│   ├── (teamMEMS_mail.csv)
│   ├── (teamMEMS_mail_zoom.csv)
│   ├── (.private_info)
│   └──
├── Executive
│   ├── README.txt
│   ├── send_mail_to_executive.csh
│   ├── send_mail_to_executive.sh
│   ├── public_holidays.sh
│   ├── (public_holidays.txt)
│   ├── (.private_info): 環境設定ファイル
│   ├── (signature.txt)
│   ├── (schedule.txt): 予定ファイル
│   ├── (log.txt): ログ
│   └── (tmp.txt)
├── TeamMEMS
│   ├── README.txt
│   ├── send_mail_to_teamMEMS.sh
│   ├── public_holidays.sh
│   ├── (public_holidays.txt)
│   ├── (.private_info): 環境設定ファイル
│   ├── (log.txt): ログ
│   └── (tmp.txt)
└── scraping_issues.sh
```

### 各ファイルの使用方法

- 主に Exetutive に関して説明する

#### Executive

- send_mail_to_executive.sh: 実行ファイル
  - 弄る必要は無い
- public_holidays.sh: 祝日取得用のスクリプト
  - `send_mail_to_executive.sh` に呼ばれている
- (public_holidays.txt)
  - `public_holidays.sh` に自動生成される祝日データであり，無くても問題無い
- (.private_info): 環境設定ファイル
  - 個人情報を含む為，これを読み込んで実行ファイルは動く
  - **自分の名前や学年を書く**
  - メールの送信先などもここで弄る
- (signature.txt): 自分の署名ファイル
  - **自分の署名にそのまま書き換える**
  - 本文の区切り文字などは不要
- (executive_mail.csv): ミーティング日程設定ファイル
  - **ミーティング日程が決まる度に書き換える**
    - 例えば次のようなデータが、それぞれの日程について改行区切りで入っている
      - e.g.) 2020/12/01,10:00,Zoom,Executive meeting
    - 各ミーティング日程は，半角スペース区切りで `yy/mm/dd,開始時刻,場所,内容` と入力
      - `yy/mm/dd` : 日付．月と日それぞれ 2 文字ずつ(0 埋め)入力
        - e.g.) 2020/03/07、2020/12/05、2020/10/14
      - `開始時刻` : そのまま:区切りで入力
        - e.g. 10:00、09:30
      - `場所` : 基本的には以下の候補いずれかを「そのまま」コピー
        - `113` `114` `128` `VDEC306` `VDEC402` `Bldg13`
          - 英語と日本語それぞれの場所文字列を出す為のキーワード
          - `Zoom` の場合は自動で外部ファイル `executive_mail_zoom.csv` を読み込む
          - これ以外の文字列の場合，書いたものがそのまま日本語でも英語でも使われるが， `send_mail_to_executive.sh` の該当箇所を加筆すれば同様に使用可能
      - `内容` : 内容
        - `Executive meeting` など
- README.txt: 初代説明書
- (log.txt): ログ
- (tmp.txt): メールの文面として自動で作成され，自動で削除される

## 使用手順

### セットアップ

1. 下記の何れかで mozart の自分のフォルダ(`$HOME`)で `git clone`

```
git clone https://github.com/miya789/MailSystem.git # デフォルト
git clone git@github.com:miya789/MailSystem.git     # GitHubとSSH通信可能な人用
```

2. 各ディレクトリ( `Executive` , `TeamMEMS` )の `.private_info` や `signature.txt` を自分用の設定に置き換える
3. `cd Executive`
4. `chmod a+x ExecutiveMail.sh`
5. `pwd`
6. pwd で出てきた結果をコピー
7. `crontab -e`
8. i キーを押す(insert モード)
9. `00 09 * * * (pwdの結果)/send_mail_to_executive.sh`　と入力
   - e.g.) `00 09 * * * $HOME/MailSystem/Executive/send_mail_to_executive.sh`
10. `:wq` と入力(エディタから抜ける)

### 普段

- Meeting で日程が決まる度に， `{Executive, TeamMEMS}/schedule.txt` を更新
  - 別のプログラムで作成した以下のような `executive_mail.csv` と `executive_mail_zoom.csv` を読み込む

```csv:executive_mail.csv
Start Date,Start Time,Location,Subject
2021/03/01,10:00,114,Executive meeting
2021/03/08,10:00,Zoom,Executive meeting
```

```csv:executive_mail_zoom.csv
Start Date,Start Time,URL,Password
2021/03/08,10:00,https://zoom.us/j/<ID>?pwd=<ENCODED_PASS>,<PASS>
```

- 決まった直後に行うと，忘れずにできて良い
- 尚，送信に失敗した場合は自分宛に cron からメールが届く
  - その際は自分でメールを書いて送信すれば良い

### 引継ぎ後

- 下記手順で cron の設定を消去する
  1. `crontab -e`
  2. 設定時に書いた一行を削除
