## プライベートCA 作成ツール genca

ニフクラ [リモートアクセスVPNゲートウェイ](https://pfs.nifcloud.com/service/ra_vpngw.htm)のクライアント証明書認証で利用出来る、プライベートCA作成ツールです。  
[プライベートCAを1clickで作成出来るツール genca を公開しました！](https://blog.pfs.nifcloud.com/20190410_genca)

## 注意事項

  * 本ツールはニフクラのサポート対象外となります。ご利用は自己責任でお願いいたします。

## 利用方法

### ツールの実行

  ご利用のOS環境にあったバイナリをダウンロードし、gencaを実行して下さい。  
  実行すると、プライベートCAと署名済みクライアント証明書が自動生成されます。  

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
* client.nifcloud.local.pfx
  * クライアント証明書(PKCS #12形式)
  * クライアント証明書認証を利用する際に必要なファイルです。Windowsなどのクライアント端末に予めインポートする必要があります。

### CA証明書のアップロード

生成されたCA証明書(nifcloud.local.CAcert.pem)をコントロールパネルからアップロードし、CA証明書の一覧に表示される事を確認します。  
[クラウドヘルプ（CA証明書：アップロード）](https://pfs.nifcloud.com/help/ca/upload.htm)

### クライアント証明書(PKCS #12形式)のインポート

クライアント証明書認証を利用するには、クライアント証明書(client.nifcloud.local.pfx)を端末にインポートする必要があります。  
インポート時にパスワードが求められますが、本ツールで作成したクライアント証明書にはパスワード設定がされていないため、未入力でインポートしてください。  

## ソースコードのビルド

    go build genca.go 
    # または
    gox
    
    # 実行
    ./genca

## ライセンス

[LICENSE](./LICENSE)をご参照下さい。

## pkcs12配下のソースコードについて

クライアント証明書をPKCS #12形式に変換する際に、 [packer の pkcs12](https://github.com/hashicorp/packer/tree/master/builder/azure/pkcs12)を利用しています。


