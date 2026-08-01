# zsb_tmux_agent_notification

Per-pane notification counter, stored as a tmux **pane user option** (`@zsb_agent_notif`).
Used by `zsb_tmux_jira_picker` to show a red `(N)` badge per session row = number of
pending agent notifications across that session's panes.

State lives in the running tmux server (in-memory, dies with tmux). No sqlite, no files,
no deps.

## Command

```
zsb_tmux_agent_notification [--decrease] <session_name> <pane_id>
```

- **increase** (no flag): `@zsb_agent_notif` on `<pane_id>` += 1
- **`--decrease`**: reset `@zsb_agent_notif` on `<pane_id>` to `0`
- `<pane_id>` identifies the pane (and its session). `<session_name>` is accepted to match
  the tmux hook format string but is **unused** — pass anything (`_`).

The picker reads all panes in one call:
```sh
tmux list-panes -a -F '#{session_name}\t#{@zsb_agent_notif}'
```
and sums per session.

## Wiring 1 — coding-agent hooks (increase)

To bump the counter whenever the agent wants attention, hook the same events that already
play a sound. For **Claude Code**, this lives in `~/.claude/settings.json` under `hooks`.
Add this command entry next to each existing `zsb_play` sound entry:

```json
{
  "type": "command",
  "command": "zsb_tmux_agent_notification _ \"$TMUX_PANE\""
}
```

`$TMUX_PANE` is set by tmux in every pane's environment, so it resolves to the pane running
the agent. Outside tmux it's empty → tmux errors are ignored (no-op).

The four hook points currently wired (each already had a `zsb_play` sound):

| Section       | matcher           | meaning                    |
|---------------|-------------------|----------------------------|
| `Stop`        | `*`               | agent finished a turn      |
| `PostToolUse` | `ExitPlanMode`    | plan approved              |
| `PreToolUse`  | `AskUserQuestion` | about to ask a question    |
| `PreToolUse`  | `ExitPlanMode`    | about to present a plan    |

Example (`Stop` section) — the notification entry sits alongside the sound entry in the
same `hooks` array:

```json
"Stop": [
  {
    "matcher": "*",
    "hooks": [
      { "type": "command", "command": "zsb_play ~/.zsh-spell-book/src/media/sounds/aoe_farm.wav" },
      { "type": "command", "command": "zsb_tmux_agent_notification _ \"$TMUX_PANE\"" }
    ]
  }
]
```

Notes:
- `Stop` fires every turn → the count climbs while you're in another pane, and resets when
  you focus the pane (see Wiring 2).
- `ExitPlanMode` matches under both `PreToolUse` and `PostToolUse` → **+2** per plan exit.
  Drop one entry if that's too much.

### Other coding agents

Any agent with lifecycle hooks works the same way: run
`zsb_tmux_agent_notification _ "$TMUX_PANE"` on its "done / needs input" events. Only
requirements: the command runs inside the tmux pane's environment (so `$TMUX_PANE` is set)
and the binary is on `PATH`.

## Wiring 2 — tmux focus-in hook (decrease/reset)

In your tmux config (`.tmux.conf`), reset a pane's counter when you look at it:

```sh
set-hook -g pane-focus-in 'run-shell "zsb_tmux_agent_notification --decrease \"#{session_name}\" \"#{pane_id}\""'
```

## Build / install

Binary must be on `PATH` (repo ships `src-go/bin` on `PATH`, see repo README). Built by the
auto-build server on change, or manually:

```sh
make go-build zsb_tmux_agent_notification
```

## Verify

```sh
pane=$(tmux display -p '#{pane_id}')
zsb_tmux_agent_notification _ "$pane"; zsb_tmux_agent_notification _ "$pane"
tmux show-options -pqv -t "$pane" @zsb_agent_notif   # -> 2
zsb_tmux_agent_notification --decrease _ "$pane"
tmux show-options -pqv -t "$pane" @zsb_agent_notif   # -> 0
```
