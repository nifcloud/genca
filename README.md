## プライベートCA 作成ツール genca

ニフクラ [リモートアクセスVPNゲートウェイ](https://pfs.nifcloud.com/service/ra_vpngw.htm)で利用できるサーバー証明書、CA証明書作成ツールです。  

## 注意事項

ツールに関するお問い合わせは[ベーシックサポート（トラブル窓口）](https://pfs.nifcloud.com/inquiry/support.htm)のサポート範囲外となります。  
ツールに関するお問い合わせはgithub上のissueを起票してください。  
コミュニティベースでのサポートとなります。  

## 利用方法

### ツールの実行

  ご利用のOS環境にあったバイナリをダウンロードし、gencaを実行して下さい。  
  実行すると、プライベートCAと署名済みクライアント証明書、サーバー証明書が自動生成されます。  

* nifcloud.local.CAcert.pem
  * プライベートCA証明書です。
  * コントロールパネルのCA証明書にアップロードし、リモートアクセスVPNゲートウェイに設定する事で、クライアント証明書認証を利用する事が出来ます。
* nifcloud.local.CAkey.pem
  * プライベートCAの秘密鍵です。
* client.nifcloud.local.crt.pem
  * クライアント証明書(署名前)
* client.nifcloud.local.pem
  * クライアント証明書の秘密鍵
* client.nifcloud.local.csr.pem
  * クライアント証明書署名要求
* client.nifcloud.local.signed.crt.pem
  * 署名済みクライアント証明書
  * プライベートCAによって電子署名されたクライアント証明書です。
* server.nifcloud.local.crt.pem
  * サーバー証明書(署名前)
* server.nifcloud.local.csr.pem
  * サーバー証明書署名要求
* server.nifcloud.local.pem
  * サーバー証明書の秘密鍵
* server.nifcloud.local.signed.crt.pem
  * 署名済みサーバー証明書
  * プライベートCAによって電子署名されたサーバー証明書です。

### サーバー証明書

生成されたサーバー証明書(server.nifcloud.local.signed.crt.pem)と、サーバー証明書の秘密鍵(server.nifcloud.local.pem)をコントロールパネルからアップロードし、  
サーバー証明書の一覧に表示される事を確認します。  
[クラウドヘルプ（サーバー証明書：アップロード）](https://pfs.nifcloud.com/help/ssl/upload.htm)

このサーバー証明書をリモートアクセスVPNゲートウェイの作成時や設定変更で指定します。  
接続するには、クライアント設定ファイルにプライベートCA証明書(nifcloud.local.CAcert.pem)を指定する必要があります。  
詳細は以下ご参照ください。
* [クラウド技術仕様（リモートアクセスVPNゲートウェイ:クライアント設定ファイル）](https://pfs.nifcloud.com/spec/ra_vpngw/client_config.htm)
* [クラウドユーザーガイド（リモートアクセスVPNゲートウェイ：Windows 11のリモートアクセスVPNゲートウェイ(v2.0.0)利用方法）](https://pfs.nifcloud.com/guide/cp/ra_vpngw/ravgwv2_setup_win11.htm)

### CA証明書のアップロード

生成されたCA証明書(nifcloud.local.CAcert.pem)をコントロールパネルからアップロードし、CA証明書の一覧に表示される事を確認します。  
[クラウドヘルプ（CA証明書：アップロード）](https://pfs.nifcloud.com/help/ca/upload.htm)



## ソースコードのビルド

    go build genca.go 
    # または
    gox
    
    # 実行
    ./genca

## ライセンス

[LICENSE](./LICENSE)をご参照下さい。



