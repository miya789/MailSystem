# MitaLab Mail System
- ラボのメールの自動送信のリポジトリです．

## ファイル構成
- (・)はGit管理外
```
.
├── README.md
├── Executive
│   ├── ExecutiveMail.sh
│   ├── public_holidays.sh
│   ├── (public_holidays.txt)
│   ├── (.private_info): 環境設定ファイル
│   ├── (signature.txt)
│   ├── (schedule.txt): 予定ファイル
│   ├── README.txt
│   ├── (log.txt): ログ
│   └── (tmp.txt)
└── TeamMEMS
    ├── TeamMEMSmail.sh
    ├── public_holidays.sh
    ├── (public_holidays.txt)
    ├── (.private_info): 環境設定ファイル
    ├── (signature.txt)
    ├── (schedule.txt): 予定ファイル
    ├── README.txt
    ├── (log.txt): ログ
    └── (tmp.txt)
```

### 各ファイルの使用方法
#### Executive
- ExecutiveMail.sh: 実行ファイル
  - 弄る必要は無い
- public_holidays.sh: 祝日取得用のスクリプト
  - `ExecutiveMail.sh` に呼ばれている
- (public_holidays.txt)
  - `public_holidays.sh` に自動生成される祝日データであり，無くても問題無い
- (.private_info): 環境設定ファイル
  - 個人情報を含む為，これを読み込んで実行ファイルは動く
  - **自分の名前や学年を書く**
  - メールの送信先などもここで弄る
- (signature.txt): 自分の署名ファイル
  - **自分の署名にそのまま書き換える**
  - 本文の区切り文字などは不要
- (schedule.txt): ミーティング日程設定ファイル
  - **ミーティング日程が決まる度に書き換える**
    - 例えば次のようなデータが、それぞれの日程について改行区切りで入っている
      - e.g.) 02/09 10:00 114 https:/zoom.us/j/<ID>?pwd=<PWD>
    - 各ミーティング日程は，半角スペース区切りで `mm/dd 開始時刻 場所 (URL)` と入力
      - `mm/dd` : 日付．月と日それぞれ2文字ずつ(0埋め)入力
        - e.g.) 03/07、12/05、10/14
      - `開始時刻` : そのまま:区切りで入力。2文字ずつ埋める必要はかならずしもない
        - e.g. 10:00、9:30
      - `場所` : 基本的には以下の候補いずれかを「そのまま」コピー
        - `113` `114` `128` `VDEC306` `VDEC402` `Bldg13`
          - 英語と日本語それぞれの場所文字列を出す為のキーワード
          - これ以外の文字列の場合，書いたものがそのまま日本語でも英語でも使われるが， `ExecutiveMail.sh` の該当箇所を加筆すれば同様に使用可能
      - `(URL)` : Webミーティング用のURL
        - 2020/4/1現在の感染症情勢により必要になったものであり，無くても動作する
- README.txt: 初代説明書
- (log.txt): ログ
- (tmp.txt): メールの文面として自動で作成され，自動で削除される

## 使用手順
### セットアップ
1. 下記の何れかでmozartの自分のフォルダ(`$HOME`)で `git clone`
```
git clone https://github.com/miya789/MailSystem.git # デフォルト
git clone git@github.com:miya789/MailSystem.git     # GitHubとSSH通信可能な人用
```
2. 各ディレクトリ( `Executive` , `TeamMEMS` )の `.private_info` や `signature.txt` を自分用の設定に置き換える
3. `cd Executive`
4. `chmod a+x ExecutiveMail.sh`
5. `pwd`
6. pwdで出てきた結果をコピー
7. `crontab -e`
8.  iキーを押す(insertモード)
9.  `00 09 * * * (pwdの結果)/ExecutiveMail.sh`　と入力
    - eg) `00 09 * * * $HOME/MailSystem/Executive/ExecutiveMail.sh`
10. `:wq` と入力(エディタから抜ける)

### 普段
- Meetingで日程が決まる度に， `{Executive, TeamMEMS}/schedule.txt` を更新
- 決まった直後に行うと，忘れずにできて良い
- 尚，送信に失敗した場合は自分宛にcronからメールが届く
  - その際は自分でメールを書いて送信すれば良い

### 引継ぎ後
- 下記手順でcronの設定を消去する
  1. `crontab -e`
  2. 設定時に書いた一行を削除
