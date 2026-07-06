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
