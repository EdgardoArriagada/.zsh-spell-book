#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 2 ]]; then
  echo "usage: reply-review-comment.sh <pr_url> <comment_id>" >&2
  exit 2
fi

pr_url=$1
comment_id=$2
path=${pr_url#https://github.com/}
owner_repo=${path%%/pull/*}
pr_number=${path##*/pull/}
pr_number=${pr_number%%/*}

if [[ $owner_repo == "$path" || ! $pr_number =~ ^[0-9]+$ || ! $comment_id =~ ^[0-9]+$ ]]; then
  echo "invalid Pull Request URL or comment ID" >&2
  exit 2
fi

body=$(</dev/stdin)
if [[ -z $body ]]; then
  echo "reply body is required on stdin" >&2
  exit 2
fi

gh api --method POST \
  "repos/$owner_repo/pulls/$pr_number/comments/$comment_id/replies" \
  --raw-field "body=$body"
