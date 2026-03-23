use std::fmt::{Display, Formatter, Result as FmtResult};

#[derive(Copy, Clone, Debug)]
pub enum StatusCode {
    Continue = 100,
    Processing = 102,
    Ok = 200,
    Created = 201,
    NotModified = 304,
    BadRequest = 400,
    NotFound = 404,
    Conflict = 409,
    InternalServerError = 500,
}

impl Display for StatusCode {
    fn fmt(&self, f: &mut Formatter) -> FmtResult {
        write!(f, "{}", *self as u16)
    }
}
