//! zsb_tmux_jira_picker
//!
//! Reads `~/temp/tickets` (pipe-delimited: parent_ticket|current_ticket|label),
//! presents an fzf picker, then for a NEW session:
//!   - Creates the tmux session
//!   - Sets ZSB_PARENT_TICKET, ZSB_CURRENT_TICKET, ZSB_CURRENT_LABEL in the session
//!   - Writes ~/temp/current-ticket.zsh with those vars (sourced by ~/.zshenv
//!     so new shells inherit the current ticket context)
//! For an EXISTING session: switches to it, does nothing else.

use std::env;
use std::fs;
use std::io::Write;
use std::process::{self, Command, Stdio};

fn tmux(args: &[&str]) {
    Command::new("tmux").args(args).status().ok();
}

fn main() {
    let home = env::var("HOME").unwrap_or_default();
    let tickets_path = format!("{}/temp/tickets", home);

    let content = fs::read_to_string(&tickets_path).unwrap_or_else(|_| {
        eprintln!("zsb_tmux_jira_picker: no tickets file at {}", tickets_path);
        process::exit(1);
    });

    if content.trim().is_empty() {
        eprintln!("zsb_tmux_jira_picker: tickets file is empty");
        process::exit(1);
    }

    let mut child = Command::new("fzf")
        .args(["--delimiter=|", "--with-nth=2..", "--prompt=ticket> "])
        .stdin(Stdio::piped())
        .stdout(Stdio::piped())
        .spawn()
        .unwrap_or_else(|_| {
            eprintln!("zsb_tmux_jira_picker: failed to spawn fzf");
            process::exit(1);
        });

    if let Some(mut stdin) = child.stdin.take() {
        let _ = stdin.write_all(content.as_bytes());
    }

    let output = child.wait_with_output().unwrap_or_else(|_| {
        eprintln!("zsb_tmux_jira_picker: fzf failed");
        process::exit(1);
    });

    if output.stdout.is_empty() {
        return;
    }

    let line = String::from_utf8_lossy(&output.stdout);
    let line = line.trim();

    let parts: Vec<&str> = line.splitn(3, '|').collect();
    if parts.len() < 3 {
        eprintln!("zsb_tmux_jira_picker: invalid format: {}", line);
        process::exit(1);
    }

    let parent_ticket = parts[0].trim();
    let current_ticket = parts[1].trim();
    let current_label = parts[2].trim();

    let label_kebab = current_label.to_lowercase().replace(' ', "-");
    let truncate_at = label_kebab.len().min(25);
    let label_short = &label_kebab[..truncate_at];
    let session_name = format!("{}-{}", current_ticket, label_short);

    let session_exists = Command::new("tmux")
        .args(["has-session", "-t", &format!("={}", session_name)])
        .stderr(Stdio::null())
        .status()
        .map(|s| s.success())
        .unwrap_or(false);

    if !session_exists {
        Command::new("tmux")
            .args(["new-session", "-d", "-s", &session_name])
            .status()
            .ok();
        tmux(&["set-environment", "-t", &session_name, "ZSB_PARENT_TICKET", parent_ticket]);
        tmux(&["set-environment", "-t", &session_name, "ZSB_CURRENT_TICKET", current_ticket]);
        tmux(&["set-environment", "-t", &session_name, "ZSB_CURRENT_LABEL", current_label]);

        let ticket_file = format!("{}/temp/current-ticket.zsh", home);
        let file_content = format!(
            "declare ZSB_PARENT_TICKET={}\ndeclare ZSB_CURRENT_TICKET={}\ndeclare ZSB_CURRENT_LABEL={}\n",
            parent_ticket, current_ticket, current_label
        );
        let _ = fs::write(&ticket_file, file_content);
    }

    if env::var("TMUX").is_ok() {
        tmux(&["switch-client", "-t", &format!("={}", session_name)]);
    } else {
        tmux(&["attach-session", "-t", &format!("={}", session_name)]);
    }
}
