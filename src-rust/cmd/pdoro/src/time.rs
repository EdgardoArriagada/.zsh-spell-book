use zsb_time::{parse_short_format, format_clock};

pub struct Time {
    seconds: u32,
}

impl Time {
    pub fn new(input: &str) -> Result<Self, String> {
        let duration = parse_short_format(input)?;
        let seconds = duration.seconds();

        if seconds == 0 || seconds >= 36000 {
            return Err("Invalid length: input must be in range (1s < input < 10h)".to_string());
        }

        Ok(Self { seconds })
    }

    pub fn get_seconds(&self) -> u32 {
        self.seconds
    }

    pub fn get_clock_from_seconds(seconds: &u32) -> String {
        let clock = format_clock(*seconds);
        // pdoro pads bare seconds with leading zero (e.g., "05" not "5")
        if *seconds < 60 {
            return format!("{:02}", seconds % 60);
        }
        clock
    }
}
