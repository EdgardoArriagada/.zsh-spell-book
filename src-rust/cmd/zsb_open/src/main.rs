use std::io::{self, BufRead};
use std::process::Command;

#[cfg(target_os = "macos")]
fn open_cmd(target: &str) -> Command {
    let mut cmd = Command::new("open");
    cmd.arg(target);
    cmd
}

#[cfg(not(target_os = "macos"))]
fn open_cmd(target: &str) -> Command {
    let mut cmd = Command::new("xdg-open");
    cmd.arg(target);
    cmd
}

fn main() {
    let target = std::env::args().nth(1).unwrap_or_else(|| {
        io::stdin()
            .lock()
            .lines()
            .next()
            .expect("no input")
            .expect("failed to read stdin")
    });

    open_cmd(&target).status().expect("failed to open");
}
