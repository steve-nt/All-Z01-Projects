#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Cell {
    Empty,
    Mine,
    Theirs,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum PlayerId {
    P1,
    P2,
}

impl PlayerId {
    pub fn mine_chars(self) -> (char, char) {
        match self {
            PlayerId::P1 => ('@', 'a'),
            PlayerId::P2 => ('$', 's'),
        }
    }

    pub fn theirs_chars(self) -> (char, char) {
        match self {
            PlayerId::P1 => ('$', 's'),
            PlayerId::P2 => ('@', 'a'),
        }
    }

    pub fn classify(self, ch: char) -> Cell {
        let (m1, m2) = self.mine_chars();
        let (t1, t2) = self.theirs_chars();
        if ch == m1 || ch == m2 {
            Cell::Mine
        } else if ch == t1 || ch == t2 {
            Cell::Theirs
        } else {
            Cell::Empty
        }
    }
}

#[derive(Debug, Clone)]
pub struct Board {
    pub w: u16,
    pub h: u16,
    pub cells: Vec<Cell>,
}

impl Board {
    pub fn new(w: u16, h: u16) -> Self {
        Self {
            w,
            h,
            cells: vec![Cell::Empty; (w as usize) * (h as usize)],
        }
    }

    pub fn idx(&self, x: u16, y: u16) -> usize {
        (y as usize) * (self.w as usize) + (x as usize)
    }

    pub fn get(&self, x: u16, y: u16) -> Cell {
        self.cells[self.idx(x, y)]
    }

    pub fn set(&mut self, x: u16, y: u16, c: Cell) {
        let i = self.idx(x, y);
        self.cells[i] = c;
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn player_id_mine_chars() {
        assert_eq!(PlayerId::P1.mine_chars(), ('@', 'a'));
        assert_eq!(PlayerId::P2.mine_chars(), ('$', 's'));
    }

    #[test]
    fn player_id_theirs_chars() {
        assert_eq!(PlayerId::P1.theirs_chars(), ('$', 's'));
        assert_eq!(PlayerId::P2.theirs_chars(), ('@', 'a'));
    }

    #[test]
    fn classify_p1() {
        let p = PlayerId::P1;
        assert_eq!(p.classify('@'), Cell::Mine);
        assert_eq!(p.classify('a'), Cell::Mine);
        assert_eq!(p.classify('$'), Cell::Theirs);
        assert_eq!(p.classify('s'), Cell::Theirs);
        assert_eq!(p.classify('.'), Cell::Empty);
    }

    #[test]
    fn classify_p2() {
        let p = PlayerId::P2;
        assert_eq!(p.classify('$'), Cell::Mine);
        assert_eq!(p.classify('s'), Cell::Mine);
        assert_eq!(p.classify('@'), Cell::Theirs);
        assert_eq!(p.classify('a'), Cell::Theirs);
        assert_eq!(p.classify('.'), Cell::Empty);
    }

    #[test]
    fn board_idx_and_get_set() {
        let mut b = Board::new(4, 3);
        assert_eq!(b.idx(0, 0), 0);
        assert_eq!(b.idx(3, 2), 11);
        assert_eq!(b.get(2, 1), Cell::Empty);
        b.set(2, 1, Cell::Mine);
        assert_eq!(b.get(2, 1), Cell::Mine);
    }
}
