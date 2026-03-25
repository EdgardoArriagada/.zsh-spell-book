use std::process::{self, Command};

fn ssh_to_https(url: &str) -> String {
    // e.g. git@github.com:user/repo.git -> https://github.com/user/repo
    let without_prefix = url.strip_prefix("git@").unwrap_or(url);
    let normalized = without_prefix.replacen(':', "/", 1);
    let without_git = normalized.strip_suffix(".git").unwrap_or(&normalized);
    format!("https://{}", without_git)
}

fn main() {
    let output = Command::new("git")
        .args(["config", "--get", "remote.origin.url"])
        .stderr(process::Stdio::null())
        .output()
        .unwrap_or_else(|_| process::exit(1));

    if !output.status.success() {
        process::exit(1);
    }

    let url = String::from_utf8_lossy(&output.stdout).trim().to_owned();

    if url.starts_with("http://") || url.starts_with("https://") {
        print!("{}", url);
    } else {
        print!("{}", ssh_to_https(&url));
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn ssh_standard() {
        assert_eq!(
            ssh_to_https("git@github.com:user/repo.git"),
            "https://github.com/user/repo"
        );
    }

    #[test]
    fn ssh_no_git_suffix() {
        assert_eq!(
            ssh_to_https("git@github.com:user/repo"),
            "https://github.com/user/repo"
        );
    }

    #[test]
    fn ssh_nested_path() {
        assert_eq!(
            ssh_to_https("git@github.com:org/suborg/repo.git"),
            "https://github.com/org/suborg/repo"
        );
    }
}
