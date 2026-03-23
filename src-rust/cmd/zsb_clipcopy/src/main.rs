use std::io::{self, Read, Write};
use std::process::{Command, Stdio};

#[cfg(target_os = "macos")]
fn clipboard_cmd() -> Command {
    Command::new("pbcopy")
}

#[cfg(not(target_os = "macos"))]
fn clipboard_cmd() -> Command {
    let mut cmd = Command::new("xclip");
    cmd.args(["-selection", "clipboard"]);
    cmd
}

fn main() {
    let args: Vec<String> = std::env::args().skip(1).collect();
    let content = if args.is_empty() {
        let mut buf = String::new();
        io::stdin().read_to_string(&mut buf).unwrap();
        buf
    } else {
        args.join("\n")
    };

    let mut child = clipboard_cmd()
        .stdin(Stdio::piped())
        .spawn()
        .expect("failed to spawn clipboard command");
    child.stdin.take().unwrap().write_all(content.as_bytes()).unwrap();
    child.wait().unwrap();
}
