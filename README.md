# Example sqlc with pgx and sqldef

refs:
- [Getting started with PostgreSQL — sqlc 1.20.0 documentation](https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html)
- [Using Go and pgx — sqlc 1.20.0 documentation](https://docs.sqlc.dev/en/stable/guides/using-go-and-pgx.html)
- https://github.com/k0kubun/sqldef
- [Dockerfileで対象プラットフォームによって処理分岐させる](https://zenn.dev/ytdrep/articles/d65c26201042eb)
- [sqlc と pgxpool でトランザクション](https://zenn.dev/shiguredo/articles/sqlc-pgxpool-transaction)


## TODO

- [ ] このプロジェクトの [Connect](https://connectrpc.com/docs/go/getting-started/) 版を作る
- [ ] main.go が分厚いのでダイエットさせる (root dir に移動させる)

## Requirements

- Go
- sqlc
- sqldef
- Docker
- Docker Compose

## Commands

DB 起動
```
make up
```

実行されるマイグレーションの確認
```
make migration-dry-run
```

sqlc によるコード生成
```
make gen
```

アプリケーションのビルド
```
make build
```

### DB定義を変更するとき

schema.sql を編集
```shell
# check
make migration-dry-run

# run
make migration
```

必要があれば、query.sql を編集して、sqlc generate
```shell
make gen
```

## Architecture

TODO: 書く
