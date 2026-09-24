---
name: solve-pr-comments
description: Address review comments on the current pull request, including code fixes and GitHub replies.
disable-model-invocation: true
---

# Solve PR Comments

1. Get the pr_url. `~/.agents/skills/solve-pr-comments/scripts/get-pr-url.sh`

2. Use receiving-code-review skill to address all pr comments <pr_url>

- work on those that need code fixes
- reply to inline comments in their thread with `scripts/reply-review-comment.sh <pr_url> <comment_id>` and pass the body on stdin; this derives the Pull Request scope from `pr_url` and prevents shell evaluation of Markdown
- use `gh pr comment <pr_url> --body-file <file>` only for Pull Request-level comments
- wait for me to commit and push before answering comments that needed code changes
