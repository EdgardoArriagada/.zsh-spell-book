use request::Request;
use response::Response;
use status_code::StatusCode;
use std::io::Read;
use std::net::TcpListener;

pub trait Handler {
    fn handle_request(&self, request: &Request) -> Response;

    fn handle_bad_request(&self, e: &str) -> Response {
        Response::new(StatusCode::BadRequest, Some(e.to_owned()))
    }
}

pub struct Server {
    addr: String,
}

impl Server {
    pub fn new(addr: &str) -> Self {
        Self {
            addr: addr.to_owned(),
        }
    }

    pub fn run(self, handler: impl Handler) {
        println!("Listening on {}", self.addr);

        let listener = TcpListener::bind(&self.addr).expect("Failed to bind address");

        loop {
            match listener.accept() {
                Ok((mut stream, _)) => {
                    let mut buffer = [0 as u8; 256];
                    match stream.read(&mut buffer) {
                        Ok(_) => {
                            let response = match Request::try_from(&buffer[..]) {
                                Ok(request) => handler.handle_request(&request),
                                Err(e) => handler.handle_bad_request(&e),
                            };

                            if let Err(e) = response.send(&mut stream) {
                                println!("Failed to send reponse: {}", e)
                            }
                        }
                        Err(e) => println!("Failed to read from connection: {}", e),
                    }
                }
                Err(e) => println!("Failed to establish a connection: {}", e),
            }
        }
    }
}

pub mod controllers;
pub mod request;
pub mod response;
pub mod router;
pub mod status_code;
pub mod tcp_handler;
