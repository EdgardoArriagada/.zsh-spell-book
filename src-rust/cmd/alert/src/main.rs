use clap::Parser;
use std::process::{Command, Stdio};

#[derive(Debug, Parser)]
#[command(about = "Cross-platform desktop notification")]
struct Args {
    /// Icon name (Linux only, notify-send). Defaults to "timer"
    #[clap(short, long, default_value = "timer")]
    icon: String,

    /// Message to display (all remaining args joined with spaces)
    #[clap(required = true)]
    message: Vec<String>,
}

#[cfg(target_os = "macos")]
fn notify(message: &str, _icon: &str) {
    let script = format!(
        r#"tell app "System Events" to display dialog "{}" buttons {{"OK"}} default button "OK""#,
        message.replace('"', "\\\"")
    );
    Command::new("osascript")
        .args(["-e", &script])
        .stdout(Stdio::null())
        .stderr(Stdio::null())
        .status()
        .ok();
}

#[cfg(not(target_os = "macos"))]
fn notify(message: &str, icon: &str) {
    Command::new("notify-send")
        .args(["--icon", icon, "--urgency", "critical", message])
        .status()
        .ok();
}

fn main() {
    let args = Args::parse();
    let message = args.message.join(" ");
    notify(&message, &args.icon);
}
