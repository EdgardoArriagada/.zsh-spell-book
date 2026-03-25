use std::path::Path;
use std::process::{self, Command};

fn get_repo_root() -> Option<String> {
    let output = Command::new("git")
        .args(["rev-parse", "--show-toplevel"])
        .stderr(process::Stdio::null())
        .output()
        .ok()?;
    if output.status.success() {
        Some(String::from_utf8_lossy(&output.stdout).trim().to_owned())
    } else {
        None
    }
}

fn main() {
    let args: Vec<String> = std::env::args().collect();
    if args.len() != 3 {
        eprintln!("Usage: zsb_charm_tmux_renametab <win_id> <current_file_path>");
        process::exit(1);
    }

    let win_id = &args[1];
    let file_path = &args[2];

    let repo_root = match get_repo_root() {
        Some(r) => r,
        None => process::exit(1),
    };

    if !file_path.is_empty() && !file_path.starts_with(&repo_root) {
        process::exit(1);
    }

    let basename = Path::new(&repo_root)
        .file_name()
        .and_then(|n| n.to_str())
        .unwrap_or("");

    let tabname = format!("  {}", basename);

    Command::new("tmux")
        .args(["rename-window", "-t", win_id, &tabname])
        .status()
        .ok();
}

#[cfg(test)]
mod tests {
    fn file_in_repo(file_path: &str, repo_root: &str) -> bool {
        file_path.is_empty() || file_path.starts_with(repo_root)
    }

    #[test]
    fn empty_file_path_is_allowed() {
        assert!(file_in_repo("", "/home/user/myrepo"));
    }

    #[test]
    fn file_inside_repo_is_allowed() {
        assert!(file_in_repo(
            "/home/user/myrepo/src/main.rs",
            "/home/user/myrepo"
        ));
    }

    #[test]
    fn file_outside_repo_is_rejected() {
        assert!(!file_in_repo(
            "/home/user/otherrepo/main.rs",
            "/home/user/myrepo"
        ));
    }

    #[test]
    fn repo_root_itself_is_allowed() {
        assert!(file_in_repo("/home/user/myrepo", "/home/user/myrepo"));
    }
}
