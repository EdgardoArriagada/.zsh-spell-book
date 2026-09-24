#!/usr/bin/env bash
set -euo pipefail

script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

gh() {
  printf '%s\n' "$@"
}
export -f gh

output=$(printf '%s' 'Fixed in `abc123`. Keep $(literal) unchanged.' | \
  "$script_dir/reply-review-comment.sh" \
    https://github.com/owner/repo/pull/1854 \
    4063904478)

[[ $output == *"repos/owner/repo/pulls/1854/comments/4063904478/replies"* ]]
[[ $output == *'body=Fixed in `abc123`. Keep $(literal) unchanged.'* ]]
