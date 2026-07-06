use std::process::{self, Command};

use zsb_git::repo_name_from_url;

fn main() {
    let output = Command::new("git")
        .args(["remote", "get-url", "origin"])
        .stderr(process::Stdio::null())
        .output()
        .unwrap_or_else(|_| process::exit(1));

    if !output.status.success() {
        process::exit(1);
    }

    let url = String::from_utf8_lossy(&output.stdout).trim().to_owned();
    println!("{}", repo_name_from_url(&url));
}

