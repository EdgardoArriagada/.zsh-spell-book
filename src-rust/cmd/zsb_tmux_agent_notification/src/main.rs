use std::process::{self, Command, Stdio};
use std::thread;
use std::time::{Duration, SystemTime, UNIX_EPOCH};

const NOTIF_VAR: &str = "@zsb_agent_notif";
const FINISH_TIMER_VAR: &str = "@zsb_agent_finish_timer";
const DEBOUNCED_FINISH: &str = "--debounced-finished";
const FINISH_DELAY: Duration = Duration::from_secs(10);
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

fn pane_id(pane: &str) -> String {
    tmux_output(&["display-message", "-p", "-t", pane, "#{pane_id}"])
}

fn finish_timer(pane: &str) -> String {
    tmux_output(&["show-options", "-pqv", "-t", pane, FINISH_TIMER_VAR])
}

fn set_finish_timer(pane: &str, token: &str) -> bool {
    Command::new("tmux")
        .args(["set-option", "-p", "-t", pane, FINISH_TIMER_VAR, token])
        .status()
        .is_ok_and(|status| status.success())
}

fn clear_finish_timer(pane: &str) {
    Command::new("tmux")
        .args(["set-option", "-pu", "-t", pane, FINISH_TIMER_VAR])
        .status()
        .ok();
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

fn play_finish_sound(flag: &str) {
    if !matches!(flag, "" | "--finished" | "--force-finished") {
        return;
    }
    if let Ok(sound) = std::env::var("ZSB_FINISH_SOUND") {
        if !sound.is_empty() {
            zsb_play::spawn(&sound).ok();
        }
    }
}

fn is_finished(flag: &str) -> bool {
    matches!(flag, "" | "--finished")
}

fn cancels_finish(flag: &str) -> bool {
    !is_finished(flag) && flag != "--clear-finished"
}

fn schedule_finish(pane: &str) -> bool {
    let token = format!(
        "{}-{}",
        process::id(),
        SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap_or_default()
            .as_nanos()
    );
    if !set_finish_timer(pane, &token) {
        return false;
    }

    let spawned = std::env::current_exe().ok().and_then(|exe| {
        Command::new(exe)
            .args([DEBOUNCED_FINISH, &token, "_", pane])
            .stdin(Stdio::null())
            .stdout(Stdio::null())
            .stderr(Stdio::null())
            .spawn()
            .ok()
    });
    if spawned.is_none() && finish_timer(pane) == token {
        clear_finish_timer(pane);
    }
    spawned.is_some()
}

fn run_debounced_finish(args: &[String]) -> i32 {
    if args.len() != 4 {
        return 1;
    }
    let token = &args[1];
    let pane = &args[3];
    thread::sleep(FINISH_DELAY);
    if finish_timer(pane) != *token {
        return 0;
    }

    apply_flag("--finished", pane);
    refresh_window_name(pane);
    play_finish_sound("--finished");
    if finish_timer(pane) == *token {
        clear_finish_timer(pane);
    }
    0
}

// zsb_tmux_agent_notification [--finished|--force-finished|--working|--clear-finished|--manual]
//     <session_name> <pane_id>
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
    if args.first().map(String::as_str) == Some(DEBOUNCED_FINISH) {
        return run_debounced_finish(args);
    }

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
    let resolved_pane = pane_id(pane);

    if is_finished(flag) && !resolved_pane.is_empty() {
        if schedule_finish(&resolved_pane) {
            return 0;
        }
    } else if cancels_finish(flag) && !resolved_pane.is_empty() {
        clear_finish_timer(&resolved_pane);
    }

    if !apply_flag(flag, pane) {
        return 1;
    }
    refresh_window_name(pane);
    play_finish_sound(flag);
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
            (
                "finished only",
                format!("myrepo{FINISHED_SUFFIX}"),
                "myrepo",
            ),
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

    #[test]
    fn only_normal_finish_is_debounced() {
        assert!(is_finished(""));
        assert!(is_finished("--finished"));
        assert!(!is_finished("--force-finished"));
        assert!(!is_finished("--working"));
        assert!(!is_finished("--clear-finished"));
        assert!(!is_finished("--manual"));

        assert!(cancels_finish("--force-finished"));
        assert!(cancels_finish("--working"));
        assert!(!cancels_finish("--clear-finished"));
        assert!(cancels_finish("--manual"));
        assert!(cancels_finish("--invalid"));
    }
}
