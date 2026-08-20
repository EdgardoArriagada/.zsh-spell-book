use std::{
    env, io,
    process::{Child, Command, Stdio},
};

fn player_for_env(zsb_macos: Option<&str>) -> &'static str {
    match zsb_macos {
        Some(value) if value.trim().parse::<i32>().unwrap_or(0) != 0 => "afplay",
        Some(_) => "aplay",
        None if cfg!(target_os = "macos") => "afplay",
        None => "aplay",
    }
}

pub fn spawn(file: &str) -> io::Result<Child> {
    Command::new(player_for_env(env::var("ZSB_MACOS").ok().as_deref()))
        .arg(file)
        .stdout(Stdio::null())
        .stderr(Stdio::null())
        .spawn()
}

pub fn play(file: &str) -> i32 {
    match spawn(file).and_then(|mut child| child.wait()) {
        Ok(status) => status.code().unwrap_or(1),
        Err(_) => 127,
    }
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
