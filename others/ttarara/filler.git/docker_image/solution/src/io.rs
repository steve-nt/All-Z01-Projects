use std::io::BufRead;

use crate::board::{Board, Cell, PlayerId};
use crate::piece::Piece;

pub fn parse_greeting(line: &str) -> Option<PlayerId> {
    let rest = line.strip_prefix("$$$ exec p")?;
    match rest.chars().next()? {
        '1' => Some(PlayerId::P1),
        '2' => Some(PlayerId::P2),
        _ => None,
    }
}

pub fn format_move(x: u16, y: u16) -> String {
    format!("{} {}\n", x, y)
}

fn read_line<R: BufRead>(r: &mut R) -> std::io::Result<Option<String>> {
    let mut buf = String::new();
    let n = r.read_line(&mut buf)?;
    if n == 0 {
        Ok(None)
    } else {
        if buf.ends_with('\n') {
            buf.pop();
            if buf.ends_with('\r') {
                buf.pop();
            }
        }
        Ok(Some(buf))
    }
}

fn parse_wh(header: &str, tag: &str) -> Option<(u16, u16)> {
    let rest = header.strip_prefix(tag)?.trim_end_matches(':').trim();
    let mut parts = rest.split_whitespace();
    let w: u16 = parts.next()?.parse().ok()?;
    let h: u16 = parts.next()?.parse().ok()?;
    Some((w, h))
}

pub fn parse_board<R: BufRead>(r: &mut R, player: PlayerId) -> std::io::Result<Option<Board>> {
    let header = match read_line(r)? {
        Some(s) => s,
        None => return Ok(None),
    };
    let (w, h) = match parse_wh(&header, "Anfield ") {
        Some(v) => v,
        None => return Ok(None),
    };
    let _ = read_line(r)?;
    let mut board = Board::new(w, h);
    for y in 0..h {
        let row = match read_line(r)? {
            Some(s) => s,
            None => return Ok(None),
        };
        let cells: &str = if row.len() >= 4 { &row[4..] } else { "" };
        for (x, ch) in cells.chars().enumerate() {
            if x >= w as usize {
                break;
            }
            let c = player.classify(ch);
            if !matches!(c, Cell::Empty) {
                board.set(x as u16, y, c);
            }
        }
    }
    Ok(Some(board))
}

pub fn parse_piece<R: BufRead>(r: &mut R) -> std::io::Result<Option<Piece>> {
    let header = match read_line(r)? {
        Some(s) => s,
        None => return Ok(None),
    };
    let (w, h) = match parse_wh(&header, "Piece ") {
        Some(v) => v,
        None => return Ok(None),
    };
    let mut hashes = Vec::new();
    for y in 0..h {
        let row = match read_line(r)? {
            Some(s) => s,
            None => return Ok(None),
        };
        for (x, ch) in row.chars().enumerate() {
            if x >= w as usize {
                break;
            }
            if ch != '.' && !ch.is_whitespace() {
                hashes.push((x as u16, y));
            }
        }
    }
    Ok(Some(Piece::new(w, h, hashes)))
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::io::Cursor;

    #[test]
    fn greeting_p1() {
        assert_eq!(
            parse_greeting("$$$ exec p1 : [./robots/bender]"),
            Some(PlayerId::P1)
        );
    }

    #[test]
    fn greeting_p2() {
        assert_eq!(
            parse_greeting("$$$ exec p2 : [./robots/terminator]"),
            Some(PlayerId::P2)
        );
    }

    #[test]
    fn greeting_rejects_garbage() {
        assert_eq!(parse_greeting("hello world"), None);
        assert_eq!(parse_greeting("$$$ exec p3 : [x]"), None);
    }

    #[test]
    fn coordinate_output_is_x_y_newline() {
        assert_eq!(format_move(7, 2), "7 2\n");
        assert_eq!(format_move(0, 0), "0 0\n");
    }

    #[test]
    fn board_parse_basic() {
        let input = "Anfield 5 3:\n    01234\n000 .....\n001 ..@..\n002 ..a..\n";
        let mut r = Cursor::new(input);
        let b = parse_board(&mut r, PlayerId::P1).unwrap().unwrap();
        assert_eq!(b.w, 5);
        assert_eq!(b.h, 3);
        assert_eq!(b.get(2, 1), Cell::Mine);
        assert_eq!(b.get(2, 2), Cell::Mine);
        assert_eq!(b.get(0, 0), Cell::Empty);
    }

    #[test]
    fn board_parse_classifies_by_player() {
        let input = "Anfield 3 1:\n    012\n000 @.$\n";
        let mut r = Cursor::new(input);
        let b = parse_board(&mut r, PlayerId::P2).unwrap().unwrap();
        assert_eq!(b.get(0, 0), Cell::Theirs);
        assert_eq!(b.get(2, 0), Cell::Mine);
    }

    #[test]
    fn piece_parse_basic() {
        let input = "Piece 4 2:\n..##\n.##.\n";
        let mut r = Cursor::new(input);
        let p = parse_piece(&mut r).unwrap().unwrap();
        assert_eq!(p.w, 4);
        assert_eq!(p.h, 2);
        assert_eq!(p.hashes, vec![(2, 0), (3, 0), (1, 1), (2, 1)]);
    }

    #[test]
    fn piece_parse_with_o_character() {
        // The bundled 01-edu engine prints piece cells as 'O', not '#'.
        let input = "Piece 4 2:\n..OO\n.OO.\n";
        let mut r = Cursor::new(input);
        let p = parse_piece(&mut r).unwrap().unwrap();
        assert_eq!(p.hashes, vec![(2, 0), (3, 0), (1, 1), (2, 1)]);
    }

    #[test]
    fn piece_parse_all_empty_row() {
        let input = "Piece 3 2:\n...\n.#.\n";
        let mut r = Cursor::new(input);
        let p = parse_piece(&mut r).unwrap().unwrap();
        assert_eq!(p.hashes, vec![(1, 1)]);
    }

    #[test]
    fn piece_parse_trailing_crlf_and_whitespace() {
        let input = "Piece 3 1:\r\n.#.  \r\n";
        let mut r = Cursor::new(input);
        let p = parse_piece(&mut r).unwrap().unwrap();
        assert_eq!(p.w, 3);
        assert_eq!(p.h, 1);
        assert_eq!(p.hashes, vec![(1, 0)]);
    }

    #[test]
    fn board_parse_larger_row_prefix_widths_are_ok() {
        // The engine may print row indices wider than 3 digits for big
        // maps; we rely on the cells starting at column 4 regardless.
        let input = "Anfield 3 1:\n    012\n000 @.$\n";
        let mut r = Cursor::new(input);
        let b = parse_board(&mut r, PlayerId::P1).unwrap().unwrap();
        assert_eq!(b.get(0, 0), Cell::Mine);
        assert_eq!(b.get(2, 0), Cell::Theirs);
    }

    #[test]
    fn eof_returns_none() {
        let mut r = Cursor::new("");
        assert!(parse_board(&mut r, PlayerId::P1).unwrap().is_none());
        let mut r = Cursor::new("");
        assert!(parse_piece(&mut r).unwrap().is_none());
    }
}
