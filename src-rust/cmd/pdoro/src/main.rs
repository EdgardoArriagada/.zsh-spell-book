mod args;
mod client;
mod server;
mod time;
mod utils;

use args::Args;
use clap::Parser;

use client::actions;
use utils::stderr;

fn main() {
    let args = Args::parse();

    if args.remaining {
        return actions::remaining();
    }

    if args.is_counter_running {
        return actions::is_counter_running();
    }

    if let Some(input) = args.is_valid_time {
        return actions::is_valid_time(&input);
    }

    match (args.time, args.callback_with_args) {
        (Some(time), Some(callback_with_args)) => {
            return actions::start(&time, &callback_with_args)
        }
        (Some(_), None) | (None, Some(_)) => {
            return stderr("Both time and callback_with_args must be provided.")
        }
        _ => {}
    }

    if args.pause_resume_counter {
        return actions::pause_resume_counter();
    }

    if args.halt_counter {
        return actions::halt_counter();
    }

    if args.start_server {
        return actions::start_server();
    }

    stderr("No arguments provided.");
}
