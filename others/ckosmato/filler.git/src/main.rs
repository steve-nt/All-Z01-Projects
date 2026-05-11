use std::io::{self, BufRead, Write};

// ─── Data Structures ────────────────────────────────────────────────────────

#[derive(Debug, Clone, PartialEq)]
pub struct Piece {
    pub rows: usize,
    pub cols: usize,
    pub cells: Vec<Vec<bool>>, // true = filled cell
}

#[derive(Debug, Clone)]
pub struct Anfield {
    pub rows: usize,
    pub cols: usize,
    pub grid: Vec<Vec<char>>,
}

#[derive(Debug, Clone, Copy, PartialEq)]
pub enum Player {
    P1, // represented by 'a'/'@'
    P2, // represented by 's'/'$'
}

impl Player {
    pub fn my_chars(&self) -> (char, char) {
        match self {
            Player::P1 => ('a', '@'),
            Player::P2 => ('s', '$'),
        }
    }
    pub fn opp_chars(&self) -> (char, char) {
        match self {
            Player::P1 => ('s', '$'),
            Player::P2 => ('a', '@'),
        }
    }
}

// ─── Parsing ────────────────────────────────────────────────────────────────

/// Parse the player number from the first line sent by game_engine.
/// Format: `$$$ exec p<number> : [<path>]`
pub fn parse_player(line: &str) -> Option<Player> {
    if line.contains("p1") {
        Some(Player::P1)
    } else if line.contains("p2") {
        Some(Player::P2)
    } else {
        None
    }
}

/// Parse anfield header and grid lines.
/// Header: `Anfield <cols> <rows>:`
/// Then rows with format: `NNN <grid_row>`
pub fn parse_anfield(lines: &[String]) -> Option<Anfield> {
    // First line: "Anfield <cols> <rows>:"
    let header = lines.first()?;
    let parts: Vec<&str> = header.split_whitespace().collect();
    if parts.len() < 3 || parts[0] != "Anfield" {
        return None;
    }
    let cols: usize = parts[1].parse().ok()?;
    let rows: usize = parts[2].trim_end_matches(':').parse().ok()?;

    let mut grid: Vec<Vec<char>> = Vec::with_capacity(rows);

    // Skip the column-index header line.
    // Grid lines have format "NNN <cells>" where NNN is exactly 3 digits followed by a space.
    let mut grid_lines = lines.iter().skip(1).filter(|l| {
        let b = l.as_bytes();
        b.len() > 3
            && b[0].is_ascii_digit()
            && b[1].is_ascii_digit()
            && b[2].is_ascii_digit()
            && b[3] == b' '
    });

    for _ in 0..rows {
        let line = grid_lines.next()?;
        // Format: "NNN <cells>"  — find first space then take the rest
        let cell_part = line
            .splitn(2, ' ')
            .nth(1)
            .unwrap_or("")
            .trim_end();
        let row: Vec<char> = cell_part.chars().take(cols).collect();
        grid.push(row);
    }

    if grid.len() != rows {
        return None;
    }

    Some(Anfield { rows, cols, grid })
}

/// Parse a piece from lines starting after "Piece <cols> <rows>:".
pub fn parse_piece(lines: &[String]) -> Option<Piece> {
    let header = lines.first()?;
    let parts: Vec<&str> = header.split_whitespace().collect();
    if parts.len() < 3 || parts[0] != "Piece" {
        return None;
    }
    let cols: usize = parts[1].parse().ok()?;
    let rows: usize = parts[2].trim_end_matches(':').parse().ok()?;

    let mut cells: Vec<Vec<bool>> = Vec::with_capacity(rows);
    for line in lines.iter().skip(1).take(rows) {
        let row: Vec<bool> = line.chars().take(cols).map(|c| c == '#' || c == 'O').collect();
        // Pad if line is shorter than cols
        let mut row = row;
        while row.len() < cols {
            row.push(false);
        }
        cells.push(row);
    }

    if cells.len() != rows {
        return None;
    }

    Some(Piece { rows, cols, cells })
}

// ─── Placement Logic ─────────────────────────────────────────────────────────

/// Check if placing `piece` at (col, row) on `anfield` is valid for `player`.
/// Rules:
///   - Piece must not go out of bounds.
///   - Exactly ONE cell of the piece overlaps the player's own territory.
///   - ZERO cells of the piece overlap the opponent's territory.
pub fn is_valid_placement(
    anfield: &Anfield,
    piece: &Piece,
    col: isize,
    row: isize,
    player: Player,
) -> bool {
    let (my_lower, my_upper) = player.my_chars();
    let (opp_lower, opp_upper) = player.opp_chars();

    let mut overlap_mine = 0usize;

    for pr in 0..piece.rows {
        for pc in 0..piece.cols {
            if !piece.cells[pr][pc] {
                continue;
            }
            let gr = row + pr as isize;
            let gc = col + pc as isize;

            // Boundary check
            if gr < 0 || gc < 0 || gr >= anfield.rows as isize || gc >= anfield.cols as isize {
                return false;
            }

            let cell = anfield.grid[gr as usize][gc as usize];

            // Opponent overlap → invalid
            if cell == opp_lower || cell == opp_upper {
                return false;
            }

            // Count own overlaps
            if cell == my_lower || cell == my_upper {
                overlap_mine += 1;
                if overlap_mine > 1 {
                    return false; // more than one own overlap
                }
            }
        }
    }

    overlap_mine == 1
}

// ─── Strategy ────────────────────────────────────────────────────────────────

/// Compute Manhattan distance from cell (r,c) to the nearest opponent cell.
fn dist_to_opponent(anfield: &Anfield, r: usize, c: usize, player: Player) -> usize {
    let (ol, ou) = player.opp_chars();
    let mut min_dist = usize::MAX;
    for gr in 0..anfield.rows {
        for gc in 0..anfield.cols {
            let cell = anfield.grid[gr][gc];
            if cell == ol || cell == ou {
                let d = r.abs_diff(gr) + c.abs_diff(gc);
                if d < min_dist {
                    min_dist = d;
                }
            }
        }
    }
    min_dist
}

/// Score a placement at (col, row).
/// Strategy: maximize territory gain (piece cells on empty squares)
/// while aggressively moving toward the opponent (minimize distance to opp).
fn score_placement(
    anfield: &Anfield,
    piece: &Piece,
    col: isize,
    row: isize,
    player: Player,
) -> i64 {
    let (opp_lower, opp_upper) = player.opp_chars();

    let mut new_cells = 0i64;
    let mut min_dist_to_opp = i64::MAX;

    for pr in 0..piece.rows {
        for pc in 0..piece.cols {
            if !piece.cells[pr][pc] {
                continue;
            }
            let gr = (row + pr as isize) as usize;
            let gc = (col + pc as isize) as usize;
            let cell = anfield.grid[gr][gc];

            if cell == '.' {
                new_cells += 1;

                // Distance from this new cell to nearest opponent
                let d = dist_to_opponent(anfield, gr, gc, player) as i64;
                if d < min_dist_to_opp {
                    min_dist_to_opp = d;
                }
            }

            // Extra bonus if we are adjacent to opponent — we're cutting them off
            for (dr, dc) in [(-1i64, 0), (1, 0), (0, -1i64), (0, 1)] {
                let nr = gr as i64 + dr;
                let nc = gc as i64 + dc;
                if nr >= 0
                    && nc >= 0
                    && nr < anfield.rows as i64
                    && nc < anfield.cols as i64
                {
                    let ncell = anfield.grid[nr as usize][nc as usize];
                    if ncell == opp_lower || ncell == opp_upper {
                        new_cells += 2; // bonus for being adjacent to opponent
                    }
                }
            }
        }
    }

    if min_dist_to_opp == i64::MAX {
        min_dist_to_opp = 0;
    }

    // Higher new_cells = better; lower dist_to_opp = better (aggressive)
    new_cells * 100 - min_dist_to_opp
}

/// Find the best (col, row) to place the piece, or None if no valid placement.
pub fn best_placement(
    anfield: &Anfield,
    piece: &Piece,
    player: Player,
) -> Option<(isize, isize)> {
    let mut best_score = i64::MIN;
    let mut best_pos: Option<(isize, isize)> = None;

    let row_range = -(piece.rows as isize)..anfield.rows as isize;
    let col_range = -(piece.cols as isize)..anfield.cols as isize;

    for row in row_range {
        for col in col_range.clone() {
            if is_valid_placement(anfield, piece, col, row, player) {
                let s = score_placement(anfield, piece, col, row, player);
                if s > best_score {
                    best_score = s;
                    best_pos = Some((col, row));
                }
            }
        }
    }

    best_pos
}

// ─── Main Game Loop ──────────────────────────────────────────────────────────

fn main() {
    let stdin = io::stdin();
    let stdout = io::stdout();
    let mut out = io::BufWriter::new(stdout.lock());

    let mut lines_iter = stdin.lock().lines().map(|l| l.expect("read error"));

    let mut player: Option<Player> = None;

    loop {
        // Collect lines for one turn
        let mut turn_lines: Vec<String> = Vec::new();

        // Read until we have both Anfield and Piece sections
        let mut anfield_lines: Vec<String> = Vec::new();
        let mut piece_lines: Vec<String> = Vec::new();
        let mut in_anfield = false;
        let mut in_piece = false;
        let mut anfield_rows_left: usize = 0;
        let mut piece_rows_left: usize = 0;

        loop {
            let line = match lines_iter.next() {
                Some(l) => l,
                None => return, // EOF
            };
            turn_lines.push(line.clone());

            // Detect player assignment
            if line.starts_with("$$$ exec") {
                if let Some(p) = parse_player(&line) {
                    player = Some(p);
                }
                continue;
            }

            if line.starts_with("Anfield") {
                in_anfield = true;
                in_piece = false;
                // Parse dimensions from header
                let parts: Vec<&str> = line.split_whitespace().collect();
                if parts.len() >= 3 {
                    let cols: usize = parts[1].parse().unwrap_or(0);
                    let rows: usize = parts[2].trim_end_matches(':').parse().unwrap_or(0);
                    anfield_rows_left = rows + 1; // +1 for col-index header line
                    anfield_lines.clear();
                    anfield_lines.push(line);
                    let _ = cols; // used indirectly via parse_anfield
                }
                continue;
            }

            if line.starts_with("Piece") {
                in_anfield = false;
                in_piece = true;
                let parts: Vec<&str> = line.split_whitespace().collect();
                if parts.len() >= 3 {
                    piece_rows_left = parts[2].trim_end_matches(':').parse().unwrap_or(0);
                }
                piece_lines.clear();
                piece_lines.push(line);
                continue;
            }

            if in_anfield && anfield_rows_left > 0 {
                anfield_lines.push(line);
                anfield_rows_left -= 1;
                continue;
            }

            if in_piece && piece_rows_left > 0 {
                piece_lines.push(line);
                piece_rows_left -= 1;
                if piece_rows_left == 0 {
                    break; // done reading this turn
                }
                continue;
            }
        }

        let p = match player {
            Some(p) => p,
            None => {
                writeln!(out, "0 0").unwrap();
                out.flush().unwrap();
                continue;
            }
        };

        let anfield = match parse_anfield(&anfield_lines) {
            Some(a) => a,
            None => {
                writeln!(out, "0 0").unwrap();
                out.flush().unwrap();
                continue;
            }
        };

        let piece = match parse_piece(&piece_lines) {
            Some(p) => p,
            None => {
                writeln!(out, "0 0").unwrap();
                out.flush().unwrap();
                continue;
            }
        };

        match best_placement(&anfield, &piece, p) {
            Some((col, row)) => writeln!(out, "{} {}", col, row).unwrap(),
            None => writeln!(out, "0 0").unwrap(),
        }
        out.flush().unwrap();
    }
}

// ─── Unit Tests ──────────────────────────────────────────────────────────────

#[cfg(test)]
mod tests {
    use super::*;

    // ── helpers ──────────────────────────────────────────────────────────────

    fn make_anfield(rows: usize, cols: usize, grid: Vec<&str>) -> Anfield {
        Anfield {
            rows,
            cols,
            grid: grid.iter().map(|r| r.chars().collect()).collect(),
        }
    }

    fn make_piece(rows: usize, cols: usize, cells: Vec<&str>) -> Piece {
        Piece {
            rows,
            cols,
            cells: cells
                .iter()
                .map(|r| r.chars().map(|c| c == '#').collect())
                .collect(),
        }
    }

    // ── parse_player ─────────────────────────────────────────────────────────

    #[test]
    fn test_parse_player_p1() {
        assert_eq!(
            parse_player("$$$ exec p1 : [robots/bender]"),
            Some(Player::P1)
        );
    }

    #[test]
    fn test_parse_player_p2() {
        assert_eq!(
            parse_player("$$$ exec p2 : [robots/wall_e]"),
            Some(Player::P2)
        );
    }

    #[test]
    fn test_parse_player_invalid() {
        assert_eq!(parse_player("$$$ exec p3 : [robots/x]"), None);
    }

    // ── parse_anfield ─────────────────────────────────────────────────────────

    #[test]
    fn test_parse_anfield_basic() {
        let lines: Vec<String> = vec![
            "Anfield 5 3:".to_string(),
            "    01234".to_string(),
            "000 .....".to_string(),
            "001 ..@..".to_string(),
            "002 .....".to_string(),
        ];
        let af = parse_anfield(&lines).expect("should parse");
        assert_eq!(af.rows, 3);
        assert_eq!(af.cols, 5);
        assert_eq!(af.grid[1][2], '@');
        assert_eq!(af.grid[0][0], '.');
    }

    #[test]
    fn test_parse_anfield_with_players() {
        let lines: Vec<String> = vec![
            "Anfield 20 4:".to_string(),
            "    01234567890123456789".to_string(),
            "000 ....................".to_string(),
            "001 .........@..........".to_string(),
            "002 ....................".to_string(),
            "003 .........$..........".to_string(),
        ];
        let af = parse_anfield(&lines).expect("should parse");
        assert_eq!(af.rows, 4);
        assert_eq!(af.cols, 20);
        assert_eq!(af.grid[1][9], '@');
        assert_eq!(af.grid[3][9], '$');
    }

    #[test]
    fn test_parse_anfield_missing_header() {
        let lines: Vec<String> = vec!["BadHeader 5 3:".to_string()];
        assert!(parse_anfield(&lines).is_none());
    }

    // ── parse_piece ───────────────────────────────────────────────────────────

    #[test]
    fn test_parse_piece_basic() {
        let lines: Vec<String> = vec![
            "Piece 4 1:".to_string(),
            ".OO.".to_string(),
        ];
        let p = parse_piece(&lines).expect("should parse");
        assert_eq!(p.rows, 1);
        assert_eq!(p.cols, 4);
        assert!(!p.cells[0][0]);
        assert!(p.cells[0][1]);
        assert!(p.cells[0][2]);
        assert!(!p.cells[0][3]);
    }

    #[test]
    fn test_parse_piece_multi_row() {
        let lines: Vec<String> = vec![
            "Piece 2 2:".to_string(),
            ".#".to_string(),
            "#.".to_string(),
        ];
        let p = parse_piece(&lines).expect("should parse");
        assert_eq!(p.rows, 2);
        assert_eq!(p.cols, 2);
        assert!(!p.cells[0][0]);
        assert!(p.cells[0][1]);
        assert!(p.cells[1][0]);
        assert!(!p.cells[1][1]);
    }

    #[test]
    fn test_parse_piece_hash_symbol() {
        let lines: Vec<String> = vec![
            "Piece 3 2:".to_string(),
            "###".to_string(),
            "#..".to_string(),
        ];
        let p = parse_piece(&lines).expect("should parse");
        assert!(p.cells[0][0]);
        assert!(p.cells[0][1]);
        assert!(p.cells[0][2]);
        assert!(p.cells[1][0]);
        assert!(!p.cells[1][1]);
    }

    // ── is_valid_placement ────────────────────────────────────────────────────

    #[test]
    fn test_valid_placement_exactly_one_overlap() {
        // Grid: player P1 has '@' at (1,1), rest empty
        let af = make_anfield(5, 5, vec![
            ".....",
            ".@...",
            ".....",
            ".....",
            ".....",
        ]);
        // Piece: single filled cell
        let piece = make_piece(1, 1, vec!["#"]);
        // Placing at (1,1) → overlaps own '@' → valid
        assert!(is_valid_placement(&af, &piece, 1, 1, Player::P1));
    }

    #[test]
    fn test_invalid_placement_no_overlap() {
        let af = make_anfield(5, 5, vec![
            ".....",
            ".@...",
            ".....",
            ".....",
            ".....",
        ]);
        let piece = make_piece(1, 1, vec!["#"]);
        // Placing at (3,3) → no own overlap → invalid
        assert!(!is_valid_placement(&af, &piece, 3, 3, Player::P1));
    }

    #[test]
    fn test_invalid_placement_two_own_overlaps() {
        let af = make_anfield(5, 5, vec![
            ".....",
            ".@@..",
            ".....",
            ".....",
            ".....",
        ]);
        // Piece that spans two cells horizontally
        let piece = make_piece(1, 2, vec!["##"]);
        // Placing at (1,1) → overlaps both '@' at (1,1) and (2,1) → invalid
        assert!(!is_valid_placement(&af, &piece, 1, 1, Player::P1));
    }

    #[test]
    fn test_invalid_placement_opponent_overlap() {
        let af = make_anfield(5, 5, vec![
            ".....",
            ".@...",
            "..$. .",
            ".....",
            ".....",
        ]);
        let piece = make_piece(1, 1, vec!["#"]);
        // Placing at (2,2) → overlaps opponent '$' → invalid
        assert!(!is_valid_placement(&af, &piece, 2, 2, Player::P1));
    }

    // ── boundary detection ────────────────────────────────────────────────────

    #[test]
    fn test_boundary_piece_out_of_bounds_right() {
        let af = make_anfield(3, 3, vec![
            "@..",
            "...",
            "...",
        ]);
        let piece = make_piece(1, 3, vec!["###"]);
        // col=1 → piece occupies cols 1,2,3 but grid only has 0..2 → out of bounds
        assert!(!is_valid_placement(&af, &piece, 1, 0, Player::P1));
    }

    #[test]
    fn test_boundary_piece_out_of_bounds_bottom() {
        let af = make_anfield(3, 3, vec![
            "@..",
            "...",
            "...",
        ]);
        // Piece 1x3 vertical
        let piece = Piece {
            rows: 3,
            cols: 1,
            cells: vec![vec![true], vec![true], vec![true]],
        };
        // row=1 → occupies rows 1,2,3 but grid only 0..2 → out of bounds
        assert!(!is_valid_placement(&af, &piece, 0, 1, Player::P1));
    }

    #[test]
    fn test_boundary_negative_coords_with_valid_overlap() {
        // Piece 2 cols wide at col=-1 with '#' at local (0,0) and (0,1)
        // Only (0,1) lands on grid col 0 which has '@'
        let af = make_anfield(3, 3, vec![
            "@..",
            "...",
            "...",
        ]);
        let piece = make_piece(1, 2, vec!["##"]);
        // col=-1, row=0 → cell(0,0)→ gc=-1 (out of bounds, but it's a '#' → invalid)
        assert!(!is_valid_placement(&af, &piece, -1, 0, Player::P1));
    }

    #[test]
    fn test_boundary_negative_empty_cells_ok() {
        // Piece: ".#" at col=-1, row=0
        // cell(0,0) is '.', cell(0,1) lands on gc=0 which is '@' → exactly one overlap
        let af = make_anfield(3, 3, vec![
            "@..",
            "...",
            "...",
        ]);
        let piece = make_piece(1, 2, vec![".#"]);
        // cell(0,0) not filled → skip; cell(0,1) at gc=0,gr=0 → '@' → overlap=1 → valid
        assert!(is_valid_placement(&af, &piece, -1, 0, Player::P1));
    }

    // ── best_placement ────────────────────────────────────────────────────────

    #[test]
    fn test_best_placement_finds_valid_move() {
        let af = make_anfield(5, 5, vec![
            ".....",
            ".@...",
            ".....",
            ".....",
            ".$...",
        ]);
        let piece = make_piece(1, 1, vec!["#"]);
        let result = best_placement(&af, &piece, Player::P1);
        assert!(result.is_some());
        let (col, row) = result.unwrap();
        // Must be a valid placement
        assert!(is_valid_placement(&af, &piece, col, row, Player::P1));
    }

    #[test]
    fn test_best_placement_no_valid_move() {
        // Tiny grid, fully occupied by opponent — no valid spot
        let af = make_anfield(1, 1, vec!["$"]);
        let piece = make_piece(1, 1, vec!["#"]);
        // No '@' exists so no own overlap possible AND opponent blocks
        assert!(best_placement(&af, &piece, Player::P1).is_none());
    }

    #[test]
    fn test_best_placement_prefers_opponent_direction() {
        // P1 starts top-left, P2 is bottom-right
        // Piece is single cell — best move should be toward opponent
        let af = make_anfield(5, 5, vec![
            "@....",
            ".....",
            ".....",
            ".....",
            "....$",
        ]);
        let piece = make_piece(1, 1, vec!["#"]);
        let result = best_placement(&af, &piece, Player::P1);
        assert!(result.is_some());
        let (col, row) = result.unwrap();
        // Placement must be valid and must touch '@' at (0,0)
        // Valid positions for a 1x1 piece overlapping '@': only (0,0) itself
        // Adjacent empty cells (1,0) and (0,1) also have exactly 1 own overlap
        assert!(is_valid_placement(&af, &piece, col, row, Player::P1));
        // The scorer should pick a move that is valid — any of (0,0),(1,0),(0,1)
        let valid_positions = [(0,0),(1,0),(0,1)];
        assert!(
            valid_positions.contains(&(col, row)),
            "Expected a position adjacent to own territory, got ({}, {})", col, row
        );
    }

    // ── player chars ─────────────────────────────────────────────────────────

    #[test]
    fn test_player_chars_p1() {
        let (lower, upper) = Player::P1.my_chars();
        assert_eq!(lower, 'a');
        assert_eq!(upper, '@');
    }

    #[test]
    fn test_player_chars_p2() {
        let (lower, upper) = Player::P2.my_chars();
        assert_eq!(lower, 's');
        assert_eq!(upper, '$');
    }

    #[test]
    fn test_opponent_chars_p1_sees_p2() {
        let (lower, upper) = Player::P1.opp_chars();
        assert_eq!(lower, 's');
        assert_eq!(upper, '$');
    }

    // ── placement with lowercase territory chars ──────────────────────────────

    #[test]
    fn test_valid_placement_lowercase_my_char() {
        // 'a' is the lowercase version of P1's territory (last placed piece)
        let af = make_anfield(3, 3, vec![
            "a..",
            "...",
            "...",
        ]);
        let piece = make_piece(1, 1, vec!["#"]);
        // Overlap with 'a' counts as own overlap
        assert!(is_valid_placement(&af, &piece, 0, 0, Player::P1));
    }

    #[test]
    fn test_invalid_placement_lowercase_opp_char() {
        // 's' is the lowercase version of P2's territory
        let af = make_anfield(3, 3, vec![
            "@..",
            "...",
            "..s",
        ]);
        let piece = make_piece(1, 1, vec!["#"]);
        // Overlapping 's' → invalid
        assert!(!is_valid_placement(&af, &piece, 2, 2, Player::P1));
    }
}