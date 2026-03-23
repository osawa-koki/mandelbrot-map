---
name: merge-pr
description: "PRをマージする。ユーザーが「PRをマージして」「プルリクをマージして」「マージして」「マージ前チェックをして」「PRをマージ可能か確認して」と言ったときに使用する。"
argument-hint: "[PR番号]"
---

# PR マージスキル

PRのマージ前チェックを行い、条件を満たしていればマージする。

## 1. PR を特定する

引数で PR 番号が指定された場合はそれを使う。指定がなければ現在のブランチに対応する PR を `gh pr view` で取得する。

## 2. PR の種類を判断する

`gh pr view <番号> --json title,body,labels,headRefName,baseRefName,statusCheckRollup` で情報を取得する。

`baseRefName` を確認する：
- `main`（またはデフォルトブランチ）→ 通常のPRまたは親PR
- それ以外 → 子PR（親ブランチへのマージ）

親PRかどうかは、対応するIssueにsub-issueが存在するかで判断する：

```bash
# 親Issueのsub-issue一覧を取得（REST API）
gh api repos/{owner}/{repo}/issues/{Issue番号}/sub_issues
```

<!-- NOTE: gh CLI にはまだ sub-issue のネイティブサポートがないため、
     gh api 経由で REST API を呼び出している。
     gh CLI に対応コマンドが追加された場合はそちらに移行すること。
     参考: https://github.com/cli/cli/issues/10298 -->

## 3. マージ前チェック

以下をすべて確認する。1つでも問題があればマージせず、ユーザーに報告する。

### CI チェック

すべての check が `SUCCESS` であることを確認する。失敗・実行中の check がある場合はマージしない。

### チェックリスト完了確認

- `## Checklist` セクションが存在すること
- すべての `- [ ]` が `- [x]` になっていること（未完了項目が残っていないこと）

### PR 本文チェック

- `closes #番号` の記載があること（Issueの自動クローズに必要）
- タイトルと概要が実際の変更内容と一致していること（古いままになっていないか）

### ブランチ・Issue 対応チェック

- PR のブランチ名に Issue 番号が含まれていること（開発ワークフローの規約に沿っているか）
- `closes #番号` で参照している Issue が実際に存在し、PR の変更内容がその Issue のスコープを適切にカバーしていること
  - Issue の内容は `gh issue view <番号>` で取得して確認する

### ラベルチェック

- PR にラベルが設定されていること
- ラベルの種類（`bug` / `feature` / `enhancement` / `documentation` / `chore` / `question`）が変更内容と合っていること

### 不要ファイルチェック

- PR の差分（`gh pr diff <番号>`）にログ・一時ファイル・ビルド成果物などが含まれていないこと

### 親PRの追加チェック

親PRの場合、以下も確認する：

- sub-issue API で子Issueを取得し、全ての子Issueがクローズ済みであることを確認する
- 未クローズの子Issueがある場合、対応する子PRが未マージの可能性があるためマージしない

```bash
# 子Issueの状態を確認（REST API）
gh api repos/{owner}/{repo}/issues/{親Issue番号}/sub_issues
```

## 4. マージ実行

すべてのチェックを通過したら、マージする。

```bash
# 通常のPR・親PR → squash merge（mainの履歴をきれいに保つ）
gh pr merge <番号> --squash

# 子PR → merge commit（親ブランチに個別の変更履歴を残す）
gh pr merge <番号> --merge
```

子PRにmerge commitを使う理由：親ブランチ上で各子タスクの変更履歴を個別に残すため。最終的に親PRをsquashマージする際に1コミットにまとまる。

## 5. 子PRのIssueを手動クローズする

子PRの場合、`closes #番号` はmainへのマージ時にのみ自動クローズされる。
親ブランチへのマージでは自動クローズされないため、**マージ後に必ず手動でIssueをクローズする**。

```bash
gh issue close <Issue番号> --comment "PR #番号で親ブランチにマージ済み"
```

通常のPR（baseがmain）の場合は自動クローズされるため不要。

## 6. main ブランチに移動して最新化する

マージ完了後、ローカルを最新の main に切り替える：

```bash
git checkout main && git pull origin main
```

子PRの場合（baseが親ブランチ）は、親ブランチに切り替えて最新化する：

```bash
git checkout <親ブランチ名> && git pull origin <親ブランチ名>
```

## 7. 結果報告

マージ成功後、以下を報告する：
- マージされた PR のタイトルと番号
- クローズされた Issue 番号

子PRをマージした場合は、残りの子PRの状況と、次に着手すべきタスクを案内する。

親PRをマージした場合は、全子Issueのクローズ状況を確認して報告する。

## チェック失敗時の対応

チェックに失敗した場合は、マージせずに以下を報告する：
- どのチェックが失敗したか
- 解消するために何が必要か

ユーザーが「それでもマージして」と明示的に指示した場合のみ、問題を再度示したうえでマージを実行する。


## エラー時の自己修正

このスキルの実行中にエラーや不具合が発生した場合：

1. ユーザーにエラー内容と原因を報告する
2. エラーの根本原因がスキルの記述（手順の誤り、コマンドの間違い、考慮漏れなど）にある場合は、このSKILL.mdを修正して再発を防止する
3. 修正内容をコミットに含める