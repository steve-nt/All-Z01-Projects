mod board;
mod io;
mod piece;
mod strategy;

use std::io::{BufRead, BufReader, Write};

fn main() {
    let stdin = std::io::stdin();
    let mut reader = BufReader::new(stdin.lock());
    let stdout = std::io::stdout();
    let mut writer = stdout.lock();

    let player = loop {
        let mut line = String::new();
        match reader.read_line(&mut line) {
            Ok(0) => return,
            Ok(_) => {
                if let Some(p) = io::parse_greeting(line.trim_end()) {
                    break p;
                }
            }
            Err(_) => return,
        }
    };

    loop {
        let board = match io::parse_board(&mut reader, player) {
            Ok(Some(b)) => b,
            _ => return,
        };
        let piece = match io::parse_piece(&mut reader) {
            Ok(Some(p)) => p,
            _ => return,
        };
        let (x, y) = strategy::choose(&board, &piece);
        if writer.write_all(io::format_move(x, y).as_bytes()).is_err() {
            return;
        }
        if writer.flush().is_err() {
            return;
        }
    }
}
