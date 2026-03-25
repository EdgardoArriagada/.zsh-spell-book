use std::io::{self, BufRead};
use std::process;
use std::sync::LazyLock;

use regex::Regex;

static URL_RE: LazyLock<Regex> =
    LazyLock::new(|| Regex::new(r"https?://[\w\-_\.%?/:+=&#%]+").unwrap());

fn extract_first_url(text: &str) -> Option<&str> {
    URL_RE.find(text).map(|m| m.as_str())
}

fn main() {
    let input = if let Some(arg) = std::env::args().nth(1) {
        arg
    } else {
        let stdin = io::stdin();
        let mut line = String::new();
        stdin.lock().read_line(&mut line).unwrap_or(0);
        line.trim_end_matches('\n').to_owned()
    };

    if input.is_empty() {
        eprintln!("Usage: zsb_charm_tmux_urlopen <url>");
        process::exit(1);
    }

    match extract_first_url(&input) {
        Some(url) => {
            if let Err(e) = open::that(url) {
                eprintln!("Error: failed to open URL: {}", e);
                process::exit(1);
            }
        }
        None => {
            eprintln!("Error: No valid URL found in the input.");
            process::exit(1);
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_extract_first_url() {
        assert_eq!(
            extract_first_url("Check this out: http://example.com"),
            Some("http://example.com")
        );
        assert_eq!(
            extract_first_url("Visit https://example.com for more info"),
            Some("https://example.com")
        );
        assert_eq!(extract_first_url("No URL here!"), None);
        assert_eq!(
            extract_first_url("Multiple URLs: http://example.com and https://example.org"),
            Some("http://example.com")
        );
        assert_eq!(
            extract_first_url("Special chars: https://example.com/path?query=1&other=2"),
            Some("https://example.com/path?query=1&other=2")
        );
        assert_eq!(extract_first_url("Edge case: just some text"), None);
        assert_eq!(
            extract_first_url("Another edge case: https://example.com/path#fragment"),
            Some("https://example.com/path#fragment")
        );
        assert_eq!(
            extract_first_url(
                "Another edge case: (https://www.example.cl/vehiculos/autos-veh%C3%ADculo/foo/)"
            ),
            Some("https://www.example.cl/vehiculos/autos-veh%C3%ADculo/foo/")
        );
        assert_eq!(extract_first_url("Empty string: "), None);
        assert_eq!(extract_first_url("Only special chars: !@#$%^&*()"), None);
    }
}
