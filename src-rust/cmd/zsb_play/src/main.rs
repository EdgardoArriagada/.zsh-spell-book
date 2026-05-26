use std::env;
use std::process::{self, Command, Stdio};

fn player_for_env(zsb_macos: Option<&str>) -> &'static str {
    match zsb_macos {
        Some(value) if value.trim().parse::<i32>().unwrap_or(0) != 0 => "afplay",
        Some(_) => "aplay",
        None if cfg!(target_os = "macos") => "afplay",
        None => "aplay",
    }
}

fn run_player(player: &str, file: &str) -> i32 {
    match Command::new(player)
        .arg(file)
        .stdout(Stdio::null())
        .stderr(Stdio::null())
        .status()
    {
        Ok(status) => status.code().unwrap_or(1),
        Err(_) => 127,
    }
}

fn main() {
    let mut args = env::args().skip(1);
    let Some(file) = args.next() else {
        eprintln!("Usage: zsb_play <sound-file>");
        process::exit(1);
    };

    let player = player_for_env(env::var("ZSB_MACOS").ok().as_deref());
    process::exit(run_player(player, &file));
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn zsb_macos_one_uses_afplay() {
        assert_eq!(player_for_env(Some("1")), "afplay");
    }

    #[test]
    fn zsb_macos_zero_uses_aplay() {
        assert_eq!(player_for_env(Some("0")), "aplay");
    }

    #[test]
    fn zsb_macos_unset_uses_platform_default() {
        let expected = if cfg!(target_os = "macos") {
            "afplay"
        } else {
            "aplay"
        };

        assert_eq!(player_for_env(None), expected);
    }
}
