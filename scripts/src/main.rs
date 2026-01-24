use std::process::Command;
use std::io::{self, Write};
use std::thread;
use std::time::Instant;

// ansi
const GREEN: &str = "\x1b[32m";
const RED: &str = "\x1b[31m";
const CYAN: &str = "\x1b[36m";
const RESET: &str = "\x1b[0m";

fn main() {
    let start = Instant::now();
    println!("{}jobs:{}", CYAN, RESET);
    let fmt_handle = thread::spawn(|| exec("go", &["fmt", "./..."]));
    let tidy_handle = thread::spawn(|| exec("go", &["mod", "tidy"]));

    print!("{:<25} ... ", "[fmt]:");
    io::stdout().flush().unwrap();

    let success = fmt_handle.join().unwrap() && tidy_handle.join().unwrap();
    
    if success {
        println!("{}[READY]{}", GREEN, RESET);
    } else {
        println!("{}[FAIL]{}", RED, RESET);
    }

    // lints code
    println!("{:<25} ... ", "[golangci-lint]:");
    if exec("golangci-lint", &["run", "./..."]) {
        println!("\n{}[PASS]{} Good ({}s total)", GREEN, RESET, start.elapsed().as_secs());
        println!("{}Nice Cock{}", GREEN, RESET);
    } else {
        println!("\n{}[ERR]{} Linter issues found..", RED, RESET);
    }

    print!("\nPress any key to exit..");
    io::stdout().flush().unwrap();
    let _ = io::stdin().read_line(&mut String::new());
}

fn exec(cmd: &str, args: &[&str]) -> bool {
    Command::new(cmd)
        .args(args)
        .stdout(std::process::Stdio::null()) // hide
        .stderr(std::process::Stdio::null())
        .status()
        .map(|s| s.success())
        .unwrap_or(false)
}
