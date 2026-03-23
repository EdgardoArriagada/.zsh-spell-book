use std::convert::TryFrom;
use std::fmt::{Display, Formatter, Result as FmtResult};
use std::str;

#[derive(Debug)]
pub struct Request<'buf> {
    path: &'buf str,
    arg1: Option<&'buf str>,
    arg2: Option<&'buf str>,
}

impl<'buf> Display for Request<'buf> {
    fn fmt(&self, f: &mut Formatter<'_>) -> FmtResult {
        write!(f, "{}", self.display())
    }
}

impl<'buf> Request<'buf> {
    pub fn display(&self) -> String {
        match (self.arg1, self.arg2) {
            (Some(arg1), Some(arg2)) => {
                format!("path: {}, arg1: {}, arg2: {}", self.path, arg1, arg2)
            }
            (Some(arg1), None) => format!("path: {}, arg1: {}", self.path, arg1),
            (None, None) => format!("path: {}", self.path),
            _ => unreachable!(),
        }
    }
}

impl<'buf> Request<'buf> {
    pub fn path(&self) -> &str {
        &self.path
    }

    pub fn arg1(&self) -> Option<&str> {
        self.arg1
    }

    pub fn arg2(&self) -> Option<&str> {
        self.arg2
    }
}

impl<'buf> TryFrom<&'buf [u8]> for Request<'buf> {
    type Error = String;

    fn try_from(buf: &'buf [u8]) -> Result<Request<'buf>, Self::Error> {
        let request = str::from_utf8(buf).map_err(|_| "Invalid request")?;

        let request = match request.rfind(";") {
            Some(i) => &request[..i],
            None => Err("Request not terminated by ';' char".to_owned())?,
        };

        let (path, arg1, arg2) = parse_request(&request);

        match path {
            Some(path) => Ok(Self { path, arg1, arg2 }),
            None => Err("Invalid request line".to_owned()),
        }
    }
}

fn parse_request(request: &str) -> (Option<&str>, Option<&str>, Option<&str>) {
    let mut parts = request.splitn(3, ' ');

    (parts.next(), parts.next(), parts.next())
}
