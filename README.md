# Poppo Server Side
ぽっぽっぽー、はとぽっぽ。

## Endpoints
```
RestFul!!!
[POST]   /apps                                     ... Register Installation
[PUT]    /apps                                     ... Update Installation
[POST]   /apps                                     ... Do Login
[DELETE] /apps                                     ... Do Logout
[GET]    /tweets{?limit,maxId,sinceId,following,q} ... Search
[POST]   /tweets                                   ... Post Tweet
[GET]    /tweets/{id}                              ... Get Tweet
[DELETE] /tweets/{id}                              ... Remove Tweet
[POST]   /tweets/{id}/like                         ... Do Tweet Like
[DELETE] /tweets/{id}/like                         ... Remove Tweet Like
[POST]   /users                                    ... Register User
[GET]    /users{?screenName}                       ... Get User From ScrennName
[GET]    /users/{id}                               ... Get User
[PUT]    /users/{id}                               ... Update User
[POST]   /users/{id}/follow                        ... Do User Follow
[DELETE] /users/{id}/follow                        ... Remove User Follow
[GET]    /users/{id}/following{?limit,offset}      ... Get Following Users
[GET]    /users/{id}/follower{?limit,offset}       ... Get Follower Users
[GET]    /users/{id}/tweets{?limit,maxId}          ... Get User Tweets
[GET]    /users/{id}/liked/tweets{?limit,maxId}    ... Get Liked Tweets
[GET]    /notifications/count                      ... Count Unread Notifications
[GET]    /notifications{?limit,maxId,sinceId}      ... Get Notifications
[POST]   /images{?data,type}                       ... Upload Image into AWS S3
```
more: https://github.com/techcampman/twitter-d-api-document/


## Package structure

```
.
|-- Godeps          ... 依存関係管理
|-- api             ... APIメインパッケージ
    |-- v1          ... コントローラー郡
|-- bridge          ... 外部サービス連携 (AWS SNS)
|-- constant        ... 定数定義
|-- db              ... DBへの接続 Interface郡
    |-- collection  ... MongoDBのCollection管理
    |-- mongo       ... MongoDBのマネージャー (thanks dogenzaka/mds !)
    |-- redis       ... Redisのマネージャー
    |-- s3          ... AWS S3の実装
|-- entity          ... 実体定義 (+ MongoDBに対するIndex指定)
|-- env             ... Enviroment Interface定義
    |-- on          ... 各環境の設定定義 (local, release)
|-- errors          ... エラーハンドリング定義
|-- jsonschema      ... Json Schema定義
|-- logger          ... ログ出力系定義
|-- middleware      ... gin.Contextに対するミドルウェア群
|-- scripts         ... 便利script群 (coverage, fixture)
    |-- fixtures    ... テストデータ生成スクリプト等
|-- service         ... サービス群
|-- utils           ... 雑多なものたち
```
