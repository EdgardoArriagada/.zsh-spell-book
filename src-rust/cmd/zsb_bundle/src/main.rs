use regex::Regex;
use std::env;
use std::fs;
use std::path::{Path, PathBuf};
use std::process::Command;

struct Bundler {
    content: String,
    zsb_dir: PathBuf,
}

impl Bundler {
    fn new(zsb_dir: PathBuf) -> Self {
        Self {
            content: String::new(),
            zsb_dir,
        }
    }

    fn load_file(&mut self, relative_path: &str) {
        let full_path = self.zsb_dir.join(relative_path);
        if let Ok(data) = fs::read_to_string(&full_path) {
            self.content.push_str(&data);
            self.content.push('\n');
        }
    }

    fn load_dir(&mut self, relative_path: &str) {
        let full_path = self.zsb_dir.join(relative_path);
        self.walk_dir(&full_path);
    }

    fn walk_dir(&mut self, dir: &Path) {
        let entries = match fs::read_dir(dir) {
            Ok(entries) => entries,
            Err(_) => return,
        };

        let mut paths: Vec<PathBuf> = entries.filter_map(|e| e.ok().map(|e| e.path())).collect();
        paths.sort();

        for path in paths {
            if path.is_dir() {
                self.walk_dir(&path);
            } else if path.extension().and_then(|e| e.to_str()) == Some("zsh") {
                if let Ok(data) = fs::read_to_string(&path) {
                    self.content.push_str(&data);
                    self.content.push('\n');
                }
            }
        }
    }

    fn bundle(self, zsb: &str, zsb_dir: &str, zsb_temp_dir: &str) -> String {
        let substituted = self
            .content
            .replace("${zsb}", zsb)
            .replace("$ZSB_DIR", zsb_dir)
            .replace("$ZSB_TEMP_DIR", zsb_temp_dir);

        compile_history_ignore(substituted)
    }
}

/// Parse arguments from a `hisIgnore` call, handling both quoted and unquoted args.
fn parse_hisignore_args(args_str: &str) -> Vec<String> {
    let mut result = Vec::new();
    let mut chars = args_str.chars().peekable();

    while let Some(&c) = chars.peek() {
        if c.is_whitespace() {
            chars.next();
            continue;
        }
        if c == '\'' {
            chars.next(); // consume opening quote
            let mut arg = String::new();
            for ch in chars.by_ref() {
                if ch == '\'' {
                    break;
                }
                arg.push(ch);
            }
            result.push(arg);
        } else {
            let mut arg = String::new();
            while let Some(&ch) = chars.peek() {
                if ch.is_whitespace() {
                    break;
                }
                arg.push(ch);
                chars.next();
            }
            result.push(arg);
        }
    }

    result
}

/// Collect all HISTORY_IGNORE patterns, compute the final value, and strip
/// the dynamic machinery from the bundled output.
fn compile_history_ignore(content: String) -> String {
    let mut patterns: Vec<String> = Vec::new();

    // 1. Extract seed patterns from `declare ZSB_HISTORY_IGNORE=(...)`
    let declare_re =
        Regex::new(r"(?ms)^declare ZSB_HISTORY_IGNORE=\(\s*\n(.*?)\)").unwrap();
    if let Some(caps) = declare_re.captures(&content) {
        let block = &caps[1];
        for line in block.lines() {
            let trimmed = line.trim().trim_matches('\'');
            if !trimmed.is_empty() {
                patterns.push(trimmed.to_string());
            }
        }
    }

    // 2. Extract patterns from all `hisIgnore ...` calls
    let call_re = Regex::new(r"(?m)^hisIgnore (.+)$").unwrap();
    for caps in call_re.captures_iter(&content) {
        let args = parse_hisignore_args(&caps[1]);
        patterns.extend(args);
    }

    // 3. Build the static export line
    let joined = patterns.join("|");
    let static_export = format!("export HISTORY_IGNORE=\"({joined})\"");

    // 4. Strip dynamic machinery
    let mut result = declare_re.replace(&content, "").to_string();

    let func_re = Regex::new(r"(?m)^hisIgnore\(\).*\n").unwrap();
    result = func_re.replace(&result, "").to_string();

    result = call_re.replace_all(&result, "").to_string();

    let export_re =
        Regex::new(r#"(?m)^export HISTORY_IGNORE="\(\$\{.*\}\)".*$"#).unwrap();
    result = export_re.replace(&result, &*static_export).to_string();

    result
}

fn main() {
    let home = env::var("HOME").expect("HOME not set");
    let zsb = env::var("zsb").unwrap_or_else(|_| "zsb".to_string());
    let zsb_dir = format!("{home}/.zsh-spell-book");
    let zsb_temp_dir = format!("{zsb_dir}/src/temp");
    let result_filepath = format!("{zsb_dir}/result.zsh");

    let mut bundler = Bundler::new(PathBuf::from(&zsb_dir));

    bundler.load_file(".env");
    bundler.load_file("src/zsh.config.zsh");
    bundler.load_file("src/globalVariables.zsh");

    bundler.load_dir("src/utils");
    bundler.load_dir("src/configurations");
    bundler.load_dir("src/spells");
    bundler.load_dir("src/temp/spells");
    bundler.load_dir("src/automatic-calls");

    let result = bundler.bundle(&zsb, &zsb_dir, &zsb_temp_dir);

    fs::write(&result_filepath, result).expect("failed to write result.zsh");
    let _ = Command::new("zsh")
        .arg("-c")
        .arg(format!("zcompile {result_filepath}"))
        .status();

    println!("ℨ𝔰𝔟 𝔖𝔭𝔢𝔩𝔩𝔟𝔬𝔬𝔨 bundled!!");
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_parse_hisignore_args_unquoted() {
        let args = parse_hisignore_args("gs pop amend glog gsp");
        assert_eq!(args, vec!["gs", "pop", "amend", "glog", "gsp"]);
    }

    #[test]
    fn test_parse_hisignore_args_quoted() {
        let args = parse_hisignore_args("'gsw -' 'gsw 0' 'gsw ='");
        assert_eq!(args, vec!["gsw -", "gsw 0", "gsw ="]);
    }

    #[test]
    fn test_parse_hisignore_args_mixed() {
        let args = parse_hisignore_args("ga 'ga .'");
        assert_eq!(args, vec!["ga", "ga ."]);
    }

    #[test]
    fn test_compile_history_ignore() {
        let input = r#"# some config
declare ZSB_HISTORY_IGNORE=(
  'l[a,l,s,h,]*'
  'neofetch'
)

hisIgnore() ZSB_HISTORY_IGNORE+=( $@ )

hisIgnore gs pop
hisIgnore 'ga .'

export HISTORY_IGNORE="(${(j:|:)ZSB_HISTORY_IGNORE})"
"#;
        let result = compile_history_ignore(input.to_string());

        assert!(!result.contains("ZSB_HISTORY_IGNORE"));
        assert!(!result.contains("hisIgnore"));
        assert!(result.contains(
            r#"export HISTORY_IGNORE="(l[a,l,s,h,]*|neofetch|gs|pop|ga .)""#
        ));
    }
}
