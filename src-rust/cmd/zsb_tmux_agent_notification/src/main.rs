use std::process::{self, Command};

const NOTIF_VAR: &str = "@zsb_agent_notif";
const FINISHED_SUFFIX: &str = " \u{f009a}"; // bell — agent finished / needs attention
const WORKING_SUFFIX: &str = " \u{f051f}"; //  hourglass — agent still working
const MANUAL_SUFFIX: &str = " \u{f0e47}"; //   flag — manually flagged

const CLEAR: &str = "0";
const FINISHED: &str = "1";
const WORKING: &str = "2";
const MANUAL: &str = "3";

// tmux_output runs tmux and returns trimmed stdout, or "" on any error.
fn tmux_output(args: &[&str]) -> String {
    Command::new("tmux")
        .args(args)
        .output()
        .map(|o| String::from_utf8_lossy(&o.stdout).trim().to_owned())
        .unwrap_or_default()
}

fn pane_notif(pane: &str) -> String {
    tmux_output(&["show-options", "-pqv", "-t", pane, NOTIF_VAR])
}

fn window_id(pane: &str) -> String {
    tmux_output(&["display-message", "-p", "-t", pane, "#{window_id}"])
}

fn window_name(pane: &str) -> String {
    tmux_output(&["display-message", "-p", "-t", pane, "#{window_name}"])
}

// pane_is_focused reports whether the user is currently looking at this pane:
// its session is attached, its window is active, and the pane is active.
fn pane_is_focused(pane: &str) -> bool {
    tmux_output(&[
        "display-message",
        "-p",
        "-t",
        pane,
        "#{&&:#{session_attached},#{&&:#{window_active},#{pane_active}}}",
    ]) == "1"
}

// base_name strips any trailing known notif suffixes (in either order).
fn base_name(mut name: &str) -> &str {
    loop {
        if let Some(s) = name.strip_suffix(FINISHED_SUFFIX) {
            name = s;
        } else if let Some(s) = name.strip_suffix(WORKING_SUFFIX) {
            name = s;
        } else if let Some(s) = name.strip_suffix(MANUAL_SUFFIX) {
            name = s;
        } else {
            return name;
        }
    }
}

// window_states reports whether any pane in the window is finished (1), working (2), or manual (3).
fn window_states(win_id: &str) -> (bool, bool, bool) {
    let out = tmux_output(&["list-panes", "-t", win_id, "-F", "#{@zsb_agent_notif}"]);
    let (mut finished, mut working, mut manual) = (false, false, false);
    for line in out.lines() {
        match line.trim() {
            FINISHED => finished = true,
            WORKING => working = true,
            MANUAL => manual = true,
            _ => {}
        }
    }
    (finished, working, manual)
}

// refresh_window_name recomputes the window-name suffixes from the window's pane
// states: working glyph first, then finished bell. Renames only if changed.
fn refresh_window_name(pane: &str) {
    let (finished, working, manual) = window_states(&window_id(pane));
    let cur = window_name(pane);
    let mut want = base_name(&cur).to_owned();
    if working {
        want.push_str(WORKING_SUFFIX);
    }
    if finished {
        want.push_str(FINISHED_SUFFIX);
    }
    if manual {
        want.push_str(MANUAL_SUFFIX);
    }
    if want != cur {
        Command::new("tmux")
            .args(["rename-window", "-t", pane, &want])
            .status()
            .ok();
    }
}

fn set_pane_state(pane: &str, state: &str) {
    Command::new("tmux")
        .args(["set-option", "-p", "-t", pane, NOTIF_VAR, state])
        .status()
        .ok();
}

// zsb_tmux_agent_notification [--finished|--working|--clear-finished] <session_name> <pane_id>
// pane_id identifies the pane (and its session); session_name matches the hook
// signature but is unused. No flag defaults to --finished.
fn apply_flag(flag: &str, pane: &str) -> bool {
    match flag {
        "" | "--finished" => {
            // Already watching: clear working state without ringing the bell.
            if pane_is_focused(pane) {
                set_pane_state(pane, CLEAR);
            } else {
                set_pane_state(pane, FINISHED);
            }
            true
        }
        "--force-finished" => {
            set_pane_state(pane, FINISHED);
            true
        }
        "--working" => {
            // working fires while the pane is still focused (at submit time); don't
            // skip on focus — the focus-out/next --clear resets it.
            set_pane_state(pane, WORKING);
            true
        }
        "--clear-finished" => {
            // Only clear the finished (1) state; leave working (2) and manual (3) alone.
            if pane_notif(pane) == FINISHED {
                set_pane_state(pane, CLEAR);
            }
            true
        }
        "--manual" => {
            // Toggle: 3 → 0, 0 → 3. Ignore if pane is working (2) or finished (1).
            match pane_notif(pane).as_str() {
                MANUAL => set_pane_state(pane, CLEAR),
                CLEAR | "" => set_pane_state(pane, MANUAL),
                _ => {}
            }
            true
        }
        _ => false,
    }
}

fn run(args: &[String]) -> i32 {
    let mut flag = "";
    let mut rest: &[String] = args;
    if let Some(first) = args.first() {
        if first.starts_with("--") {
            flag = first;
            rest = &args[1..];
        }
    }
    if rest.len() < 2 {
        return 1;
    }
    let pane = &rest[1];

    if !apply_flag(flag, pane) {
        return 1;
    }
    refresh_window_name(pane);
    0
}

fn main() {
    let args: Vec<String> = std::env::args().skip(1).collect();
    process::exit(run(&args));
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn base_name_strips_suffixes() {
        let cases = [
            ("no suffix", "myrepo".to_owned(), "myrepo"),
            ("finished only", format!("myrepo{FINISHED_SUFFIX}"), "myrepo"),
            ("working only", format!("myrepo{WORKING_SUFFIX}"), "myrepo"),
            ("manual only", format!("myrepo{MANUAL_SUFFIX}"), "myrepo"),
            (
                "both working then finished",
                format!("myrepo{WORKING_SUFFIX}{FINISHED_SUFFIX}"),
                "myrepo",
            ),
            (
                "both finished then working",
                format!("myrepo{FINISHED_SUFFIX}{WORKING_SUFFIX}"),
                "myrepo",
            ),
            (
                "double finished",
                format!("myrepo{FINISHED_SUFFIX}{FINISHED_SUFFIX}"),
                "myrepo",
            ),
            (
                "all three suffixes",
                format!("myrepo{WORKING_SUFFIX}{FINISHED_SUFFIX}{MANUAL_SUFFIX}"),
                "myrepo",
            ),
        ];
        for (name, input, want) in cases {
            assert_eq!(base_name(&input), want, "case: {name}");
        }
    }
}
