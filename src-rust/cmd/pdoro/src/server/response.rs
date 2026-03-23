use std::io::{Result as IoResult, Write};

use super::status_code::StatusCode;

#[derive(Debug)]
pub struct Response {
    status_code: StatusCode,
    body: Option<String>,
}

impl Response {
    pub fn new(status_code: StatusCode, body: Option<String>) -> Self {
        Self { status_code, body }
    }

    pub fn send(&self, stream: &mut impl Write) -> IoResult<()> {
        match &self.body {
            Some(body) => write!(stream, "{} {};", self.status_code, body),
            None => write!(stream, "{};", self.status_code),
        }
    }
}
