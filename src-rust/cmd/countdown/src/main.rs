use std::io::{self, Write};
use std::process::Command;
use std::{env, thread, time::Duration};

use clap::Parser;

const MAX_SECONDS: u32 = 215999; // 59:59:59

#[derive(Debug, Parser)]
#[command(about = "Countdown timer with notification")]
struct Args {
    /// Time duration (e.g., 30, 5m, 1h, 05:30, 1:05:30)
    time: String,

    /// Custom alert message
    #[clap(short, long)]
    message: Option<String>,
}

fn main() {
    let args = Args::parse();

    let seconds = zsb_time::parse_duration(&args.time).unwrap_or_else(|e| {
        eprintln!(
            "{}\n  Examples:\n    countdown 30      (30 seconds)\n    countdown 5m      (5 minutes)\n    countdown 1h      (1 hour)\n    countdown 05:30   (5 minutes 30 seconds)\n    countdown 1:05:30 (1 hour 5 minutes 30 seconds)",
            e
        );
        std::process::exit(1);
    });

    if seconds == 0 || seconds > MAX_SECONDS {
        eprintln!(
            "Bad argument.\nTry with (hh:)?mm:ss (min 1, max 59:59:59)\nOR {{s : s ∈ Z and 1 ≤ s ≤ {}}}\nOR ^n[hHmMsS]$ (min 1s, max 59h)",
            MAX_SECONDS
        );
        std::process::exit(1);
    }

    run_countdown(seconds);
    play_notification_sound();

    let time_message = zsb_time::human_duration(seconds);
    let current_time = chrono_free_time();
    println!("The timer for {} was up at {}", time_message, current_time);

    let notification_msg = match &args.message {
        Some(m) => m.clone(),
        None => format!("The timer for {} is over", time_message),
    };
    show_notification(&notification_msg);
}

fn run_countdown(total_seconds: u32) {
    let mut stdout = io::stdout();
    for i in (1..=total_seconds).rev() {
        let time_str = zsb_time::format_clock(i);
        print!("\r⏳ {}   ", time_str);
        stdout.flush().ok();
        thread::sleep(Duration::from_secs(1));
    }
    print!("\r");
    stdout.flush().ok();
}

fn play_notification_sound() {
    let home = env::var("HOME").unwrap_or_default();
    let sound_file = format!("{}/.zsh-spell-book/src/media/sounds/xylofon.wav", home);
    Command::new("afplay").arg(&sound_file).spawn().ok();
}

fn chrono_free_time() -> String {
    let output = Command::new("date").arg("+%H:%M:%S").output();
    match output {
        Ok(o) => String::from_utf8_lossy(&o.stdout).trim().to_string(),
        Err(_) => String::from("??:??:??"),
    }
}

fn show_notification(message: &str) {
    let bin_dir = env::current_exe()
        .ok()
        .and_then(|p| p.parent().map(|d| d.to_path_buf()));

    let alert_path = bin_dir
        .map(|d| d.join("alert"))
        .filter(|p| p.exists())
        .map(|p| p.to_string_lossy().to_string())
        .unwrap_or_else(|| "alert".to_string());

    Command::new(&alert_path).arg(message).status().ok();
}

#[cfg(test)]
mod tests {
    #[test]
    fn test_notification_message_custom() {
        let custom = Some("Take a break!".to_string());
        let msg = match &custom {
            Some(m) => m.clone(),
            None => format!("The timer for {} is over", zsb_time::human_duration(30)),
        };
        assert_eq!(msg, "Take a break!");
    }

    #[test]
    fn test_notification_message_default() {
        let custom: Option<String> = None;
        let msg = match &custom {
            Some(m) => m.clone(),
            None => format!("The timer for {} is over", zsb_time::human_duration(30)),
        };
        assert_eq!(msg, "The timer for 30 seconds is over");
    }

    #[test]
    fn test_notification_message_default_clock() {
        let custom: Option<String> = None;
        let msg = match &custom {
            Some(m) => m.clone(),
            None => format!("The timer for {} is over", zsb_time::human_duration(90)),
        };
        assert_eq!(msg, "The timer for 01:30 is over");
    }
}
