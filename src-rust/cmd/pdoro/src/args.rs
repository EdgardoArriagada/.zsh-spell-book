use clap::Parser;

/// pdoro
#[derive(Debug, Parser)]
#[command(author, about, long_about = None)]
pub struct Args {
    /// time duration of session
    #[clap(long, short)]
    pub time: Option<String>,

    /// callback program with args
    #[clap(long, short)]
    pub callback_with_args: Option<String>,

    /// remaining duration of session
    #[clap(long, short)]
    pub remaining: bool,

    /// start pdoro server
    #[clap(long, short)]
    pub start_server: bool,

    /// halt pomodoro counter
    #[clap(long)]
    pub halt_counter: bool,

    /// toggle pause/resume pomodoro counter
    #[clap(long, short)]
    pub pause_resume_counter: bool,

    /// validate time duration
    #[clap(long)]
    pub is_valid_time: Option<String>,

    /// check if counter is running
    #[clap(long, short)]
    pub is_counter_running: bool,
}
