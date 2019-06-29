

PostgreSQLを扱おうという人がLinuxを利用できない状況にあることは稀である気がするので、CentOS前提で話を進める。
CentOS以外を利用している人はCentOSと利用しているOSの差異を理解しているであろうと勝手に推測する。

1. どうにかしてCentOSをセットアップする

    すでにパソコンがあるならば、VirtualBoxを使うのが手軽かもしれない。

    https://www.virtualbox.org/

    ネットワークやファイヤーウォールの設定に注意。

2. PostgreSQLをインストールする。

   https://yum.postgresql.org/repopackages.php
   にアクセスしてリポジトリのアドレスを調べてから、例えば

   yum -y install https://download.postgresql.org/pub/repos/yum/reporpms/EL-7-x86_64/pgdg-redhat-repo-latest.noarch.rpm
    
   とかやってから（URLは適宜書き換えが必要）、

   yum -y install postgresql11-server

   とかやればいい。

3. 初期設定

    データベースクラスタの作成

    環境変数類の設定

    設定ファイルの編集

      自分がアクセスしたい方法でつながるようにしておいたほうがよい。

4. 起動する、自動起動設定をする

    systemctl start postgresql-11.service

    systemctl enable postgresql-11.service




WALバッファ → WALファイル → アーカイブファイル
ベースバックアップ+アーカイブファイルでPITR可能