pub enum TimeUnit {
    Hours,
    Minutes,
    Seconds,
}

pub struct ShortDuration {
    pub unit: TimeUnit,
    pub value: u32,
}

impl ShortDuration {
    pub fn seconds(&self) -> u32 {
        match self.unit {
            TimeUnit::Hours => self.value * 3600,
            TimeUnit::Minutes => self.value * 60,
            TimeUnit::Seconds => self.value,
        }
    }
}

/// Parse short format like "5m", "1h", "30s" into a ShortDuration.
pub fn parse_short_format(input: &str) -> Result<ShortDuration, String> {
    let mut digits = String::new();
    let mut has_num = false;

    for c in input.chars() {
        match c {
            c if c.is_alphabetic() && !has_num => return Err("Invalid time format".to_string()),
            c if c.is_numeric() => {
                has_num = true;
                digits.push(c);
            }
            'h' | 'H' => return Ok(make_duration(TimeUnit::Hours, &digits)),
            'm' | 'M' => return Ok(make_duration(TimeUnit::Minutes, &digits)),
            's' | 'S' => return Ok(make_duration(TimeUnit::Seconds, &digits)),
            _ => return Err("Invalid time format".to_string()),
        }
    }

    Err("Invalid time format".to_string())
}

fn make_duration(unit: TimeUnit, digits: &str) -> ShortDuration {
    ShortDuration {
        unit,
        value: digits.parse::<u32>().expect("Failed to parse time"),
    }
}

/// Parse clock format "MM:SS" or "HH:MM:SS" into total seconds.
fn parse_clock_format(input: &str) -> Result<u32, String> {
    let parts: Vec<&str> = input.split(':').collect();

    let (hours, minutes, seconds) = match parts.len() {
        2 => {
            let m = parse_component(parts[0], 59)?;
            let s = parse_component(parts[1], 59)?;
            (0, m, s)
        }
        3 => {
            let h = parse_component(parts[0], 59)?;
            let m = parse_component(parts[1], 59)?;
            let s = parse_component(parts[2], 59)?;
            (h, m, s)
        }
        _ => return Err("Invalid time format".to_string()),
    };

    Ok(hours * 3600 + minutes * 60 + seconds)
}

fn parse_component(s: &str, max: u32) -> Result<u32, String> {
    let val = s
        .parse::<u32>()
        .map_err(|_| "Invalid time format".to_string())?;
    if val > max {
        return Err("Invalid time format".to_string());
    }
    Ok(val)
}

/// Parse any supported duration format into total seconds.
/// Tries short format (5m, 1h, 30s), then clock format (MM:SS, HH:MM:SS), then raw integer.
pub fn parse_duration(input: &str) -> Result<u32, String> {
    if let Ok(d) = parse_short_format(input) {
        return Ok(d.seconds());
    }

    if let Ok(secs) = parse_clock_format(input) {
        return Ok(secs);
    }

    input
        .parse::<u32>()
        .map_err(|_| "Invalid time format. Use formats like: 30 (seconds), 5m, 1h, 05:30, or 1:05:30".to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_parse_short_format() {
        assert_eq!(parse_short_format("5s").unwrap().seconds(), 5);
        assert_eq!(parse_short_format("5S").unwrap().seconds(), 5);
        assert_eq!(parse_short_format("10m").unwrap().seconds(), 600);
        assert_eq!(parse_short_format("10M").unwrap().seconds(), 600);
        assert_eq!(parse_short_format("2h").unwrap().seconds(), 7200);
        assert_eq!(parse_short_format("2H").unwrap().seconds(), 7200);
        assert!(parse_short_format("invalid").is_err());
        assert!(parse_short_format("5").is_err());
        assert!(parse_short_format("5x").is_err());
        assert!(parse_short_format("").is_err());
    }

    #[test]
    fn test_parse_clock_format() {
        assert_eq!(parse_clock_format("05:30").unwrap(), 330);
        assert_eq!(parse_clock_format("1:05:30").unwrap(), 3930);
        assert_eq!(parse_clock_format("00:01").unwrap(), 1);
        assert_eq!(parse_clock_format("59:59").unwrap(), 3599);
        assert_eq!(parse_clock_format("1:00:00").unwrap(), 3600);
        assert!(parse_clock_format("invalid").is_err());
        assert!(parse_clock_format("60:00").is_err());
        assert!(parse_clock_format("00:60").is_err());
        assert!(parse_clock_format("").is_err());
    }

    #[test]
    fn test_parse_duration_short() {
        assert_eq!(parse_duration("5s").unwrap(), 5);
        assert_eq!(parse_duration("10m").unwrap(), 600);
        assert_eq!(parse_duration("2h").unwrap(), 7200);
    }

    #[test]
    fn test_parse_duration_clock() {
        assert_eq!(parse_duration("05:30").unwrap(), 330);
        assert_eq!(parse_duration("1:05:30").unwrap(), 3930);
    }

    #[test]
    fn test_parse_duration_raw() {
        assert_eq!(parse_duration("30").unwrap(), 30);
        assert_eq!(parse_duration("3600").unwrap(), 3600);
    }

    #[test]
    fn test_parse_duration_invalid() {
        assert!(parse_duration("abc").is_err());
        assert!(parse_duration("").is_err());
    }
}
