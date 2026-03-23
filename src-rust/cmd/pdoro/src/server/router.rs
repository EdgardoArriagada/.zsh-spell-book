use super::controllers::{
    halt_counter, health_check, is_counter_running, not_found, pause_resume_counter,
    remaining_pomodoro, start_pomodoro,
};
use super::request::Request;
use super::response::Response;

pub fn router(request: &Request) -> Response {
    match request.path() {
        "healthcheck" => health_check(),
        "start" => start_pomodoro(request),
        "halt-counter" => halt_counter(),
        "remaining" => remaining_pomodoro(),
        "is-counter-running" => is_counter_running(),
        "pause-resume-counter" => pause_resume_counter(),
        _ => not_found(),
    }
}
