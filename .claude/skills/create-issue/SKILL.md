---
name: create-issue
description: "GitHub Issueを作成する。ユーザーが「Issueを作って」「Issue作成して」「GitHubのIssueを追加して」「課題を登録して」「Issueを立てて」と言ったときに使用する。"
argument-hint: "[概要]"
---

# Issue 作成スキル

GitHub Issue を作成する。開発ワークフローに沿った形で、必要な情報を揃えて作成する。

## 1. 情報を収集する

以下が不明な場合はユーザーに確認する。自然な会話の流れで判断できる場合は確認不要。

- **タイトル**: 何をするか（簡潔に）
- **概要**: なぜ必要か、何を実現するか
- **ラベル**: 以下から最も適切なものを1つ選ぶ

| ラベル | 用途 |
|-------|------|
| `bug` | バグ報告・修正 |
| `feature` | 新機能の実装 |
| `enhancement` | 既存機能の改善 |
| `documentation` | ドキュメントの追加・修正 |
| `chore` | 保守・依存関係更新などの雑務 |
| `question` | 質問・調査 |

ラベルはユーザーが明示しない場合でも、内容から推測して提案する。

## 2. 分割を判断する

タスクが以下のいずれかに該当する場合、親子Issue構造に分割する：

- 独立してレビュー・テストできる複数の作業が含まれている
- 1つのPRが大きくなりすぎる見込みがある（目安: 複数の関心事にまたがる）
- ユーザーが分割を指示している

分割が不要な場合はステップ3に進む。分割する場合はステップ4に進む。

## 3. 単独 Issue を作成する（分割不要の場合）

```bash
gh issue create --title "タイトル" --body "概要" --label "ラベル" --assignee "@me"
```

担当者は `@me`（自分自身）を設定する。

作成後、ブランチ名フォーマットを案内する：

```
issues/番号_やることをざっくりとケバブケース
```

例：`issues/5_fix-login-error`

## 4. 親子 Issue を作成する（分割する場合）

GitHub の sub-issue 機能を使って親子関係を構築する。

<!-- NOTE: gh CLI にはまだ sub-issue のネイティブサポートがないため、
     gh api 経由で REST API を呼び出している。
     gh CLI に --parent オプション等が追加された場合はそちらに移行すること。
     参考: https://github.com/cli/cli/issues/10298 -->

### 4-1. 子タスクを洗い出す

タスクを独立した作業単位に分割し、実装順序を決める。各子タスクは単独でレビュー・テスト可能な粒度にする。

### 4-2. 親 Issue を作成する

```bash
gh issue create --title "親タイトル" --body "概要" --label "ラベル" --assignee "@me"
```

### 4-3. 子 Issue を作成し、親に紐づける

各子タスクごとに Issue を作成し、sub-issue API で親に紐づける：

```bash
# 子Issueを作成
gh issue create --title "子タスクのタイトル" --body "概要" --label "ラベル" --assignee "@me"

# 子IssueのID（番号ではなくID）を取得し、親Issueに紐づける（REST API）
# sub_issue_id にはIssue番号ではなくIssue ID（整数）が必要
CHILD_ID=$(gh api repos/{owner}/{repo}/issues/{子Issue番号} --jq '.id')
gh api repos/{owner}/{repo}/issues/{親Issue番号}/sub_issues \
  -F sub_issue_id=$CHILD_ID
```

`{owner}/{repo}` は `gh repo view --json nameWithOwner -q .nameWithOwner` で取得できる。
`-F`（大文字）を使うことで値が整数として送信される（`-f` だと文字列になりAPIエラーになる）。

### 4-4. 紐づけを確認する

```bash
# 親Issueのsub-issue一覧を取得して確認
gh api repos/{owner}/{repo}/issues/{親Issue番号}/sub_issues
```

### 4-5. ブランチ構造を案内する

```
main
 └── issues/親番号_親の説明        ← 親ブランチ（from main）
      ├── issues/子番号1_子の説明   ← 子ブランチ1（from 親ブランチ）
      └── issues/子番号2_子の説明   ← 子ブランチ2（from 親ブランチ）
```

- 親ブランチは `main` から作成する
- 子ブランチは親ブランチから作成する
- 子PRは親ブランチに向けて作成する
- 全ての子PRをマージ後、親PRを `main` にマージする

## 5. 結果を報告する

作成した Issue の URL と番号をユーザーに伝える。親子構造の場合は全Issueの一覧と、次に着手すべき子Issueを案内する。


## エラー時の自己修正

このスキルの実行中にエラーや不具合が発生した場合：

1. ユーザーにエラー内容と原因を報告する
2. エラーの根本原因がスキルの記述（手順の誤り、コマンドの間違い、考慮漏れなど）にある場合は、このSKILL.mdを修正して再発を防止する
3. 修正内容をコミットに含める