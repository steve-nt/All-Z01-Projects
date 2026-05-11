use std::collections::VecDeque;

use crate::board::{Board, Cell};
use crate::piece::Piece;

pub fn legal_placements(board: &Board, piece: &Piece) -> Vec<(u16, u16)> {
    let mut out = Vec::new();
    if piece.w > board.w || piece.h > board.h {
        return out;
    }
    let max_x = board.w - piece.w;
    let max_y = board.h - piece.h;
    for y in 0..=max_y {
        'next: for x in 0..=max_x {
            let mut overlap_mine = 0u32;
            for &(dx, dy) in &piece.hashes {
                match board.get(x + dx, y + dy) {
                    Cell::Theirs => continue 'next,
                    Cell::Mine => {
                        overlap_mine += 1;
                        if overlap_mine > 1 {
                            continue 'next;
                        }
                    }
                    Cell::Empty => {}
                }
            }
            if overlap_mine == 1 {
                out.push((x, y));
            }
        }
    }
    out
}

pub fn chebyshev_field(board: &Board) -> Vec<u32> {
    let w = board.w as usize;
    let mut dist = vec![u32::MAX; w * (board.h as usize)];
    let mut q: VecDeque<(u16, u16)> = VecDeque::new();
    let mut any_source = false;
    for y in 0..board.h {
        for x in 0..board.w {
            if board.get(x, y) == Cell::Theirs {
                dist[(y as usize) * w + (x as usize)] = 0;
                q.push_back((x, y));
                any_source = true;
            }
        }
    }
    if !any_source {
        return vec![0; w * (board.h as usize)];
    }
    while let Some((x, y)) = q.pop_front() {
        let d = dist[(y as usize) * w + (x as usize)];
        for dy in -1i32..=1 {
            for dx in -1i32..=1 {
                if dx == 0 && dy == 0 {
                    continue;
                }
                let nx = x as i32 + dx;
                let ny = y as i32 + dy;
                if nx < 0 || ny < 0 || nx >= board.w as i32 || ny >= board.h as i32 {
                    continue;
                }
                let idx = (ny as usize) * w + (nx as usize);
                if dist[idx] == u32::MAX {
                    dist[idx] = d + 1;
                    q.push_back((nx as u16, ny as u16));
                }
            }
        }
    }
    dist
}

fn score(
    field: &[u32],
    board_w: u16,
    piece: &Piece,
    x: u16,
    y: u16,
) -> (u64, u32, u16, u16) {
    let mut sum: u64 = 0;
    let mut max: u32 = 0;
    for &(dx, dy) in &piece.hashes {
        let bx = x + dx;
        let by = y + dy;
        let d = field[(by as usize) * (board_w as usize) + (bx as usize)];
        sum = sum.saturating_add(d as u64);
        if d > max {
            max = d;
        }
    }
    (sum, max, y, x)
}

pub fn choose(board: &Board, piece: &Piece) -> (u16, u16) {
    let legal = legal_placements(board, piece);
    if legal.is_empty() {
        return (0, 0);
    }
    let field = chebyshev_field(board);
    legal
        .into_iter()
        .map(|(x, y)| (score(&field, board.w, piece, x, y), (x, y)))
        .min_by_key(|(s, _)| *s)
        .map(|(_, xy)| xy)
        .unwrap_or((0, 0))
}

#[cfg(test)]
mod tests {
    use super::*;

    fn board_from_rows(rows: &[&str]) -> Board {
        let h = rows.len() as u16;
        let w = rows[0].chars().count() as u16;
        let mut b = Board::new(w, h);
        for (y, row) in rows.iter().enumerate() {
            for (x, ch) in row.chars().enumerate() {
                let c = match ch {
                    '@' => Cell::Mine,
                    '$' => Cell::Theirs,
                    _ => Cell::Empty,
                };
                b.set(x as u16, y as u16, c);
            }
        }
        b
    }

    #[test]
    fn accepts_exactly_one_overlap() {
        let b = board_from_rows(&["@...", "....", "....", "...$"]);
        let p = Piece::new(2, 1, vec![(0, 0), (1, 0)]);
        let legal = legal_placements(&b, &p);
        assert!(legal.contains(&(0, 0)));
    }

    #[test]
    fn rejects_zero_overlap() {
        let b = board_from_rows(&["@...", "....", "....", "...$"]);
        let p = Piece::new(2, 1, vec![(0, 0), (1, 0)]);
        let legal = legal_placements(&b, &p);
        assert!(!legal.contains(&(2, 2)));
    }

    #[test]
    fn rejects_two_mine_overlaps() {
        let b = board_from_rows(&["@@..", "....", "....", "...$"]);
        let p = Piece::new(2, 1, vec![(0, 0), (1, 0)]);
        let legal = legal_placements(&b, &p);
        assert!(!legal.contains(&(0, 0)));
    }

    #[test]
    fn rejects_opponent_overlap() {
        let b = board_from_rows(&["@$..", "....", "....", "...$"]);
        let p = Piece::new(2, 1, vec![(0, 0), (1, 0)]);
        let legal = legal_placements(&b, &p);
        assert!(!legal.contains(&(0, 0)));
    }

    #[test]
    fn oob_handled_by_top_left_range() {
        let b = board_from_rows(&["@...", "....", "....", "...$"]);
        let p = Piece::new(2, 1, vec![(0, 0), (1, 0)]);
        let legal = legal_placements(&b, &p);
        for &(x, _) in &legal {
            assert!(x + 2 <= b.w);
        }
    }

    #[test]
    fn piece_larger_than_board_yields_none() {
        let b = board_from_rows(&["@.", ".."]);
        let p = Piece::new(3, 3, vec![(0, 0)]);
        assert!(legal_placements(&b, &p).is_empty());
    }

    #[test]
    fn chebyshev_field_single_source() {
        let b = board_from_rows(&["...", "..$", "..."]);
        let f = chebyshev_field(&b);
        assert_eq!(f[(0) * 3 + 0], 2);
        assert_eq!(f[(0) * 3 + 2], 1);
        assert_eq!(f[(1) * 3 + 2], 0);
        assert_eq!(f[(2) * 3 + 0], 2);
    }

    #[test]
    fn chebyshev_field_no_sources() {
        let b = board_from_rows(&["...", "..."]);
        let f = chebyshev_field(&b);
        assert!(f.iter().all(|&d| d == 0));
    }

    #[test]
    fn heuristic_picks_closest_to_opponent() {
        // 1x1 piece, mine at (0,0) and (4,0), opponent at (6,0).
        // The placement closer to (6,0) should be chosen.
        // But 1x1 piece with one # requires exactly-one overlap with mine,
        // so legal placements are only at (0,0) and (4,0).
        let b = board_from_rows(&["@...@.$"]);
        let p = Piece::new(1, 1, vec![(0, 0)]);
        let (x, y) = choose(&b, &p);
        assert_eq!((x, y), (4, 0));
    }

    #[test]
    fn heuristic_fallback_when_no_legal_move() {
        let b = board_from_rows(&["..", ".."]);
        let p = Piece::new(1, 1, vec![(0, 0)]);
        assert_eq!(choose(&b, &p), (0, 0));
    }

    #[test]
    fn tiebreak_lowest_y_then_x() {
        // Two mines, no opponent → field is all zeros → tiebreak picks
        // lowest (y, x).
        let b = board_from_rows(&["@..", "...", "..@"]);
        let p = Piece::new(1, 1, vec![(0, 0)]);
        let (x, y) = choose(&b, &p);
        assert_eq!((x, y), (0, 0));
    }
}
