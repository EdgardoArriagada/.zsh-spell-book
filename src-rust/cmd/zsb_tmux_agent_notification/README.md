# zsb_tmux_agent_notification

Per-pane notification **state**, stored as a tmux **pane user option** (`@zsb_agent_notif`).
Quad-state: `0` idle · `1` finished · `2` working · `3` manual. `zsb_tmux_jira_picker` reads
it per session and shows a **blue** count (working), **red** count (finished), and **yellow**
count (manually flagged).

State lives in the running tmux server (in-memory, dies with tmux). No sqlite, no files,
no deps.

## Command

```
zsb_tmux_agent_notification [--finished|--force-finished|--working|--clear-finished|--manual] <session_name> <pane_id>
```

- **`--finished`** (default when no flag): wait 10 seconds, then set `@zsb_agent_notif` = `1`
  (red). Another `--finished` restarts the timer; any notification except `--clear-finished`
  cancels it. If the pane is focused when the timer fires, clears to `0` instead (you're already
  watching).
- **`--force-finished`**: cancel any pending timer and finish immediately, even when focused.
- **`--working`**: set `@zsb_agent_notif` = `2` (blue). Does **not** skip on focus — its hook
  fires while the pane is still focused at submit time, before you switch away.
- **`--clear-finished`**: reset `@zsb_agent_notif` = `0` **only if it's `1` (finished)**. A
  working (`2`) or manually-flagged (`3`) pane is left untouched.
- **`--manual`**: toggle `@zsb_agent_notif` between `0` and `3` (yellow). No-op if the pane
  is working (`2`) or finished (`1`) — those states take priority. `--finished` and `--working`
  override `3` normally.
- `<pane_id>` identifies the pane (and its session). `<session_name>` is accepted to match
  the tmux hook format string but is **unused** — pass anything (`_`).

Each state also updates the tmux window name: a working glyph (`󰔟`) and/or the finished bell
(`󰂟`) are appended, recomputed from all panes in the window and stripped on `--clear`.

The picker reads all panes in one call:
```sh
tmux list-panes -a -F '#{session_name}\t#{@zsb_agent_notif}'
```
and tallies panes per session by state (`1`→finished, `2`→working).

## Wiring 1 — coding-agent hooks (finished / working)

Hook the same events that already play a sound. For **Claude Code**, this lives in
`~/.claude/settings.json` under `hooks`. Add a command entry next to each existing `zsb_play`
sound entry.

`$TMUX_PANE` is set by tmux in every pane's environment, so it resolves to the pane running
the agent. Outside tmux it's empty → tmux errors are ignored (no-op).

**Finished** — the agent is done / wants attention (`--finished`, red):

| Section       | matcher           | meaning                    |
|---------------|-------------------|----------------------------|
| `Stop`        | `*`               | agent finished a turn      |
| `PostToolUse` | `ExitPlanMode`    | plan approved              |
| `PreToolUse`  | `AskUserQuestion` | about to ask a question    |
| `PreToolUse`  | `ExitPlanMode`    | about to present a plan    |

```json
{
  "type": "command",
  "command": "zsb_tmux_agent_notification --finished _ \"$TMUX_PANE\""
}
```

**Working** — the agent just started churning (`--working`, blue). Fire on the turn's start,
e.g. `UserPromptSubmit`:

```json
"UserPromptSubmit": [
  {
    "hooks": [
      { "type": "command", "command": "zsb_tmux_agent_notification --working _ \"$TMUX_PANE\"" }
    ]
  }
]
```

`--working` sets `2` and the next `--finished` (Stop) overwrites it with `1`, so a session
reads blue while an agent runs, then red once it's done — until you focus the pane (Wiring 2).

Example (`Stop` section) — the notification entry sits alongside the sound entry in the
same `hooks` array:

```json
"Stop": [
  {
    "matcher": "*",
    "hooks": [
      { "type": "command", "command": "zsb_play ~/.zsh-spell-book/src/media/sounds/aoe_farm.wav" },
      { "type": "command", "command": "zsb_tmux_agent_notification --finished _ \"$TMUX_PANE\"" }
    ]
  }
]
```

Note: `ExitPlanMode` matches under both `PreToolUse` and `PostToolUse` → the second
`--finished` restarts the 10-second timer.

### Other coding agents

Any agent with lifecycle hooks works the same way: run
`zsb_tmux_agent_notification --working _ "$TMUX_PANE"` on "turn started" and
`--finished` on "done / needs input". Only requirements: the command runs inside the tmux
pane's environment (so `$TMUX_PANE` is set) and the binary is on `PATH`.

## Wiring 2 — tmux focus-in hook (clear finished)

In your tmux config (`.tmux.conf`), clear a pane's **finished** state when you look at it (a
still-working pane keeps its blue badge):

```sh
set-hook -g pane-focus-in 'run-shell "zsb_tmux_agent_notification --clear-finished \"#{session_name}\" \"#{pane_id}\""'
```

## Build / install

Binary must be on `PATH` (repo ships `src-rust/bin` on `PATH`, see repo README). Built by the
auto-build server on change, or manually:

```sh
make rust-build zsb_tmux_agent_notification
```

## Verify

```sh
pane=$(tmux display -p '#{pane_id}')
zsb_tmux_agent_notification --working _ "$pane"
tmux show-options -pqv -t "$pane" @zsb_agent_notif   # -> 2  (window name gets 󰔟)
zsb_tmux_agent_notification --clear-finished _ "$pane"
tmux show-options -pqv -t "$pane" @zsb_agent_notif   # -> 2  (working kept — only clears 1)
```

`--finished` on the **focused** pane clears to `0` (removes hourglass, no bell). To see it set
`1` (swap glyph to `󰂟`), target a pane in another window/session that you're not currently on.
`--clear-finished` on that finished pane then resets it to `0`.
