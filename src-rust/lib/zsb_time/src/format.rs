pub fn format_clock(seconds: u32) -> String {
    let hours = seconds / 3600;
    let minutes = seconds / 60 % 60;
    let secs = seconds % 60;

    match (hours, minutes) {
        (0, 0) => format!("{}", secs),
        (0, _) => format!("{:02}:{:02}", minutes, secs),
        _ => format!("{:02}:{:02}:{:02}", hours, minutes, secs),
    }
}

pub fn human_duration(seconds: u32) -> String {
    match seconds {
        1 => "1 second".to_string(),
        s if s < 60 => format!("{} seconds", s),
        _ => format_clock(seconds),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_format_clock() {
        assert_eq!(format_clock(5), "5");
        assert_eq!(format_clock(30), "30");
        assert_eq!(format_clock(60), "01:00");
        assert_eq!(format_clock(90), "01:30");
        assert_eq!(format_clock(3600), "01:00:00");
        assert_eq!(format_clock(3661), "01:01:01");
        assert_eq!(format_clock(7200), "02:00:00");
    }

    #[test]
    fn test_human_duration() {
        assert_eq!(human_duration(1), "1 second");
        assert_eq!(human_duration(30), "30 seconds");
        assert_eq!(human_duration(60), "01:00");
        assert_eq!(human_duration(90), "01:30");
        assert_eq!(human_duration(3600), "01:00:00");
    }
}
