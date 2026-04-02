use std::process::{self, Command};

fn tmux_yank_args(key: &str) -> Option<Vec<&str>> {
    match key {
        "w" | "e" => Some(vec![
            "send", "-X", "select-word", ";",
            "send", "-X", "other-end", ";",
            "send", "-X", "copy-selection",
        ]),
        "y" => Some(vec!["send", "-X", "copy-line"]),
        "b" => Some(vec![
            "send", "-X", "begin-selection", ";",
            "send", "-X", "previous-word", ";",
            "send", "-X", "copy-selection",
        ]),
        "$" => Some(vec!["send", "-X", "copy-end-of-line"]),
        "q" => None,
        _ => {
            process::exit(1);
        }
    }
}

fn main() {
    let args: Vec<String> = std::env::args().collect();

    if args.len() == 2 {
        if let Some(tmux_args) = tmux_yank_args(&args[1]) {
            Command::new("tmux")
                .args(&tmux_args)
                .status()
                .ok();
        }
        return;
    }

    let bin = &args[0];
    Command::new("tmux")
        .args([
            "command-prompt",
            "-1",
            "-p",
            "(w/e/b/y/$/q)",
            &format!("run-shell '{} \"%%\"'", bin),
        ])
        .status()
        .ok();
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn word_selection() {
        let args = tmux_yank_args("w").unwrap();
        assert_eq!(
            args,
            vec![
                "send", "-X", "select-word", ";",
                "send", "-X", "other-end", ";",
                "send", "-X", "copy-selection",
            ]
        );
        assert_eq!(tmux_yank_args("e"), tmux_yank_args("w"));
    }

    #[test]
    fn copy_line() {
        let args = tmux_yank_args("y").unwrap();
        assert_eq!(args, vec!["send", "-X", "copy-line"]);
    }

    #[test]
    fn backward_word() {
        let args = tmux_yank_args("b").unwrap();
        assert_eq!(
            args,
            vec![
                "send", "-X", "begin-selection", ";",
                "send", "-X", "previous-word", ";",
                "send", "-X", "copy-selection",
            ]
        );
    }

    #[test]
    fn end_of_line() {
        let args = tmux_yank_args("$").unwrap();
        assert_eq!(args, vec!["send", "-X", "copy-end-of-line"]);
    }

    #[test]
    fn quit_returns_none() {
        assert!(tmux_yank_args("q").is_none());
    }
}
