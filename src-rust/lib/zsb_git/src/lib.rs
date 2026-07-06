pub fn git_repo_name(cwd: &str) -> Option<String> {
    use std::process::Command;
    let out = Command::new("git")
        .args(["-C", cwd, "remote", "get-url", "origin"])
        .stderr(std::process::Stdio::null())
        .output()
        .ok()?;
    if !out.status.success() {
        return None;
    }
    let url = String::from_utf8_lossy(&out.stdout).trim().to_string();
    if url.is_empty() {
        return None;
    }
    Some(repo_name_from_url(&url).to_string())
}

pub fn repo_name_from_url(url: &str) -> &str {
    let base = url.rsplit(|c| c == '/' || c == ':').next().unwrap_or(url);
    base.strip_suffix(".git").unwrap_or(base)
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
