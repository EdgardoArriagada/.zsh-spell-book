### Pipeline subagent

- Model: `gpt-6.1-sol`. Goal: **"solve pr pipeline"** for `<PR_URL>`.
- Verify the checkout matches the PR head. Inspect failed checks, fix the root cause, verify the fix, stage only its changes, commit, and push. Do not post the chat message for pipeline work.
