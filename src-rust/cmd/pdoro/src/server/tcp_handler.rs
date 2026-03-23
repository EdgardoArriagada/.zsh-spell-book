use super::request::Request;
use super::response::Response;
use super::router::router;
use super::Handler;

pub struct TCPHandler;

impl Handler for TCPHandler {
    fn handle_request(&self, request: &Request) -> Response {
        router(request)
    }
}
