---
name: solve-pr-comments
description: Address review comments on the current pull request, including code fixes and GitHub replies.
disable-model-invocation: true
---

# Solve PR Comments

1. Get the pr_url. `~/.agents/skills/solve-pr-comments/scripts/get-pr-url.sh`

2. Use receiving-code-review skill to address all pr comments <pr_url>

- work on those that need code fixes
- for each comment that does not need code changes, post a GitHub reply via `gh pr comment <pr_url> --body "<response>"`
- wait for me to commit and push before answering comments that needed code changes
