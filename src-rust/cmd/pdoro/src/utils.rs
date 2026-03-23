use std::process;

use std::thread;
use std::time::Duration;

pub fn sleep(secs: u64) {
    thread::sleep(Duration::from_secs(secs))
}

pub fn stderr(msg: &str) {
    eprintln!("{}", msg);
    process::exit(1);
}

pub fn stdout(msg: &str) {
    println!("{}", msg);
    process::exit(0);
}
