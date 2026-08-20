use std::{env, process};

fn main() {
    let mut args = env::args().skip(1);
    let Some(file) = args.next() else {
        eprintln!("Usage: zsb_play <sound-file>");
        process::exit(1);
    };

    process::exit(zsb_play::play(&file));
}
