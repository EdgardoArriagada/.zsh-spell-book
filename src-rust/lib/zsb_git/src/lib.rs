use std::process::Command;

/// Repo name for `cwd`. Prefers the `origin` remote URL; falls back to the
/// main worktree's directory name for local-only repos.
pub fn git_repo_name(cwd: &str) -> Option<String> {
    remote_repo_name(cwd).or_else(|| local_repo_name(cwd))
}

fn git_out(cwd: &str, args: &[&str]) -> Option<String> {
    let out = Command::new("git")
        .args(["-C", cwd])
        .args(args)
        .stderr(std::process::Stdio::null())
        .output()
        .ok()?;
    if !out.status.success() {
        return None;
    }
    let s = String::from_utf8_lossy(&out.stdout).trim().to_string();
    if s.is_empty() {
        return None;
    }
    Some(s)
}

fn remote_repo_name(cwd: &str) -> Option<String> {
    let url = git_out(cwd, &["remote", "get-url", "origin"])?;
    Some(repo_name_from_url(&url).to_string())
}

// `git worktree list` works from any subdir and always lists the main worktree
// first, so this handles both deep cwd's and secondary worktrees.
fn local_repo_name(cwd: &str) -> Option<String> {
    let out = git_out(cwd, &["worktree", "list", "--porcelain"])?;
    repo_name_from_worktree_list(&out).map(str::to_string)
}

pub fn repo_name_from_url(url: &str) -> &str {
    let base = url.rsplit(|c| c == '/' || c == ':').next().unwrap_or(url);
    base.strip_suffix(".git").unwrap_or(base)
}

pub fn repo_name_from_worktree_list(out: &str) -> Option<&str> {
    let path = out.lines().next()?.strip_prefix("worktree ")?.trim_end();
    let name = path.trim_end_matches('/').rsplit('/').next()?;
    if name.is_empty() {
        None
    } else {
        Some(name)
    }
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

    #[test]
    fn worktree_list_main_repo() {
        let out = "worktree /Users/me/code/myrepo\nHEAD abc123\nbranch refs/heads/master\n";
        assert_eq!(repo_name_from_worktree_list(out), Some("myrepo"));
    }

    #[test]
    fn worktree_list_from_secondary_worktree() {
        // main worktree is always listed first, even when cwd is a worktree
        let out = "worktree /Users/me/code/myrepo\nHEAD abc\nbranch refs/heads/master\n\
                   \n\
                   worktree /Users/me/code/myrepo_gitworktree/feature/x\nHEAD def\nbranch refs/heads/feature/x\n";
        assert_eq!(repo_name_from_worktree_list(out), Some("myrepo"));
    }

    #[test]
    fn worktree_list_bare_repo() {
        let out = "worktree /Users/me/code/myrepo\nbare\n";
        assert_eq!(repo_name_from_worktree_list(out), Some("myrepo"));
    }

    #[test]
    fn worktree_list_garbage() {
        assert_eq!(repo_name_from_worktree_list(""), None);
        assert_eq!(repo_name_from_worktree_list("not a worktree line"), None);
    }
}
