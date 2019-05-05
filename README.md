:jp: [:us:](./README_en.md)

# My First Go Server

## Go Simple Server

`simple_server/` ディレクトリ

### 使い方

    cd simple_server
    docker build .
    docker run -p 8080:8080 [DOCKER_IMAGE]

### 動作

    curl -XGET -H 'Content-Type:application/json' http://localhost:8080/

response:

    {"message":"Hello World!!"}

## インターネットへの公開

http://35.247.11.48:8080/ (一定期間後削除)

### GCP Cloud Registry

- Cloud Build を利用してDockerイメージをアップロード
- [GCPのリファレンス](https://cloud.google.com/cloud-build/docs/quickstart-docker)

Dockerfileによるビルド

    gcloud builds submit --tag gcr.io/[PROJECT_ID]/[DOCKER_IMAGE] .

### デプロイ

1. Kubernetes Engine > クラスタ タブで「コンテナをデプロイ」から上記でpushしたDockerイメージを選択
2. ワークロード タブで「サービスの公開」からロードバランサを有効化、ポート番号8080に設定

## Go RESTful Server

`restful_server/` ディレクトリ

### 使い方

    cd restful_server
    docker-compose build
    docker-compose up

### Request サンプル

#### GET /users

    curl -XGET -H 'Content-Type:application/json' http://localhost:8080/users

#### GET /users/:id

    curl -XGET -H 'Content-Type:application/json' http://localhost:8080/users/1

#### POST /users

    curl -XPOST -H 'Content-Type:application/json' http://localhost:8080/users -d '{ "name": "test", "email": "hoge@example.com" }'

#### PUT /users/:id

    curl -XPUT -H 'Content-Type:application/json' http://localhost:8080/users/1 -d '{ "name": "test2", "email": "fuga@example.com" }'

#### DELETE /users/:id

    curl -XDELETE -H 'Content-Type:application/json' http://localhost:8080/users/1
