use std::process;

use zsb_git::git_repo_name;

fn main() {
    match git_repo_name(".") {
        Some(name) => println!("{}", name),
        None => process::exit(1),
    }
}

