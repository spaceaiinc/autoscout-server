# Rechub Server

## Prerequirements
Make sure this project depends on `golang`, `docker` and `docker-compose`.
So you need these items on your local machine before development.

## Get Started


1. Setup

```
$ make setup
```

2. Run docker-compose

```
$ make up
```

3. Run API application

```
$ make run
```

## Run Test

```
$ make test
```

## ブランチ構成
- `main` : 本番用。常にデプロイ可能な状態を保つ。

- `develop` : 開発用。開発の最新状態を反映。

- `feature/` : 新機能開発用。developから派生し、完成後にdevelopにマージ。

- `fix/` : バグ修正用。developから派生し、修正後にdevelopにマージ。

- `hotfix/` : 緊急修正用。mainから派生し、修正後にmainとdevelopにマージ。

- `release/` : リリース準備用。developから派生し、完成後にmainとdevelopにマージ。


## GitHub ラベル説明

GitHubにおける各ラベルは、以下の用途に使用されます：

- `enhancement`（機能追加）: 新たに追加された機能や、既存の機能の拡張を示します。主に新たな要素がプロジェクトに導入されたときに利用されます。

- `bug`（バグ）: ソフトウェア内のバグを示します。何らかの不具合や予期せぬ動作が発見されたときに利用されます。

- `emergency`（緊急）: 緊急を要する問題に対して使われます。これは通常、致命的なエラーや即時対応が必要な重大な問題を示します。

- `test`（テスト）: テストに関連する変更を示します。これは新たなテストの追加や、既存のテストの改良を示すために使用されます。

- `documentation`（ドキュメンテーション）: ドキュメンテーションに関する変更を示します。これは新たなドキュメンテーションの追加や、既存のドキュメンテーションの更新を示すために使用されます。

- `action`（アクション）: GitHub Actionsの設定ファイルやワークフローに対する変更を示します。これは新たなワークフローの追加や、既存のワークフローの改良を示すために使用されます。

バージョン変更ラベル:

- `major`（メジャー）: 互換性を破るような大きな変更があったときに付けられます。新しい機能の追加や既存の機能の大幅な変更などが含まれます。

- `minor`（マイナー）: 後方互換性を保ちつつ、新しい機能を追加したときに付けられます。

- `patch`（パッチ）: 後方互換性を保ったままバグ修正を行ったときに付けられます。

以上のラベルは、プロジェクトの進行状況を明確にし、メンバー間でのコミュニケーションを助けるためのものです。適切なラベルの使用は、プロジェクト管理を円滑にし、作業の効率を向上させます。
