use std::env;
use std::fs;
use std::io::ErrorKind;
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
        let data = fs::read_to_string(&full_path)
            .unwrap_or_else(|e| panic!("failed to read {}: {}", full_path.display(), e));
        self.content.push_str(&data);
        self.content.push('\n');
    }

    fn load_file_optional(&mut self, relative_path: &str) {
        let full_path = self.zsb_dir.join(relative_path);
        match fs::read_to_string(&full_path) {
            Ok(data) => {
                self.content.push_str(&data);
                self.content.push('\n');
            }
            Err(e) if e.kind() == ErrorKind::NotFound => {}
            Err(e) => panic!("failed to read {}: {}", full_path.display(), e),
        }
    }

    fn load_dir(&mut self, relative_path: &str) {
        let full_path = self.zsb_dir.join(relative_path);
        self.walk_dir(&full_path, false);
    }

    fn load_dir_optional(&mut self, relative_path: &str) {
        let full_path = self.zsb_dir.join(relative_path);
        self.walk_dir(&full_path, true);
    }

    fn walk_dir(&mut self, dir: &Path, optional: bool) {
        let entries = match fs::read_dir(dir) {
            Ok(entries) => entries,
            Err(e) if optional && e.kind() == ErrorKind::NotFound => return,
            Err(e) => panic!("failed to read dir {}: {}", dir.display(), e),
        };

        let mut paths: Vec<PathBuf> = entries
            .map(|e| {
                e.unwrap_or_else(|err| {
                    panic!("failed to read entry in {}: {}", dir.display(), err)
                })
                .path()
            })
            .collect();
        paths.sort();

        for path in paths {
            if path.is_dir() {
                self.walk_dir(&path, false);
            } else if path.extension().and_then(|e| e.to_str()) == Some("zsh") {
                let data = fs::read_to_string(&path)
                    .unwrap_or_else(|e| panic!("failed to read {}: {}", path.display(), e));
                self.content.push_str(&data);
                self.content.push('\n');
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
    // Pass 1: collect all hisIgnore patterns.
    let mut patterns: Vec<String> = Vec::new();
    for line in content.lines() {
        if let Some(args) = line.strip_prefix("hisIgnore ") {
            patterns.extend(parse_hisignore_args(args));
        }
    }
    let static_export = format!("export HISTORY_IGNORE=\"({})\"", patterns.join("|"));

    // Pass 2: strip dynamic machinery and replace the export line.
    let mut out = String::with_capacity(content.len());
    for line in content.lines() {
        if line.starts_with("declare ZSB_HISTORY_IGNORE=()")
            || line.starts_with("hisIgnore()")
            || line.starts_with("hisIgnore ")
        {
            continue;
        }
        if line.starts_with("export HISTORY_IGNORE=\"(${") {
            out.push_str(&static_export);
            out.push('\n');
            continue;
        }
        out.push_str(line);
        out.push('\n');
    }
    out
}

fn main() {
    let home = env::var("HOME").expect("HOME not set");
    let zsb = env::var("zsb").unwrap_or_else(|_| "zsb".to_string());
    let zsb_dir = format!("{home}/.zsh-spell-book");
    let zsb_temp_dir = format!("{zsb_dir}/src/temp");
    let result_filepath = format!("{zsb_dir}/result.zsh");

    let mut bundler = Bundler::new(PathBuf::from(&zsb_dir));

    bundler.load_file_optional(".env");
    bundler.load_file("src/zsh.config.zsh");
    bundler.load_file("src/globalVariables.zsh");

    bundler.load_dir("src/utils");
    bundler.load_dir("src/configurations");
    bundler.load_dir("src/spells");
    bundler.load_dir_optional("src/temp/spells");
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
declare ZSB_HISTORY_IGNORE=() # zsb_bundle: exact match, do not modify

hisIgnore() ZSB_HISTORY_IGNORE+=( $@ ) # zsb_bundle: exact match, do not modify

hisIgnore 'l[a,l,s,h,]*' 'neofetch'
hisIgnore gs pop
hisIgnore 'ga .'

export HISTORY_IGNORE="(${(j:|:)ZSB_HISTORY_IGNORE})" # zsb_bundle: exact match, do not modify
"#;
        let result = compile_history_ignore(input.to_string());

        assert!(!result.contains("ZSB_HISTORY_IGNORE"));
        assert!(!result.contains("hisIgnore"));
        assert!(result.contains(
            r#"export HISTORY_IGNORE="(l[a,l,s,h,]*|neofetch|gs|pop|ga .)""#
        ));
    }
}
