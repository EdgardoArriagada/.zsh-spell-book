use response::Response;
use std::io::{Read, Write};
use std::net::TcpStream;

#[derive(Debug)]
pub enum ClientError {
    ServerNotStarted,
    ReadError,
    WriteError,
    DecodeError,
}

#[derive(Clone)]
pub struct Client {
    addr: String,
}

impl Client {
    pub fn new(addr: &str) -> Self {
        Self {
            addr: addr.to_owned(),
        }
    }

    pub fn run(&self, request_line: &str) -> Result<Response, ClientError> {
        match TcpStream::connect(&self.addr) {
            Ok(mut stream) => {
                if let Err(_) = stream.write(request_line.as_bytes()) {
                    return Err(ClientError::WriteError);
                }

                let mut buffer = Vec::new();
                match stream.read_to_end(&mut buffer) {
                    Ok(_) => Response::try_from(&buffer[..]).map_err(|_| ClientError::DecodeError),

                    Err(_) => Err(ClientError::ReadError),
                }
            }
            Err(_) => Err(ClientError::ServerNotStarted),
        }
    }
}

pub mod actions;
pub mod response;
