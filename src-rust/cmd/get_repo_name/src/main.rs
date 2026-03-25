use std::process::{self, Command};

fn repo_name_from_url(url: &str) -> &str {
    let base = url.rsplit(|c| c == '/' || c == ':').next().unwrap_or(url);
    base.strip_suffix(".git").unwrap_or(base)
}

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

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn https_with_git_suffix() {
        assert_eq!(repo_name_from_url("https://github.com/user/repo.git"), "repo");
    }

    #[test]
    fn https_without_git_suffix() {
        assert_eq!(repo_name_from_url("https://github.com/user/repo"), "repo");
    }

    #[test]
    fn ssh_url() {
        assert_eq!(repo_name_from_url("git@github.com:user/repo.git"), "repo");
    }

    #[test]
    fn ssh_url_no_suffix() {
        assert_eq!(repo_name_from_url("git@github.com:user/repo"), "repo");
    }
}
