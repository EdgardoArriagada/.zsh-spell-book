pub enum TimeFormat {
    Hours,
    Minutes,
    Seconds,
}

pub struct Time {
    pub format: TimeFormat,
    pub value: u32,
}

impl Time {
    pub fn new(input: &str) -> Result<Self, String> {
        let time = Self::input_to_time(input)?;

        match Self::is_valid_len(&time) {
            true => Ok(time),
            false => Err("Invalid length: input must be in range (1s < input < 10h)".to_string()),
        }
    }

    fn input_to_time(input: &str) -> Result<Self, String> {
        let mut value = String::new();
        let mut has_num = false;

        for c in input.chars() {
            match c {
                c if c.is_alphabetic() && !has_num => return Err("Invalid time format".to_string()),
                c if c.is_numeric() => {
                    has_num = true;
                    value.push(c)
                }
                'h' | 'H' => return Ok(Self::parsed(TimeFormat::Hours, &value)),
                'm' | 'M' => return Ok(Self::parsed(TimeFormat::Minutes, &value)),
                's' | 'S' => return Ok(Self::parsed(TimeFormat::Seconds, &value)),
                _ => return Err("Invalid time format".to_string()),
            }
        }

        Err("Invalid time format".to_string())
    }

    fn is_valid_len(time: &Self) -> bool {
        if time.value == 0 {
            return false;
        }

        match time.format {
            TimeFormat::Hours => time.value < 10,
            TimeFormat::Minutes => time.value < 60 * 10,
            TimeFormat::Seconds => time.value < 60 * 60 * 10,
        }
    }

    pub fn parsed(format: TimeFormat, value: &str) -> Self {
        Self {
            format,
            value: value.parse::<u32>().expect("Failed to parse time"),
        }
    }

    pub fn get_seconds(&self) -> u32 {
        match self.format {
            TimeFormat::Hours => self.value * 60 * 60,
            TimeFormat::Minutes => self.value * 60,
            TimeFormat::Seconds => self.value,
        }
    }

    pub fn get_clock_from_seconds(seconds: &u32) -> String {
        let hours = seconds / 60 / 60;
        let minutes = seconds / 60 % 60;
        let seconds = seconds % 60;

        match (hours, minutes, seconds) {
            (0, 0, _) => return format!("{:02}", seconds),
            (0, _, _) => return format!("{:02}:{:02}", minutes, seconds),
            _ => return format!("{:02}:{:02}:{:02}", hours, minutes, seconds),
        }
    }
}
