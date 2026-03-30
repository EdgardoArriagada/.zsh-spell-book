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
        self.content
            .replace("${zsb}", zsb)
            .replace("$ZSB_DIR", zsb_dir)
            .replace("$ZSB_TEMP_DIR", zsb_temp_dir)
    }
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
