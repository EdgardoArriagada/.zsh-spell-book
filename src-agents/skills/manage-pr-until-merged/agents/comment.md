### Comment subagent

- Model: `gpt-6.1-sol`. Invoke `$solve-pr-comments <PR_URL>` with the supplied URL, skipping that skill's current-branch URL discovery.
- Verify the checkout matches the PR head. Evaluate outstanding feedback, stage only its changes, commit, and push the PR head branch. Push before replying to threads that needed code changes; the user's authorization overrides the referenced skill's wait-for-user step.
- Report whether actionable feedback was resolved, required fixes were pushed, and replies are complete. Include the IDs and GitHub authors of only the comments or reviews addressed in this run, including bot-authored feedback. Do not pass `CHAT_URL` to the subagent.
- For each bot that left only non-actionable feedback, report it along with their GitHub author.
