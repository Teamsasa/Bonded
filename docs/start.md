# はじめに

## 前提条件

- Dockerのインストール
- Goのインストール
- aws-sam-cliのインストール

## アプリケーションの起動

下記コマンドでアプリケーションをBuild、起動します。
```sh
# Build
make build

# run
make start-all
```

## Makefile

開発に必要なコマンドは基本Makefileに記述されています
```
  build                Build SAM application
  compose-down         Stop and remove Docker containers
  compose-up           Start Docker containers
  dynamodb-init        Initialize DynamoDB Local using an external script
  fmt                  Format all Go code files
  help                 Display this help message
  start-all            Start and initialize DynamoDB, then start SAM API
  sam-api              Start SAM API
```
