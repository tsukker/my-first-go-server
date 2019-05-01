:jp: [:us:](./README_en.md)

# My First Go Server

## Go Simple Server

`simple_server/` ディレクトリ

### Usage

    cd simple_server
    docker build .
    docker run -p 8080:8080 [DOCKER_IMAGE]

## インターネットへの公開

(IPアドレス):8080 (一定期間後削除)

### GCP Cloud Registry

- Cloud Build を利用してDockerイメージをアップロード
- [GCPのリファレンス](https://cloud.google.com/cloud-build/docs/quickstart-docker)

Dockerfileによるビルド

    gcloud builds submit --tag gcr.io/[PROJECT_ID]/[DOCKER_IMAGE] .

### デプロイ

1. Kubernetes Engine > クラスタ タブで「コンテナをデプロイ」から上記でpushしたDockerイメージを選択
2. ワークロード タブで「サービスの公開」からロードバランサを有効化、ポート番号8080に設定