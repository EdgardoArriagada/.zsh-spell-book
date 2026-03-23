use std::str;

pub struct Response {
    status: u16,
    msg: String,
}

impl Response {
    pub fn status(&self) -> u16 {
        self.status
    }

    pub fn msg(&self) -> &str {
        &self.msg
    }

    pub fn valid_msg(&self) -> Result<&str, &str> {
        match self.msg().is_empty() {
            false => Ok(&self.msg),
            true => Err("No message"),
        }
    }
}

impl TryFrom<&[u8]> for Response {
    type Error = String;

    fn try_from(buf: &[u8]) -> Result<Response, Self::Error> {
        let response = str::from_utf8(buf).map_err(|_| "Invalid response")?;

        let mut parts = response.splitn(2, ' ');

        let status = parts.next().unwrap_or("422");
        let msg = parts.next().unwrap_or("");

        let status = status.parse::<u16>().expect("Could not parse status code");
        let msg = match msg.rfind(';') {
            Some(i) => &msg[..i],
            None => "",
        };

        Ok(Response {
            status,
            msg: msg.to_owned(),
        })
    }
}
