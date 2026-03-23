use std::fs::File;

use daemonize::Daemonize;

use crate::client::{response::Response, Client};
use crate::server::tcp_handler::TCPHandler;
use crate::server::Server;
use crate::time::Time;
use crate::utils::{stderr, stdout};

use super::ClientError;

trait HandledRun {
    fn safe_run(&self, path: &str, callback: fn(res: &Response) -> ());
}

impl HandledRun for Client {
    fn safe_run(&self, path: &str, callback: fn(res: &Response) -> ()) {
        match self.run(path) {
            Ok(res) => callback(&res),
            Err(ClientError::ServerNotStarted) => stderr("Pdoro server has not been started."),
            Err(e) => stderr(format!("Error: {:?}", e).as_str()),
        }
    }
}

static IP: &'static str = "127.0.0.1:51789";

pub fn remaining() {
    Client::new(IP).safe_run("remaining;", |res| {
        let digits = res.valid_msg().expect("Failed to retrieve remaining time.");

        let seconds = digits
            .parse::<u32>()
            .expect("Failed to parse remaining time.");

        let clock = Time::get_clock_from_seconds(&seconds);

        match (clock.as_str(), res.status()) {
            ("00", _) => stdout("No pomodoro timer is running."),
            (_, 304) => stdout(format!("{} (paused)", &clock).as_str()),
            _ => stdout(&clock),
        }
    });
}

pub fn is_counter_running() {
    Client::new(IP).safe_run("is-counter-running;", |res| match res.status() {
        100 | 102 => stdout(res.msg()),
        _ => stderr(res.msg()),
    });
}

pub fn is_valid_time(input: &str) {
    match Time::new(&input) {
        Ok(_) => stdout("true"),
        Err(_) => stdout("false"),
    }
}

pub fn start(time: &str, callback_with_args: &str) {
    let start_request = match get_start_request(time, callback_with_args) {
        Ok(req) => req,
        Err(e) => return stderr(e.as_str()),
    };

    Client::new(IP).safe_run(start_request.as_str(), |res| match res.status() {
        201 => stdout(res.msg()),
        _ => stderr(res.msg()),
    });
}

pub fn pause_resume_counter() {
    Client::new(IP).safe_run("pause-resume-counter;", |res| match res.status() {
        200 => stdout(res.msg()),
        _ => stderr(res.msg()),
    });
}

pub fn halt_counter() {
    Client::new(IP).safe_run("halt-counter;", |res| match res.status() {
        200 => stdout(res.msg()),
        _ => stderr(res.msg()),
    });
}

pub fn start_server() {
    match Client::new(IP).run("healthcheck;") {
        Ok(_) => return stderr("Pomodoro server already running."),
        Err(ClientError::ServerNotStarted) => {
            println!("starting...");
            start_daemon_server()
        }
        Err(e) => return stderr(format!("Error: {:?}", e).as_str()),
    }
}

fn get_start_request(time_arg: &str, callback_with_args: &str) -> Result<String, String> {
    let seconds = Time::new(time_arg)?.get_seconds();

    Ok(format!("start {} {};", seconds, callback_with_args))
}

fn start_daemon_server() {
    let stdout_file = File::create("/tmp/pdoro.out").expect("Failed to create stdout file.");
    let stderr_file = File::create("/tmp/pdoro.err").expect("Failed to create stderr file.");

    let daemonize = Daemonize::new()
        .working_directory("/tmp")
        .user("nobody")
        .group("pdoro_daemon")
        .umask(0o777)
        .stdout(stdout_file)
        .stderr(stderr_file)
        .privileged_action(|| "Executed before drop privileges");

    match daemonize.start() {
        Ok(_) => println!("Success, daemonized"),
        Err(e) => eprintln!("Error, {}", e),
    }

    Server::new(IP).run(TCPHandler);
}
