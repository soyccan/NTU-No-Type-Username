use std::env;
use std::error;
use std::fs::File;
use std::io::Read;
use std::process::exit;

mod tui;

fn main() -> Result<(), Box<dyn error::Error>> {
    // let cookie_path = get_env("COOKIE_PATH");
    // let cookies = read_file_to_string(&cookie_path);

    Ok(tui::tui_start()?)
}

fn read_file_to_string(path: &str) -> String {
    let mut file = File::open(path).unwrap_or_else(|_| {
        eprintln!("open {path} failed");
        exit(1)
    });

    let mut buffer = String::new();
    file.read_to_string(&mut buffer).unwrap_or_else(|_| {
        eprintln!("read {path} failed");
        exit(1)
    });

    buffer
}

fn get_env(key: &str) -> String {
    env::var(key).unwrap_or_else(|_| {
        eprintln!("{key} not set");
        exit(1)
    })
}
