use std::process;

fn main() {
    let args: Vec<String> = std::env::args().skip(1).collect();
    if args.len() != 2 {
        eprintln!("Usage: repeatstr <number> <word>");
        process::exit(1);
    }

    let n: usize = args[0].parse().unwrap_or_else(|_| {
        eprintln!("First argument must be a number");
        process::exit(1);
    });

    println!("{}", args[1].repeat(n));
}
