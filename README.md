# go-backend-token-auth-middleware

## やりたいこと

* フロントエンドフレームワークのバックエンドサーバーとして認証と認可を行うためのミドルウェアをつくる
* フローとしてはフロントエンドのみで OAuth のアクセストークン発行まで行う
  * フロントエンドでのアクセストークン取得フローでは PKCE に対応する
* フロントエンドはアクセストークンを使ってバックエンドに認証を行う
* バックエンドはアクセストークンの詳細情報をプロバイダに問い合わせてユーザー名や対象アプリケーションを取得する
* バックエンドはアクセストークンが正常なら当アプリケーション用アクセストークンを発行する
* フロントエンドはバックエンドのAPI問い合わせ時に当アプリケーション用アクセストークンを添付し、バックエンドはこれを検証して認可を行う

## やること

* アプリ用トークンの発行と検証は以前つくったものを流用できそう
  * https://github.com/nirasan/go-jwt-handler

* プロバイダ毎にアクセストークンを検証する
  * Google, github, Amazon でのアクセストークンの検証について
    * https://stackoverflow.com/questions/12296017/how-to-validate-an-oauth-2-0-access-token-for-a-resource-server
  * OpenID Connect の ID Token 検証のコードは参考にできそう
    * https://github.com/coreos/go-oidc/blob/a4973d9a4225417aecf5d450a9522f00c1f7130f/example/idtoken/app.go
