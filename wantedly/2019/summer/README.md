# Coding Challenge Docker WebApp

## 課題1

Dockerイメージの作成 & 起動方法

```
$ docker build . -t hello_world
$ docker run -d -p 8080:8080 hello_world
```

## 課題2

https://coding-challenge-docker-webapp.herokuapp.com

HerokuのFree Dynoを利用しているため基本的にはスリープ状態になっており、
その場合はレスポンスが返るまでに若干の時間を要する

## 課題3

Dockerイメージの作成 & 起動方法

```
$ docker-compose build
$ docker-compose up -d
```

課題2でHeroku上に作成したコンテナの上でも同様のアプリが動いている (PostgreSQLはHerokuのアドオンを使っている)

また、権限の都合上、ローカルのDBのタイムゾーンはJST、HerokuのDBのタイムゾーンはUTCになっている
